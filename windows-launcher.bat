@echo off
setlocal

REM Start ffmpeg in background (if needed)
echo Starting FFmpeg...
start "" ffmpeg.exe

REM Create shortcut on Desktop for video-cutter.exe
set SHORTCUT_NAME=Video Cutter
set TARGET=%~dp0video-cutter.exe
set ICON=%~dp0internal\frontend\public\favicon.png
set SHORTCUT_PATH=%USERPROFILE%\Desktop\%SHORTCUT_NAME%.lnk

REM Check if shortcut already exists
if not exist "%SHORTCUT_PATH%" (
    echo Creating shortcut on Desktop...

    powershell -Command ^
    "$WshShell = New-Object -ComObject WScript.Shell; ^
     $Shortcut = $WshShell.CreateShortcut('%SHORTCUT_PATH%'); ^
     $Shortcut.TargetPath = '%TARGET%'; ^
     $Shortcut.IconLocation = '%ICON%'; ^
     $Shortcut.Save()"
) else (
    echo Shortcut already exists: "%SHORTCUT_PATH%"
)

endlocal
