@ECHO OFF

echo Executing configure_network_interfaces.bat

SETLOCAL ENABLEDELAYEDEXPANSION

:: List available Network Interfaces and set dnsservers to 127.0.0.1 (localhost) and ::1 (IPv6 localhost)
:: This is required for Acrylic to work properly
for /f "skip=1 delims=" %%A in ('wmic nic where "netenabled=true" get netconnectionID') do @for /f "delims=" %%B in ("%%A") do (
    CALL :TRIM %%B
    SET networkAdapterName=!TRIMRESULT!

    ECHO Configuring !networkAdapterName!...

    NETSH interface ipv4 set dnsservers "!networkAdapterName!" static 127.0.0.1 primary
    if !errorlevel! NEQ 0 (
        ECHO "Failed to configure !networkAdapterName!"
        EXIT 1
    )

    NETSH interface ipv6 set dnsservers "!networkAdapterName!" static ::1 primary
    if !errorlevel! NEQ 0 (
        ECHO "Failed to configure !networkAdapterName!"
        EXIT 1
    )
)

ENDLOCAL

EXIT 0


:: Removes leading and trailing spaces from a string passed as arguments. The result is stored in TRIMRESULT.
:TRIM
    SET TRIMRESULT=%*
GOTO :EOF
