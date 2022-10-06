::curl -L https://github.com/massalabs/thyra/releases/latest/download/thyra-server_windows_amd64 --output thyra-server.exe

IF EXIST "C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt" GOTO :install_completed
ELSE
::curl -L https://sourceforge.net/projects/acrylic/files/Acrylic/2.1.1/Acrylic.exe/download --output acrylic_binary.exe
start /wait %~dp0acrylic_binary.exe

:install_completed

ECHO 127.0.0.1 *.massa >> "C:\Program Files (x86)\Acrylic DNS Proxy\AcrylicHosts.txt"

::net stop Acrylic DNS Proxy && net start Acrylic DNS Proxy 

ECHO Installation and setup successfull

PAUSE