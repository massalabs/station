' Create a WScript Shell Object
Set oShell = CreateObject("Wscript.Shell")

' Command to be run
Dim strArgs
strArgs = "massastation.exe"

' Use the Run method of the WScript Shell object to run the command.
' The second parameter 0 means "hidden window"
' The third parameter false means the script will not wait for the program to finish execution before continuing to the next statement
oShell.Run strArgs, 0, false