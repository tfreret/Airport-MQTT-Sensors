package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

type DataRecord struct {
	From        time.Time `json:"beginning time"`
	To          time.Time `json:"ending time"`
	MeasureType string    `json:"type"`
	AirportId   string    `json:"id"`
	Points      []Point   `json:"tab of points"`
}

type Point struct {
	Time     time.Time
	Value    interface{}
	sensorID string
}

var dbClient influxdb2.Client
var dbLock = sync.Mutex{}

func getDBClient() influxdb2.Client {
	if dbClient == nil {
		dbClient = influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"))
	}
	return dbClient
}

// Fonction pour charger les variables d'environnement au démarrage du programme
func init() {
	envFile := filepath.Join("./internal/config", "influxdb.env")
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier influxdb.env :", err)
	}
}

func influxRequest(airportID, sensorID, measureType string, from, to time.Time) (DataRecord, error) {
	InfluxDBClient := influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"))
	defer InfluxDBClient.Close()

	queryAPI := InfluxDBClient.QueryAPI(os.Getenv("INFLUXDB_ORG"))

	var builder strings.Builder
	builder.WriteString("from(bucket:\"" + os.Getenv("INFLUXDB_BUCKET") + "\") ")

	// Layout à respecter pour formatter les dates
	timeLayout := `2006-01-02T15:04:05Z`

	// Test des paramètres
	if to.IsZero() && !from.IsZero() {
		builder.WriteString("|> range(start: " + from.Format(timeLayout) + ") ")
	} else if !to.IsZero() && !from.IsZero() {
		builder.WriteString("|> range(start: " + from.Format(timeLayout) + ", stop: " + to.Format(timeLayout) + ") ")
	} else {
		builder.WriteString("|> range(start: -1h) ")
	}

	builder.WriteString("|> filter(fn: (r) => r._measurement == \"sensor_data\" ")

	appendFilter(&builder, "airport_id", airportID)
	appendFilter(&builder, "sensor_id", sensorID)
	appendFilter(&builder, "sensor_category", measureType)

	builder.WriteString(")")

	query := builder.String()
	fmt.Println(query)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return DataRecord{}, err
	}
	defer result.Close()

	response := DataRecord{}

	// est ce qu'on les met s'il n'y a pas de from/ to ?
	response.From = from
	response.To = to
	response.AirportId = airportID
	response.MeasureType = measureType

	// Gérer le cas où il faut retrouver l'ID depuis la response
	for result.Next() {
		point := Point{
			Time:     result.Record().Time(),
			Value:    result.Record().Value(),
			sensorID: sensorID,
		}
		response.Points = append(response.Points, point)
	}

	return response, nil
}

func appendFilter(builder *strings.Builder, field, value string) {
	if value != "" {
		builder.WriteString("and r." + field + " == \"" + value + "\" ")
	}
}

func getAirports(w http.ResponseWriter, r *http.Request) {
	dbLock.Lock()
	defer dbLock.Unlock()

	client := getDBClient()
	queryAPI := client.QueryAPI(os.Getenv("INFLUXDB_ORG"))

	var builder strings.Builder
	builder.WriteString("from(bucket:\"" + os.Getenv("INFLUXDB_BUCKET") + "\") ")

	builder.WriteString(
		"|> range(start: 0)\n" +
			"|> group()\n" +
			"|> distinct(column: \"airport_id\")\n" +
			"|> keep(columns: [\"_value\"])")

	query := builder.String()

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		handleError(w, err, "Erreur lors de la récupération des aéroports", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	var response []string
	// Gérer le cas où il faut retrouver l'ID depuis la response
	for result.Next() {
		response = append(response, fmt.Sprint(result.Record().Value()))
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		handleError(w, err, "Erreur lors du formatage de la réponse en JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		handleError(w, err, "Erreur dans l'écriture de la réponse", http.StatusInternalServerError)
		return
	}
}

// TODO changer les erreurs err.Error en internal server error pour sécurité
func dataFromSensorCatAirportIDSensorIDHandler(w http.ResponseWriter, r *http.Request) {
	// On récupère les variables de chemin
	vars := mux.Vars(r)
	airportID := vars["airportID"]
	sensorID := vars["sensorID"]
	sensorCat := vars["sensorCat"]

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	parsedFrom, parsedTo, err := checkDates(from, to)
	if err != nil {
		handleError(w, err, "Erreur lors de la vérification des dates", http.StatusInternalServerError)
		return
	}

	// on appelle la BD
	response, err := influxRequest(airportID, sensorID, sensorCat, parsedFrom, parsedTo)
	if err != nil {
		handleError(w, err, "Erreur lors de la requête à la base de données", http.StatusInternalServerError)
		return
	}

	// on formatte la réponse
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		handleError(w, err, "Erreur lors du formatage de la réponse en JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		handleError(w, err, "Erreur dans l'écriture de la réponse", http.StatusInternalServerError)
		return
	}
}

func handleError(w http.ResponseWriter, err error, message string, status int) {
	log.Println(message, ":", err)
	http.Error(w, message, status)
}

func checkDates(from, to string) (time.Time, time.Time, error) {
	layout := "2006-01-02T15:04:05Z"

	parsedFrom, errFrom := time.Parse(layout, from)
	if errFrom != nil && from != "" {
		return time.Time{}, time.Time{}, fmt.Errorf("erreur lors de la conversion de la date from en time.Time: %w", errFrom)
	}

	parsedTo, errTo := time.Parse(layout, to)
	if errTo != nil && to != "" {
		return time.Time{}, time.Time{}, fmt.Errorf("erreur lors de la conversion de la date to en time.Time: %w", errTo)
	}

	if !parsedTo.IsZero() && parsedTo.Before(parsedFrom) {
		return time.Time{}, time.Time{}, errors.New("date de fin avant la date de début")
	}

	//si date de fin après time.Now() -> la query request time.Now(), mais qu'est ce que je renvoie dans la DataStruct ?

	return parsedFrom, parsedTo, nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{sensorCat}/{airportID}/{sensorID}", dataFromSensorCatAirportIDSensorIDHandler).Methods("GET")
	r.HandleFunc("/airports", getAirports).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
		return
	}
}
