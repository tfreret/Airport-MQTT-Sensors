@echo off
setlocal enabledelayedexpansion

set "source_dir=./cmd"
set "build_dir=./exe"

if not exist "%build_dir%" mkdir "%build_dir%"

del /Q "%build_dir%\*"

set "files=/sensors/pressureSensor /sensors/tempSensor /sensors/windSensor databaseRecorder fileRecorder api alertManager"

for %%f in (%files%) do (
    echo Building %%f
    go build -o "%build_dir%\%%f.exe" "%source_dir%\%%f"
)

echo Build process completed.