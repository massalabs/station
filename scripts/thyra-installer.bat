@echo off

@REM Execute a command requiring admin rights that is present on most Windows versions 
fltmc >nul 2>&1 && (
    ECHO "Admin rights successfully detected."
) || (
  ECHO "Couldn't detect admin rights. Please execute this script as an administator."
  PAUSE
  EXIT
)

curl -L https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64 --output thyra-server.exe

SET ACRYLIC_PATH="C:\Program Files (x86)\Acrylic DNS Proxy"
@REM %~dp0 is alias for current working directory
SET STARTING_WORKING_DIR=%~dp0

IF NOT EXIST %STARTING_WORKING_DIR%thyra-server.exe (
    ECHO "The installer couldn't download Thyra from Github. Please retry or check your internet connection."
    PAUSE
    EXIT
)

IF NOT EXIST %ACRYLIC_PATH%\AcrylicHosts.txt (
    curl -L https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic-Portable.zip/download --output Acrylic-Portable.zip
    IF NOT EXIST %STARTING_WORKING_DIR%Acrylic-Portable.zip (
        ECHO "The installer couldn't download DNS Acrylic Proxy. Please retry or check your internet connection."
        PAUSE
        EXIT
    )

    MKDIR %ACRYLIC_PATH%
    IF NOT EXIST %ACRYLIC_PATH%\ (
        ECHO "Couldn't create DNS Acrylic Proxy installation folder. Please make sure this script can access C:\Program Files (x86)."
        PAUSE
        EXIT
    )
    MOVE %STARTING_WORKING_DIR%Acrylic-Portable.zip %ACRYLIC_PATH%\
    tar -xf  %ACRYLIC_PATH%\Acrylic-Portable.zip -C %ACRYLIC_PATH%
    @REM Ideally we should check here that Acrylic has been successfully unarchived.
    CD %ACRYLIC_PATH%

    CALL "InstallAcrylicService.bat"

    @REM ENABLEDELAYEDEXPANSION Allows for variable value assertion inside for loops at runtime
    @REM This for loop iterate over all connected network adapters and set their DNS.
    @REM Data is forwarded in another for loop to automatically clean wmic output of trailing <CR> tags
    SETLOCAL ENABLEDELAYEDEXPANSION
    for /f "skip=1 delims=" %%A in ('wmic nic where "netenabled=true" get netconnectionID') do @for /f "delims=" %%B in ("%%A") do (
        SET "networkAdapterName=%%B"
        @REM !networkAdapterName: =! remove all spaces from networkAdapterName variable
        NETSH interface ipv4 set dnsservers "!networkAdapterName: =!" static 127.0.0.1 primary
        NETSH interface ipv6 set dnsservers "!networkAdapterName: =!" static ::1 primary)

    DEL %ACRYLIC_PATH%\Acrylic-Portable.zip
    CD %STARTING_WORKING_DIR%
)

@REM Ideally before doing this we should check that "127.0.0.1 *.massa" is not already written in the file
ECHO: >> %ACRYLIC_PATH%\AcrylicHosts.txt
ECHO 127.0.0.1 *.massa >> %ACRYLIC_PATH%\AcrylicHosts.txt

NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"

if not exist %homedrive%%homepath%\.config\thyra mkdir %homedrive%%homepath%\.config\thyra

ECHO Installation and setup successfull

PAUSE
EXIT