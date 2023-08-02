@ECHO OFF

:: Set Acrylic DNS to resolve `.massa` TLD to localhost

ECHO Executing configure_acrylic.bat

SETLOCAL ENABLEDELAYEDEXPANSION

SET "acrylic_config=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt"
SET "acrylic_configuration=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicConfiguration.ini"
SET "source_file=config.ini"
SET "target_file=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicConfiguration.ini"

:: Remove the line containing "LocalIPv4BindingAddress"
findstr /V "LocalIPv4BindingAddress" "%acrylic_configuration%" > temp.txt

:: Replace the original file with the temporary one
MOVE /Y temp.txt "%acrylic_configuration%"

:: Add the new LocalIPv4BindingAddress setting at the end of the file
ECHO LocalIPv4BindingAddress=127.0.0.1 >> "%acrylic_configuration%"

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
    ECHO TLD already in the file
    EXIT 0
)

:: Replace the target file with the source file
COPY /Y "%source_file%" "%target_file%"


:: Add the TLD to the file
ECHO: >> "%acrylic_config%"
ECHO 127.0.0.1 *.massa >> "%acrylic_config%"
ECHO ::1 *.massa >> "%acrylic_config%"

:: Restart the Acrylic DNS Proxy Service
NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"
if %errorlevel% NEQ 0 (
    ECHO "Failed to restart Acrylic DNS Proxy Service"
    EXIT 1
)

ENDLOCAL

EXIT 0
