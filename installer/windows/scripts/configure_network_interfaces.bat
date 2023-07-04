@echo off

SETLOCAL ENABLEDELAYEDEXPANSION

:: List available Network Interfaces and set dnsservers to 127.0.0.1 (localhost) and ::1 (IPv6 localhost)
:: This is required for Acrylic to work properly
for /f "skip=1 delims=" %%A in ('wmic nic where "netenabled=true" get netconnectionID') do @for /f "delims=" %%B in ("%%A") do (
    SET "networkAdapterName=%%B"

    call :WriteToLog "Configuring !networkAdapterName: =!..."

    NETSH interface ipv4 set dnsservers "!networkAdapterName: =!" static 127.0.0.1 primary
    NETSH interface ipv6 set dnsservers "!networkAdapterName: =!" static ::1 primary
    if %errorlevel% NEQ 0 (
        call :WriteToLog "Failed to configure !networkAdapterName: =!"
        EXIT 1
    )
)

ENDLOCAL

call :WriteToLog "Success"
EXIT 0

:: decalre a function to log and print a message
:WriteToLog
  setlocal
  set LOG_MESSAGE=%~1
  echo %LOG_MESSAGE%
  set LOG_FILE=%TEMP%\massa-station-install-configure-network-interfaces.log

  echo %LOG_MESSAGE% >> %LOG_FILE%

  endlocal
exit /b
