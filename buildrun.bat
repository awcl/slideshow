@echo off
setlocal

set "EXE_NAME=run.exe"

if exist "%EXE_NAME%" (
    del "%EXE_NAME%"
)

echo Building...

go build -o "%EXE_NAME%" slideshow.go

if errorlevel 1 (
    echo Build failed. Exiting...
    exit /b 1
)

echo Build succeeded

echo Running...

start "" /b "%EXE_NAME%"

setlocal enabledelayedexpansion
set "IP_ADDRESSES="

for /f "tokens=2 delims=:" %%I in ('ipconfig ^| findstr /i "IPv4"') do (
    set "IP_ADDRESS=%%I"
    set "IP_ADDRESS=!IP_ADDRESS:~1!"
    if defined IP_ADDRESS (
        if not "!IP_ADDRESSES!"=="" set "IP_ADDRESSES=!IP_ADDRESSES!, "
        set "IP_ADDRESSES=!IP_ADDRESSES!!IP_ADDRESS!"
    )
)

echo.
echo Users on the LAN(s) can access it at:
for %%A in (%IP_ADDRESSES%) do echo http://%%A:3000/
echo.
echo If you want to stop the application, please close this command window

pause
