@echo off

set LOG_FILE=%TEMP%\massastation_install.log

:: redirect err and std output of all intructions bellow to the log file 
(

echo Executing configure_network_interfaces.bat

SETLOCAL ENABLEDELAYEDEXPANSION

:: List available Network Interfaces and set dnsservers to 127.0.0.1 (localhost) and ::1 (IPv6 localhost)
:: This is required for Acrylic to work properly
for /f "skip=1 delims=" %%A in ('wmic nic where "netenabled=true" get netconnectionID') do @for /f "delims=" %%B in ("%%A") do (
    SET "networkAdapterName=%%B"

    echo "Configuring !networkAdapterName: =!..."

    NETSH interface ipv4 set dnsservers "!networkAdapterName: =!" static 127.0.0.1 primary
    NETSH interface ipv6 set dnsservers "!networkAdapterName: =!" static ::1 primary
    if %errorlevel% NEQ 0 (
        ECHO "Failed to configure !networkAdapterName: =!"
        EXIT 1
    )
)

ENDLOCAL

) >> %LOG_FILE% 2>&1
