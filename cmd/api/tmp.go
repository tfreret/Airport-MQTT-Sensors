package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func main() {

	envFile := filepath.Join("../../internal/config", "influxdb.env")

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier influxdb.env :", err)
	}

	InfluxDBClient := influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"))
	defer InfluxDBClient.Close()

	queryAPI := InfluxDBClient.QueryAPI(os.Getenv("INFLUXDB_ORG"))

	result, err := queryAPI.Query(context.Background(), "from(bucket:\""+os.Getenv("INFLUXDB_BUCKET")+"\") |> range(start: -1h)")
	if err == nil {
		for result.Next() {
			fmt.Println(result.Record())
		}
	} else {
		fmt.Println(err)
	}
}
