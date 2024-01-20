package main

import (
	"airport/internal/config"
	"airport/internal/mqttTools"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	var (
		influxEnvFile = flag.String("influx", "influxdb.env", ".env file for influx db")
	)
	var (
		configMQTTFile = flag.String("config", "database-recorder.yaml", "Config file of the database recorder")
	)
	flag.Parse()

	configMQTT := config.ReadConfig[mqttTools.MonoConfigMqtt](*configMQTTFile)
	configInfluxDB := config.ReadEnv[mqttTools.ConfigInfluxDB](*influxEnvFile)

	fmt.Println("Using config :", configMQTT, configInfluxDB)

	InfluxDBClient := influxdb2.NewClient(configInfluxDB.InfluxDBURL, configInfluxDB.InfluxDBToken)
	writeAPI := InfluxDBClient.WriteAPI(configInfluxDB.InfluxDBOrg, configInfluxDB.InfluxDBBucket)

	brokerClient := mqttTools.NewBrokerClient(
		configMQTT.Mqtt,
	)

	brokerClient.Subscribe("data/#", configMQTT.Mqtt.MqttQOS, func(topic string, message []byte) {
		iata, measure, sensorId, err := mqttTools.ParseTopic(topic)
		if err != nil {
			fmt.Println("Couldn't extract IATA code; measure, and sensorId type from string : " + topic)
			return
		}

		value, date, err := mqttTools.ParseData(string(message))
		if err != nil {
			fmt.Println("Couldn't extract value and time from payload : " + string(message))
		}

		valuefloat, _ := strconv.ParseFloat(value, 64)
		timestamp, _ := time.Parse(time.RFC3339, date)

		tags := map[string]string{
			"airport_id":      iata,
			"sensor_category": measure,
			"sensor_id":       sensorId,
		}
		fields := map[string]interface{}{
			"value": valuefloat,
		}
		point := influxdb2.NewPoint("sensor_data", tags, fields, timestamp)

		fmt.Println("Insert in inlfuxDB : ", tags, valuefloat, timestamp)

		writeAPI.WritePoint(point)
	})

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal
	InfluxDBClient.Close()
}
