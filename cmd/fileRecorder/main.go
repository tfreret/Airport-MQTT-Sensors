package main

import (
	"airport/internal/config"
	"airport/internal/mqttTools"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	configMQTTFile := flag.String("config", "file-recorder.yaml", "Config file of the file recorder")
	outputDirectory := flag.String("outputdir", "./outputs", "Output directory where to store the files")
	flag.Parse()

	configMQTT := config.ReadConfig[mqttTools.MonoConfigMqtt](*configMQTTFile)

	log.Println("Using config : ", configMQTT, outputDirectory)

	err := os.MkdirAll(*outputDirectory, os.ModePerm)
	if err != nil {
		log.Errorf("Couldn't find or create directory '%s' : %s", *outputDirectory, err)
		os.Exit(1)
	}
	brokerClient := mqttTools.NewBrokerClient(
		configMQTT.Mqtt,
	)

	brokerClient.Subscribe("data/#", func(topic string, message []byte) {
		saveMessage(topic, message, *outputDirectory)
	}, configMQTT.Mqtt.MqttQOS)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal
}

func saveMessage(topic string, payload []byte, outputDir string) {
	log.Printf("Received message on topic : '%s' value : '%s'", topic, string(payload))

	iata, measure, _, err := mqttTools.ParseTopic(topic)
	if err != nil {
		log.Errorf("Couldn't extract IATA code; measure, and sensorId type from string : " + topic)
		return
	}

	value, time, err := mqttTools.ParseData(string(payload))
	if err != nil {
		log.Errorf("Couldn't extract value and time from payload : " + string(payload))
	}
	save(iata, measure, value, time, outputDir)
}

func save(iata, measure, value, time, outputDir string) {
	date := strings.Split(time, "T")[0]
	file := fmt.Sprintf("%s/%s-%s.csv", outputDir, iata, date)
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Couldn't open file '%s' : %s\n", file, err)
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("Couldn't close properly file '%s' : %s\n", file, err)
		}
	}(f)

	line := fmt.Sprintf("%s;%s;%s\n", time, measure, value)
	if _, err := f.WriteString(line); err != nil {
		log.Errorf("Couldn't log data to file  '%s' : %s\n", file, err)
		return
	}
}
