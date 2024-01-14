package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

type DataRecord struct {
	Time   time.Time          `json:"time"`
	Tags   map[string]string  `json:"tags"`
	Fields map[string]float64 `json:"fields"`
}

// Fonction pour charger les variables d'environnement au démarrage du programme
func init() {
	envFile := filepath.Join("../internal/config", "influxdb.env")
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier influxdb.env :", err)
	}
}

func influxRequest(airportID, sensorID, sensorCat string) (map[time.Time]interface{}, error) {
	InfluxDBClient := influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"))
	defer InfluxDBClient.Close()

	queryAPI := InfluxDBClient.QueryAPI(os.Getenv("INFLUXDB_ORG"))

	query := fmt.Sprintf(`
		from(bucket:"%s")
		|> range(start: -1h)
		|> filter(fn: (r) => r._measurement == "sensor_data" and r.airport_id == "%s" and r.sensor_id == "%s" and r.sensor_category == "%s")
	`, os.Getenv("INFLUXDB_BUCKET"), airportID, sensorID, sensorCat)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	response := make(map[time.Time]interface{})
	for result.Next() {
		response[result.Record().Time()] = result.Record().Value()
	}

	return response, nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// On récupère les variables de chemin
	vars := mux.Vars(r)
	airportID := vars["airportID"]
	sensorID := vars["sensorID"]
	sensorCat := vars["sensorCat"]

	// on appelle la BD
	response, err := influxRequest(airportID, sensorID, sensorCat)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// on formatte la réponse
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{sensorCat}/{airportID}/{sensorID}", testHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}
