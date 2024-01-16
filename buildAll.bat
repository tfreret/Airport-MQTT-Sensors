@echo off
setlocal enabledelayedexpansion

set "source_dir=./cmd"
set "build_dir=./exe"

if not exist "%build_dir%" mkdir "%build_dir%"

del /Q "%build_dir%\*"

set "files=/sensors/pressureSensor /sensors/tempSensor /sensors/windSensor databaseRecorder fileRecorder api"

for %%f in (%files%) do (
    echo Building %%f
    go build -o "%build_dir%\%%f.exe" "%source_dir%\%%f"
)

echo Build process completed.

 start "" "./%build_dir%/sensors/pressureSensor.exe"
 start "" "./%build_dir%/sensors/tempSensor.exe"
 start "" "./%build_dir%/sensors/windSensor.exe"
 start "" "./%build_dir%/database-recorder.exe"
 start "" "./%build_dir%/api.exe"

::cleanup
::echo Sensors stopping...
::for %%f in (%files%) do (
::    taskkill /IM "%%f.exe" /F
::)
::echo Stopped
::exit /B 0