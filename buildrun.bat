@echo off
setlocal

set "EXE_NAME=run.exe"

if exist %EXE_NAME% (
    del %EXE_NAME%
)

echo Building Go application...

go build -o %EXE_NAME% slideshow.go

if errorlevel 1 (
    echo Build failed. Exiting...
    exit /b 1
)

echo Build succeeded.

echo Running the application...

start "" /b %EXE_NAME%

for /f "tokens=2 delims=:" %%I in ('"ipconfig | findstr /i "IPv4""') do set IP_ADDRESS=%%I

set IP_ADDRESS=%IP_ADDRESS:~1%

echo.
echo Your application is now running.
echo Users on the LAN can access it at: http://%IP_ADDRESS%:3000/
echo.
echo If you want to stop the application, please close this command window.
echo.
