@ECHO OFF

:: Set Acrylic DNS to resolve `.massa` TLD to localhost

set LOG_FILE=%TEMP%\massastation_install.log

:: redirect err and std output of all intructions bellow to the log file 
(

echo Executing configure_acrylic.bat

SETLOCAL ENABLEDELAYEDEXPANSION

SET "acrylic_config=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt"

:: If it doesn't exist, create it
if not exist "%acrylic_config%" (
    :: Create the file
    ECHO # Acrylic DNS Proxy configuration file > "%acrylic_config%"
    ECHO # For more information please visit http://mayakron.altervista.org >> "%acrylic_config%"
    ECHO # >> "%acrylic_config%"
)

:: Check if the TLD is already in the file
FINDSTR /c:".massa" "%acrylic_config%" >nul 2>&1
if %errorlevel%==0 (
    echo "TLD already in the file"
    EXIT 0
)


:: Add the TLD to the file
ECHO: >> "%acrylic_config%"
ECHO 127.0.0.1 *.massa >> "%acrylic_config%"
ECHO ::1 *.massa >> "%acrylic_config%"

:: Restart the Acrylic DNS Proxy Service
NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"
if %errorlevel% NEQ 0 (
    echo "Failed to restart Acrylic DNS Proxy Service"
    EXIT 1
)

ENDLOCAL

EXIT 0

) >> %LOG_FILE% 2>&1
