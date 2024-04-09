# Distributed Architecture Project

## Group Members
- **Tom Freret**
- **Louis Painter**
- **Matthis Bleuet**
- **Antoine Otegui**

## Overview
This project, implemented in Go, aims to simulate sensors (pressure, temperature, wind) in the airport context. It employs a distributed architecture where microservices communicate via MQTT.

![Architecture Diagram](./assets/architecture.png).

## Project Organization
Utilization of a Kanban board (JIRA) for project management, in conjunction with Discord for communication.

## Technical Choices
- **InfluxDB**
- **HiveMQ** as MQTT broker (alternative brokers such as Mosquitto are possible)
- A few libraries:
  - **Logrus** for logging
  - **Paho** for MQTT
  - **Swaggo** for Swagger generation
  - **Viper** for loading configurations and others
  - **Mux** for the API
- **Next.js** project for the Single Page Application (SPA)

## Project Structure:
- `assets`: Images and other assets
- `/cmd/`: Source files for services
- `/configs`: Configuration files for microservices
- `/docs`: Generated Swagger data
- `/exe`: Binaries for all services
- `/internal`: Common source code
- `/outputs`: Saved .csv files
- `/web`: The SPA
- `/outputs`: Data .csv files

## Compilation and Execution Instructions

### Docker InfluxDB

To launch the Docker InfluxDB container, use the following command:
```bash
docker compose --env-file ../configs/influxdb.env up
```

### Build and Run

All scripts are available in .bat and .sh.

- `./BuildAll`: Builds the binaries
- `./runAll`: Builds then launches all services with default configurations
- `./demo`: Launches demo sensors at different airports
- `./demoAlert`: Launches a sensor with abnormal values //TODO

To launch services separately (at the root and ensure InfluxDB instance is running and binaries are compiled):

- `./exe/sensors/pressureSensor -config="path/to/config/file"`
- `./exe/sensors/tempSensor -config="path/to/config/file"`
- `./exe/sensors/windSensor -config="path/to/config/file"`
- `./exe/alertManager -config="path/to/config/file"`
- `./exe/fileRecorder -config="path/to/config/file"`
- `./exe/databaseRecorder -config="path/to/config/file" -influx="path/to/env/file"`
- `./exe/api -config="path/to/config/file" -influx="path/to/env/file"`

### Swagger Documentation

Generate Swagger documentation by following these steps:

Ensure the go/bin directory is in your

    $PATH (export PATH=$PATH:$GOPATH/bin).

Run the following command at the root of the project:

    swag init -g cmd/api/main.go

## API:

localhost:8080: The API with various routes:
- .../airports
- .../sensors/sensors/{airportID}
- .../average/{airportID}
- .../average/{airportID}/{sensorType}
- .../data/{airportID}/{sensorType}/{sensorID}
- .../swagger/: For more information and a playground to test different routes

## SPA:

To launch the SPA, navigate to the .web repo
then install dependencies:

```bash
npm i
```

and then run the app

```bash
npm run dev
```
