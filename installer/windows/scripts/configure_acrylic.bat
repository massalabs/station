@ECHO OFF

:: Set Acrylic DNS to resolve `.massa` TLD to localhost

ECHO Executing configure_acrylic.bat

SETLOCAL ENABLEDELAYEDEXPANSION

SET "acrylic_config=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt"
SET "acrylic_configuration=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicConfiguration.ini"
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


:: Add the TLD to the file
ECHO: >> "%acrylic_config%"
ECHO 127.0.0.1 *.massa >> "%acrylic_config%"
ECHO ::1 *.massa >> "%acrylic_config%"

:: get line with LocalIPv4BindingAddress=0.0.0.0 and set it to LocalIPv4BindingAddress=127.0.0.1
FOR /F "delims=:" %%A IN ('findstr /n /c:"LocalIPv4BindingAddress=0.0.0.0" "%acrylic_configuration%"') DO SET line=%%A
IF DEFINED line (
    :: Create a temporary file
    SET "temp_config=%temp%\temp_acrylic_config.ini"
    :: Use `more` command to read the file line by line
    FOR /F "tokens=1* delims=:" %%A IN ('more +0 "%acrylic_configuration%"') DO (
        IF %%A EQU %line% (
            ECHO LocalIPv4BindingAddress=127.0.0.1 >> "%temp_config%"
        ) ELSE (
            ECHO %%A:%%B >> "%temp_config%"
        )
    )
    :: Replace the original configuration file with the temporary one
    MOVE /Y "%temp_config%" "%acrylic_configuration%"
)

:: Restart the Acrylic DNS Proxy Service
NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"
if %errorlevel% NEQ 0 (
    ECHO "Failed to restart Acrylic DNS Proxy Service"
    EXIT 1
)

ENDLOCAL

EXIT 0
