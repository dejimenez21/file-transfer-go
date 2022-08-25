@echo off
cd src\client\
go build
rmdir /s /q "C:\CFTP-DEMO\1\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\2\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\3\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\4\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\5\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\6\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\7\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\8\ReceivedFiles\"
rmdir /s /q "C:\CFTP-DEMO\9\ReceivedFiles\"

xcopy client.exe "C:\CFTP-DEMO\1\" /y  
xcopy client.exe "C:\CFTP-DEMO\2\" /y
xcopy client.exe "C:\CFTP-DEMO\3\" /y
xcopy client.exe "C:\CFTP-DEMO\4\" /y
xcopy client.exe "C:\CFTP-DEMO\5\" /y
xcopy client.exe "C:\CFTP-DEMO\6\" /y
xcopy client.exe "C:\CFTP-DEMO\7\" /y
xcopy client.exe "C:\CFTP-DEMO\8\" /y
xcopy client.exe "C:\CFTP-DEMO\9\" /y


