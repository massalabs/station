@ECHO OFF

@REM Set Acrylic DNS to resolve `.massa` TLD to localhost

ECHO Executing configure_acrylic.bat

if "%~1"=="" (
    ECHO Please provide the path to Acrylic DNS Proxy installation folder
    EXIT 1
)

SETLOCAL ENABLEDELAYEDEXPANSION

@REM Get Acrylic Path from command line argument
SET acrylic_path=%~1

SET acrylic_config_host=%acrylic_path%\AcrylicHosts.txt
SET acrylic_configuration=%acrylic_path%\AcrylicConfiguration.ini

@REM AcrylicHosts.txt file does not exist, we create it
if not exist "%acrylic_config_host%" (
    ECHO # Acrylic DNS Proxy configuration file > "%acrylic_config_host%"
    ECHO # For more information please visit http://mayakron.altervista.org >> "%acrylic_config_host%"
    ECHO # >> "%acrylic_config_host%"
)

@REM If LocalIPv4BindingAddress is set to default 0.0.0.0, we set it to 127.0.0.1 for backward compatibility with Windows 10
FINDSTR /c:"LocalIPv4BindingAddress=127.0.0.1" "%acrylic_configuration%" >nul 2>&1
if %errorlevel%==0 (
    ECHO LocalIPv4BindingAddress already set to 127.0.0.1, skipping
) else (
    ECHO LocalIPv4BindingAddress not set to 127.0.0.1, setting it
    powershell -Command "Get-Content '%acrylic_configuration%' | %%{$_ -replace 'LocalIPv4BindingAddress=0.0.0.0','LocalIPv4BindingAddress=127.0.0.1'} | Out-File -Encoding UTF8 '%acrylic_configuration%.tmp'"
    FINDSTR /c:"LocalIPv4BindingAddress=127.0.0.1" "%acrylic_configuration%.tmp" >nul 2>&1
    if !errorlevel!==0 (
        MOVE /Y "%acrylic_configuration%.tmp" "%acrylic_configuration%"
    ) else (
        ECHO Failed to set LocalIPv4BindingAddress to 127.0.0.1, aborting this step
    )
)

@REM If .massa TLD is already in AcrylicHosts.txt, we can exit
FINDSTR /c:".massa" "%acrylic_config_host%" >nul 2>&1
if %errorlevel%==0 (
    ECHO TLD already in the file
    EXIT 0
)

@REM Add .massa TLD to AcrylicHosts.txt
ECHO: >> "%acrylic_config_host%"
ECHO 127.0.0.1 *.massa >> "%acrylic_config_host%"
ECHO ::1 *.massa >> "%acrylic_config_host%"

@REM Restart Acrylic DNS Proxy Service
NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"
if %errorlevel% NEQ 0 (
    ECHO "Failed to restart Acrylic DNS Proxy Service"
    EXIT 1
)

ENDLOCAL

EXIT 0
