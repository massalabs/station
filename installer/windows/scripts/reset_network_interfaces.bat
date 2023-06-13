@echo off

SETLOCAL ENABLEDELAYEDEXPANSION

:: List available Network Interfaces and set dnsservers to DHCP
for /f "skip=1 delims=" %%A in ('wmic nic where "netenabled=true" get netconnectionID') do @for /f "delims=" %%B in ("%%A") do (
    SET "networkAdapterName=%%B"

    echo Configuring !networkAdapterName: =!...

    NETSH interface ipv4 set dnsservers "!networkAdapterName: =!" dhcp
    NETSH interface ipv6 set dnsservers "!networkAdapterName: =!" dhcp
    if %errorlevel% NEQ 0 (
        ECHO "Failed to configure !networkAdapterName: =!"
        EXIT 1
    )

    NETSH interface set interface "!networkAdapterName: =!" disable
    NETSH interface set interface "!networkAdapterName: =!" enable
)

ENDLOCAL

EXIT 0
