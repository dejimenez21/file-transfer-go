@echo off
cd src\client\
go build
xcopy client.exe "C:\CFTP-DEMO\1\" /y  
xcopy client.exe "C:\CFTP-DEMO\2\" /y
xcopy client.exe "C:\CFTP-DEMO\3\" /y
xcopy client.exe "C:\CFTP-DEMO\4\" /y
xcopy client.exe "C:\CFTP-DEMO\5\" /y
xcopy client.exe "C:\CFTP-DEMO\6\" /y
xcopy client.exe "C:\CFTP-DEMO\7\" /y
xcopy client.exe "C:\CFTP-DEMO\8\" /y
xcopy client.exe "C:\CFTP-DEMO\9\" /y
