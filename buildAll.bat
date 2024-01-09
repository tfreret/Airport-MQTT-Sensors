@echo off
setlocal enabledelayedexpansion

set "source_dir=./cmd"
set "build_dir=./exe"

if not exist "%build_dir%" mkdir "%build_dir%"

del /Q "%build_dir%\*"

set "files=/sensors/pressureSensor /sensors/tempSensor /sensors/windSensor"

for %%f in (%files%) do (
    echo Building %%f
    go build -o "%build_dir%\%%f.exe" "%source_dir%\%%f"
)

echo Build process completed.

 start "" "./%build_dir%/sensors/pressureSensor.exe" -id 44 -airport KJFK -frequency 10
 start "" "./%build_dir%/sensors/tempSensor.exe" -id 55 -airport KJFK -frequency 10 &
 start "" "./%build_dir%/sensors/windSensor.exe" -id 66 -airport KJFK -frequency 10 &

::cleanup
::echo Sensors stopping...
::for %%f in (%files%) do (
::    taskkill /IM "%%f.exe" /F
::)
::echo Stopped
::exit /B 0