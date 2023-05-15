@ECHO OFF

:: Generate a certificate for the `.massa` TLD

SETLOCAL ENABLEDELAYEDEXPANSION

:: Download the latest mkcert release
curl -L "https://dl.filippo.io/mkcert/latest?for=windows/amd64" --output mkcert.exe

:: Install the local CA
mkcert.exe --install

:: Make sure the config directory exists
SET CONFIG_DIR=%homedrive%%homepath%\.config\thyra\certs
IF NOT EXIST "%CONFIG_DIR%" (
    MKDIR "%CONFIG_DIR%"
)

:: Generate a certificate for the TLD
mkcert.exe --cert-file %CONFIG_DIR%\cert.pem --key-file %CONFIG_DIR%\cert-key.pem my.massa

:: Delete the mkcert executable
DEL mkcert.exe

ENDLOCAL

EXIT 0
