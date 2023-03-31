@echo off

set "url=%1"
set "tempFile=YzLauncher-windows_temp.exe"
set "selfFile=YzLauncher-windows.exe"

:: 下载最新版本的程序文件到本地
powershell -command "& { (New-Object System.Net.WebClient).DownloadFile('%url%', '%tempFile%') }"
if %errorlevel% neq 0 (
    echo Download failed.
    exit /b %errorlevel%
)

:: 替换程序本身
copy /y "%tempFile%" "%selfFile%" >nul
if %errorlevel% neq 0 (
    echo Replace failed.
    exit /b %errorlevel%
)

:: 启动替换后的程序
start "" "%selfFile%"

:: 关闭当前程序
exit /b 0
