#!/bin/bash

./exe/sensors/pressureSensor -config="../configs/demo/sensor1.yaml" &
./exe/sensors/tempSensor -config="../configs/demo/sensor2.yaml" &
./exe/sensors/windSensor -config="../configs/demo/sensor3.yaml" &
./exe/sensors/pressureSensor -config="../configs/demo/sensor4.yaml" &
./exe/sensors/tempSensor -config="../configs/demo/sensor5.yaml" &
./exe/sensors/windSensor -config="../configs/demo/sensor6.yaml" &

main_pid=$$

wait

pkill -f 'pressureSensor|tempSensor|windSensor'
