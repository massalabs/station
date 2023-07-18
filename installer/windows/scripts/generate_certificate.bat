@ECHO OFF

:: Generate a certificate for the `.massa` TLD

ECHO Executing generate_certificate.bat

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

:: Delete the mkcert executable
DEL mkcert.exe
if %ERRORLEVEL% NEQ 0 (
    ECHO "Failed to delete mkcert"
    EXIT 1
)

ENDLOCAL

EXIT 0
