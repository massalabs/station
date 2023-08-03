@ECHO OFF

:: Set Acrylic DNS to resolve `.massa` TLD to localhost

ECHO Executing configure_acrylic.bat

SETLOCAL ENABLEDELAYEDEXPANSION

SET "acrylic_config_host=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt"
SET "acrylic_configuration=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicConfiguration.ini"
SET "target_file=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicConfiguration.ini"

:: If it doesn't exist, create it
if not exist "%acrylic_config_host%" (
    :: Create the file
    ECHO # Acrylic DNS Proxy configuration file > "%acrylic_config_host%"
    ECHO # For more information please visit http://mayakron.altervista.org >> "%acrylic_config_host%"
    ECHO # >> "%acrylic_config_host%"
)

:: Check if the TLD is already in the file
FINDSTR /c:".massa" "%acrylic_config_host%" >nul 2>&1
if %errorlevel%==0 (
    ECHO TLD already in the file
    EXIT 0
)

:: change LocalIPv4BindingAddress to 127.0.0.1
findstr /V "LocalIPv4BindingAddress" "%acrylic_configuration%" > temp.txt
MOVE /Y temp.txt "%acrylic_configuration%"
ECHO. >> "%acrylic_configuration%"
ECHO LocalIPv4BindingAddress=127.0.0.1 >> "%acrylic_configuration%"


:: Add the TLD to the file
ECHO: >> "%acrylic_config_host%"
ECHO 127.0.0.1 *.massa >> "%acrylic_config_host%"
ECHO ::1 *.massa >> "%acrylic_config_host%"

:: Restart the Acrylic DNS Proxy Service
NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"
if %errorlevel% NEQ 0 (
    ECHO "Failed to restart Acrylic DNS Proxy Service"
    EXIT 1
)

ENDLOCAL

EXIT 0
