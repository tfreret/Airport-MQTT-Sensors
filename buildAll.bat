@echo off
setlocal enabledelayedexpansion

set "./cmd"
set "./exe"

if not exist "%build_dir%" mkdir "%build_dir%"

del /Q "%build_dir%\*"

set "files=sensors/pressureSensor sensors/tempSensor sensors/windSensor"

for %%f in (%files%) do (
    echo Building %%f
    go build -o "%build_dir%\%%f.exe" "%source_dir%\%%f"
)

echo "Build process completed."

for %%f in (%files%) do (
    echo Running %%f
    start "" "%build_dir%\%%f.exe" -id 44 -airport A234C -frequency 10
)

:: ."/$build_dir/tempSensor" -id 44 -airport A234C -frequency 10 &
:: ."/$build_dir/tempSensor" -id 55 -airport A333C -frequency 10 &
:: ."/$build_dir/tempSensor" -id 66 -airport A555C -frequency 10 &

:cleanup
echo Sensors stopping...
for %%f in (%files%) do (
    taskkill /IM "%%f.exe" /F
)
echo Stopped
exit /B 0