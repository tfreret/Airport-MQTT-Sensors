package main

import (
	"airport/internal/config"
	"airport/internal/mqttTools"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	envFile := filepath.Join("./internal/config", "influxdb.env")

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier influxdb.env :", err)
	}

	InfluxDBClient := influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"))
	writeAPI := InfluxDBClient.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

	brokerClient := mqttTools.NewBrokerClient(
		"database-recorder-tom",
		config.BROKER_URL,
		config.BROKER_PORT,
		config.BROKER_USERNAME,
		config.BROKER_PASSWORD,
	)

	brokerClient.Subscribe("data/#", func(topic string, message []byte) {
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
