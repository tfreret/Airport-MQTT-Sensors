package main

import (
	"airport/internal/config"
	"airport/internal/mqttTools"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	var (
		influxEnvFile = flag.String("influx", "influxdb.env", ".env file for influx db")
	)
	var (
		configMQTTFile = flag.String("config", "database-recorder.yaml", "Config file of the database recorder")
	)
	flag.Parse()

	configMQTT := config.ReadConfig[mqttTools.MonoConfigMqtt](*configMQTTFile)
	configInfluxDB := config.ReadEnv[mqttTools.ConfigInfluxDB](*influxEnvFile)

	log.Println("Using config : ", configMQTT, configInfluxDB)

	InfluxDBClient := influxdb2.NewClient(configInfluxDB.InfluxDBURL, configInfluxDB.InfluxDBToken)
	writeAPI := InfluxDBClient.WriteAPI(configInfluxDB.InfluxDBOrg, configInfluxDB.InfluxDBBucket)

	brokerClient := mqttTools.NewBrokerClient(
		configMQTT.Mqtt,
	)

	brokerClient.Subscribe("data/#", func(topic string, message []byte) {
		iata, measure, sensorId, err := mqttTools.ParseTopic(topic)
		if err != nil {
			log.Errorf("Couldn't extract IATA code; measure, and sensorId type from string : " + topic)
			return
		}

		value, date, err := mqttTools.ParseData(string(message))
		if err != nil {
			log.Errorf("Couldn't extract value and time from payload : " + string(message))
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

		log.Println("Insert in inlfuxDB :", tags, valuefloat, timestamp)

		writeAPI.WritePoint(point)
	}, configMQTT.Mqtt.MqttQOS)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal
	InfluxDBClient.Close()
}
