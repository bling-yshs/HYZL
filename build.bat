@echo off
setlocal

:: 执行go build命令
go build

:: 计算当前目录下YzLauncher-windows.exe的md5值
certutil -hashfile YzLauncher-windows.exe MD5 | findstr /v "CertUtil" | findstr /v "MD5" > yzMD5.txt

echo Done.
