@ECHO OFF

:: Generate a certificate for the `.massa` TLD

set LOG_FILE=%TEMP%\massa-station-install-generate-certificate.log

:: redirect err and std output of all intructions bellow to the log file 
(

SETLOCAL ENABLEDELAYEDEXPANSION

:: Download the latest mkcert release
curl -L "https://dl.filippo.io/mkcert/latest?for=windows/amd64" --output mkcert.exe
if %ERRORLEVEL% NEQ 0 (
    ECHO "Failed to download mkcert"
    EXIT 1
)

:: Install the local CA
mkcert.exe --install
if %ERRORLEVEL% NEQ 0 (
    ECHO "Failed to install mkcert"
    EXIT 1
)

SET "CERTS_DIR=%~dp0certs"

:: Generate a certificate for the TLD
mkcert.exe --cert-file "%CERTS_DIR%\cert.pem" --key-file "%CERTS_DIR%\cert-key.pem" station.massa
if %ERRORLEVEL% NEQ 0 (
    ECHO "Failed to generate certificate"
    EXIT 1
)

:: Delete the mkcert executable
DEL mkcert.exe
if %ERRORLEVEL% NEQ 0 (
    ECHO "Failed to delete mkcert"
    EXIT 1
)

ENDLOCAL

EXIT 0

) >> %LOG_FILE% 2>&1
