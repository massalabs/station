@ECHO OFF

:: Set Acrylic DNS to resolve `.massa` TLD to localhost

SETLOCAL ENABLEDELAYEDEXPANSION

SET "acrylic_config=C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt"

:: If it doesn't exist, create it
if not exist "%acrylic_config%" (
    :: Create the file
    ECHO # Acrylic DNS Proxy configuration file > "%acrylic_config%"
    ECHO # For more information please visit http://mayakron.altervista.org >> "%acrylic_config%"
    ECHO # >> "%acrylic_config%"
)

:: Check if the TLD is already in the file
FINDSTR /c:".massa" "%acrylic_config%" >nul 2>&1
if %errorlevel%==0 (
    call :WriteToLog "TLD already in the file"
    goto :EOF
)


:: Add the TLD to the file
ECHO: >> "%acrylic_config%"
ECHO 127.0.0.1 *.massa >> "%acrylic_config%"
ECHO ::1 *.massa >> "%acrylic_config%"

:: Restart the Acrylic DNS Proxy Service
NET STOP "AcrylicDNSProxySvc"
NET START "AcrylicDNSProxySvc"
if %errorlevel% NEQ 0 (
    call :WriteToLog "Failed to restart Acrylic DNS Proxy Service"
    pause
    EXIT 1
)

ENDLOCAL

call :WriteToLog "Success"
pause
EXIT 0

:: decalre a function to log and print a message
:WriteToLog
  setlocal
  set LOG_MESSAGE=%~1
  echo %LOG_MESSAGE%
  set LOG_FILE=%TEMP%\massa-station-install-configure-acrylic.log

  echo %LOG_MESSAGE% >> %LOG_FILE%

  endlocal
exit /b
