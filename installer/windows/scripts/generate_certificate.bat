@ECHO OFF

:: Generate a certificate for the `.massa` TLD

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

:: Generate a certificate for the TLD
mkcert.exe --cert-file "%CONFIG_DIR%\cert.pem" --key-file "%CONFIG_DIR%\cert-key.pem" station.massa
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
