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

:: change LocalIPv4BindingAddress to 127.0.0.1
findstr /V "LocalIPv4BindingAddress" "%acrylic_configuration%" > %acrylic_configuration%
ECHO. >> "%acrylic_configuration%"
ECHO LocalIPv4BindingAddress=127.0.0.1 >> "%acrylic_configuration%"

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
