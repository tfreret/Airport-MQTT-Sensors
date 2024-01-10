package main

import (
	"airport/internal/config"
	"airport/internal/mqttTools"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

func main() {

	outputDirectory := flag.String("outputdir", "./outputs", "Output directory where to store the files")
	flag.Parse()

	err := os.MkdirAll(*outputDirectory, os.ModePerm)
	if err != nil {
		fmt.Printf("Couldn't find or create directory '%s' : %s", outputDirectory, err)
		os.Exit(1)
	}
	brokerClient := mqttTools.NewBrokerClient(
		"file-recorder",
		config.BROKER_URL,
		config.BROKER_PORT,
		config.BROKER_USERNAME,
		config.BROKER_PASSWORD,
	)

	brokerClient.Subscribe("data/#", func(topic string, message []byte) {
		saveMessage(topic, message, *outputDirectory)
	})

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal
}

func saveMessage(topic string, payload []byte, outputDir string) {
	fmt.Printf("Received message on topic : '%s'\n%s", topic, string(payload))

	iata, measure, err := getIata(topic)
	if err != nil {
		fmt.Println("Couldn't extract IATA code and measure type from string : " + topic)
		return
	}

	value, time, err := getPoint(string(payload))
	if err != nil {
		fmt.Println("Couldn't extract value and time from payload : " + string(payload))
	}
	save(iata, measure, value, time, outputDir)
}

func save(iata, measure, value, time, outputDir string) {
	date := strings.Split(time, "T")[0]
	file := fmt.Sprintf("%s/%s-%s.csv", outputDir, iata, date)
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Couldn't open file '%s' : %s\n", file, err)
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Couldn't close properly file '%s' : %s\n", file, err)
		}
	}(f)

	line := fmt.Sprintf("%s;%s;%s\n", time, measure, value)
	if _, err := f.WriteString(line); err != nil {
		fmt.Printf("Couldn't log data to file  '%s' : %s\n", file, err)
		return
	}
}

func getIata(topic string) (iata string, measure string, err error) {
	r := regexp.MustCompile(`^data/(?P<IATA>[A-Z]*)/(?P<Measure>Pres|Temp|Wind)/.*$`)
	matches := r.FindStringSubmatch(topic)
	if len(matches) == 0 {
		err = errors.New("Invalid topic : " + topic)
	} else {
		iata = matches[r.SubexpIndex("IATA")]
		measure = matches[r.SubexpIndex("Measure")]
	}
	return

}

func getPoint(payload string) (value string, time string, err error) {
	r := regexp.MustCompile(`value:(?P<Value>\d*.\d*)\ntime:(?P<Time>\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\dZ)`)
	matches := r.FindStringSubmatch(payload)
	if len(matches) == 0 {
		err = errors.New("Invalid payload : " + payload)
	} else {
		value = matches[r.SubexpIndex("Value")]
		time = matches[r.SubexpIndex("Time")]
	}
	return
}
