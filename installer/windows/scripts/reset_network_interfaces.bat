@ECHO OFF

ECHO Executing reset_network_interfaces.bat

SETLOCAL ENABLEDELAYEDEXPANSION

:: List available Network Interfaces and set dnsservers to DHCP
for /f "skip=1 delims=" %%A in ('wmic nic where "netenabled=true" get netconnectionID') do @for /f "delims=" %%B in ("%%A") do (
    CALL :TRIM %%B
    SET networkAdapterName=!TRIMRESULT!

    ECHO Configuring !networkAdapterName!...

    NETSH interface ipv4 set dnsservers "!networkAdapterName!" dhcp
    if !errorlevel! NEQ 0 (
        ECHO "Failed to configure !networkAdapterName!"
        EXIT 1
    )

    NETSH interface ipv6 set dnsservers "!networkAdapterName!" dhcp
    if !errorlevel! NEQ 0 (
        ECHO "Failed to configure !networkAdapterName!"
        EXIT 1
    )

    NETSH interface set interface "!networkAdapterName!" disable
    NETSH interface set interface "!networkAdapterName!" enable
)

ENDLOCAL

EXIT 0


:: Removes leading and trailing spaces from a string passed as arguments. The result is stored in TRIMRESULT.
:TRIM
    SET TRIMRESULT=%*
GOTO :EOF
