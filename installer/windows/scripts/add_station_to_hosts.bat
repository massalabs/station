@ECHO OFF

ECHO Executing add_station_to_hosts.bat

SETLOCAL ENABLEDELAYEDEXPANSION

@REM Adds station.massa host to Windows hosts file if not already present

SET windows_hosts_file=%windir%\System32\drivers\etc\hosts

@REM If station.massa is already in Windows hosts file, we can exit
FINDSTR /c:"127.0.0.1 station.massa" "%windows_hosts_file%" >nul 2>&1
if %errorlevel%==0 (
    ECHO station.massa host already configured
    EXIT 0
)

@REM Add station.massa to Windows hosts file
ECHO: >> "%windows_hosts_file%"
ECHO 127.0.0.1 station.massa >> "%windows_hosts_file%"

ENDLOCAL

EXIT 0
