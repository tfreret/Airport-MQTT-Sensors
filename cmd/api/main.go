package main

import (
	"airport/internal/config"
	"airport/internal/mqttTools"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mrz1836/go-sanitize"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var log = logrus.New()

var ConfigInflux mqttTools.ConfigInfluxDB

// DataRecord
// @Summary Represents a data record
// @Description Represents a data record with information about time, measurement type, airport ID, and points.
// @ID dataRecord
// @Accept json
// @Produce json
// @Example ExampleDataRecord
type DataRecord struct {
	// beginning time
	// Example: "2022-01-01T00:00:00Z"
	From time.Time
	// ending time
	// Example: "2022-01-02T00:00:00Z"
	To time.Time
	// type of measurement
	// Example: "temperature"
	MeasureType string
	// ID of the airport
	// Example: "JFK"
	AirportId string
	// array of points
	// Example: [{"Time": "2022-01-01T12:00:00Z", "Value": 25.5, "SensorID": "123"}]
	Points []Point
}

// Point
// @Summary Represents a data point
// @Description Represents a data point with time, value, and sensor ID.
// @ID point
// @Accept json
// @Produce json
// @Example ExamplePoint
type Point struct {
	// Time of the data point
	// Example: "2022-01-01T12:00:00Z"
	Time time.Time
	// Value of the data point
	// Example: 25.5
	Value interface{}
	// Sensor ID
	// Example: "123"
	SensorID string
}

// Sensor
// @Summary Represents a sensor
// @Description Represents a sensor with ID and measurement type.
// @ID sensor
// @Accept json
// @Produce json
// @Example ExampleSensor
type Sensor struct {
	// Sensor ID
	// Example: "123"
	ID string
	// Sensor category or measurement type
	// Example: "temperature"
	MeasureType string
}

// AverageResponse
// @Summary Represents the response containing average value
// @Description Represents the response containing the average value for a sensor.
// @ID averageResponse
// @Accept json
// @Produce json
// @Example ExampleAverageResponse
type AverageResponse struct {
	// Average value
	// Example: 25.5
	Average float64
}

// AverageMultipleResponse
// @Summary Represents the response containing multiple averages
// @Description Represents the response containing average values for temperature, pressure, and wind.
// @ID averageMultipleResponse
// @Accept json
// @Produce json
// @Example ExampleAverageMultipleResponse
type AverageMultipleResponse struct {
	// Average value for temperature
	// Example: 25.5
	TempAverage float64
	// Average value for pressure
	// Example: 1013.2
	PresAverage float64
	// Average value for wind speed
	// Example: 10.2
	WindAverage float64
}

var dbClient influxdb2.Client
var dbLock = sync.Mutex{}

func getDBClient() influxdb2.Client {
	if dbClient == nil {
		dbClient = influxdb2.NewClient(ConfigInflux.InfluxDBURL, ConfigInflux.InfluxDBToken)
	}
	return dbClient
}

func influxRequest(airportID, sensorID, measureType string, from, to time.Time) (DataRecord, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	InfluxDBClient := getDBClient()

	queryAPI := InfluxDBClient.QueryAPI(ConfigInflux.InfluxDBOrg)

	var builder strings.Builder
	builder.WriteString("from(bucket:\"" + ConfigInflux.InfluxDBBucket + "\") ")

	// Layout à respecter pour formatter les dates
	timeLayout := `2006-01-02T15:04:05Z`

	// Test des paramètres
	if to.IsZero() && !from.IsZero() {
		builder.WriteString("|> range(start: " + from.Format(timeLayout) + ") ")
	} else if !to.IsZero() && !from.IsZero() {
		builder.WriteString("|> range(start: " + from.Format(timeLayout) + ", stop: " + to.Format(timeLayout) + ") ")
	} else {
		builder.WriteString("|> range(start: -1d) ")
	}

	builder.WriteString("|> filter(fn: (r) => r._measurement == \"sensor_data\" ")

	appendFilter(&builder, "airport_id", airportID)
	if sensorID != "" {
		appendFilter(&builder, "sensor_id", sensorID)
	}
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
			SensorID: sensorID,
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

func writeJson(w *http.ResponseWriter, response any) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		handleError(*w, err, "Erreur lors du formatage de la réponse en JSON", http.StatusInternalServerError)
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	_, err = (*w).Write(jsonResponse)
	if err != nil {
		handleError(*w, err, "Erreur dans l'écriture de la réponse", http.StatusInternalServerError)
		return
	}
}

// @Summary Get a list of sensors for a specific airport.
// @Description This endpoint retrieves a list of sensors for a specific airport.
// @Tags Sensor
// @ID getSensors
// @Produce json
// @Param airportID path string true "ID of the airport"
// @Success 200 {array} Sensor
// @Failure 500 {string} string
// @Router /sensors/{airportID} [get]
func getSensors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	airportID := vars["airportID"]

	dbLock.Lock()
	defer dbLock.Unlock()

	client := getDBClient()
	queryAPI := client.QueryAPI(ConfigInflux.InfluxDBOrg)

	var builder strings.Builder
	builder.WriteString("from(bucket:\"" + ConfigInflux.InfluxDBBucket + "\") ")
	builder.WriteString(
		"|> range(start: 0)\n" +
			"|> filter(fn: (r) => r.airport_id == \"" + sanitize.Alpha(airportID, false) + "\")\n" +
			"|> group()\n" +
			"|> unique(column: \"sensor_id\")\n" +
			"|> keep(columns: [\"sensor_id\", \"sensor_category\"])")

	query := builder.String()

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		handleError(w, err, "Erreur lors de la récupération des capteurs", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	var response = make([]Sensor, 0)
	for result.Next() {
		sensor := Sensor{
			ID:          fmt.Sprint(result.Record().ValueByKey("sensor_id")),
			MeasureType: fmt.Sprint(result.Record().ValueByKey("sensor_category")),
		}
		response = append(response, sensor)
	}

	writeJson(&w, response)
}

// @Summary Get a list of airports.
// @Description This endpoint retrieves a list of airports.
// @ID getAirports
// @Tags Airport
// @Produce json
// @Success 200 {array} string
// @Failure 500 {string} string
// @Router /airports [get]
func getAirports(w http.ResponseWriter, _ *http.Request) {
	dbLock.Lock()
	defer dbLock.Unlock()

	client := getDBClient()
	queryAPI := client.QueryAPI(ConfigInflux.InfluxDBOrg)

	var builder strings.Builder
	builder.WriteString("from(bucket:\"" + ConfigInflux.InfluxDBBucket + "\") ")

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
	for result.Next() {
		response = append(response, fmt.Sprint(result.Record().Value()))
	}

	writeJson(&w, response)
}

// @Summary Get data from a specific sensor for a given sensor type and airport.
// @Description This endpoint retrieves data from a specific sensor based on the sensor type and airport ID.
// @Tags Data
// @ID getDataFromSensorTypeAirportIDSensorID
// @Produce json
// @Param sensorType path string true "Type of sensor"
// @Param airportID path string true "ID of the airport"
// @Param sensorID path string true "ID of the sensor"
// @Param from query string false "Start date (format: 2006-01-02T15:04:05Z)"
// @Param to query string false "End date (format: 2006-01-02T15:04:05Z)"
// @Success 200 {object} DataRecord
// @Failure 500 {string} string
// @Router /data/{airportID}/{sensorType}/{sensorID} [get]
func getDataFromSensorTypeAirportIDSensorID(w http.ResponseWriter, r *http.Request) {
	// On récupère les variables de chemin
	vars := mux.Vars(r)
	airportID := vars["airportID"]
	sensorID := vars["sensorID"]
	sensorType := vars["sensorType"]

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	parsedFrom, parsedTo, err := checkDates(from, to)
	if err != nil {
		handleError(w, err, "Erreur lors de la vérification des dates", http.StatusInternalServerError)
		return
	}

	// on appelle la BD
	response, err := influxRequest(airportID, sensorID, sensorType, parsedFrom, parsedTo)
	if err != nil {
		handleError(w, err, "Erreur lors de la requête à la base de données", http.StatusInternalServerError)
		return
	}

	writeJson(&w, response)
}

// @Summary Get the average value for a specific sensor type at a given airport.
// @Description This endpoint calculates and returns the average value for a specific sensor type at a given airport.
// @ID getAverageBySensorType
// @Tags Average
// @Produce json
// @Param sensorType path string true "Type of sensor"
// @Param airportID path string true "ID of the airport"
// @Success 200 {object} AverageResponse
// @Failure 500 {string} string
// @Router /average/{airportID}/{sensorType} [get]
func getAverageBySensorType(w http.ResponseWriter, r *http.Request) {
	// On récupère les variables de chemin
	vars := mux.Vars(r)
	sensorType := vars["sensorType"]
	airportID := vars["airportID"]

	averageValue, err := calculateAverage(airportID, sensorType)
	if err != nil {
		handleError(w, err, "Erreur lors du calcul de la moyenne", http.StatusInternalServerError)
	}

	result := AverageResponse{Average: averageValue}

	writeJson(&w, result)
}

func calculateAverage(airportID string, sensorType string) (float64, error) {
	// on appelle la BD
	response, err := influxRequest(airportID, "", sensorType, time.Time{}, time.Time{})
	if err != nil {
		return 0.0, err
	}

	// on calcule la moyenne
	if len(response.Points) == 0 {
		return 0.0, fmt.Errorf("no points available to calculate average")
	}

	sum := 0.0

	for _, point := range response.Points {
		// Assuming that the Value field is a float64
		value, ok := point.Value.(float64)
		if !ok {
			return 0.0, fmt.Errorf("value is not a valid float64")
		}

		sum += value
	}

	// on arondit au centième supéreur
	average := math.Ceil(sum/float64(len(response.Points))*100) / 100
	return average, nil
}

// @Summary Get the average values for temperature, pressure, and wind at a given airport.
// @Description This endpoint calculates and returns the average values for temperature, pressure, and wind at a given airport.
// @ID getAllAverage
// @Tags Average
// @Produce json
// @Param airportID path string true "ID of the airport"
// @Success 200 {object} AverageMultipleResponse
// @Failure 500 {string} string
// @Router /average/{airportID} [get]
func getAllAverage(w http.ResponseWriter, r *http.Request) {
	// On récupère les variables de chemin
	vars := mux.Vars(r)
	airportID := vars["airportID"]

	// on appelle la BD
	tempAverage, err := calculateAverage(airportID, "Temp")
	if err != nil {
		handleError(w, err, "Erreur lors du calcul de la moyenne pour la température", http.StatusInternalServerError)
		return
	}

	presAverage, err := calculateAverage(airportID, "Pres")
	if err != nil {
		handleError(w, err, "Erreur lors du calcul de la moyenne pour la pression", http.StatusInternalServerError)
		return
	}

	windAverage, err := calculateAverage(airportID, "Wind")
	if err != nil {
		handleError(w, err, "Erreur lors du calcul de la moyenne pour la vitesse du vent", http.StatusInternalServerError)
		return
	}

	result := AverageMultipleResponse{
		TempAverage: tempAverage,
		PresAverage: presAverage,
		WindAverage: windAverage,
	}

	writeJson(&w, result)
}

func handleError(w http.ResponseWriter, err error, message string, status int) {
	log.Println(message, ":", err)
	http.Error(w, message, status)
}

func checkDates(from, to string) (time.Time, time.Time, error) {
	layout := "2006-01-02T15:04:05Z"

	if from == "" || to == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("l'une des dates est vide")
	}

	parsedFrom, errFrom := time.Parse(layout, from)
	if errFrom != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("erreur lors de la conversion de la date from en time.Time: %w", errFrom)
	}

	parsedTo, errTo := time.Parse(layout, to)
	if errTo != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("erreur lors de la conversion de la date to en time.Time: %w", errTo)
	}

	if !parsedTo.IsZero() && parsedTo.Before(parsedFrom) {
		return time.Time{}, time.Time{}, errors.New("date de fin avant la date de début")
	}

	return parsedFrom, parsedTo, nil
}

// @title Airport Data API
// @description This API provides endpoints to retrieve data from airport sensors.
// @version 1.0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	var (
		influxEnvFile = flag.String("influx", "influxdb.env", ".env file for influx db")
	)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	flag.Parse()
	ConfigInflux = config.ReadEnv[mqttTools.ConfigInfluxDB](*influxEnvFile)

	r := mux.NewRouter()

	r.HandleFunc("/data/{airportID}/{sensorType}/{sensorID}", getDataFromSensorTypeAirportIDSensorID).Methods("GET", "OPTIONS")
	r.HandleFunc("/airports", getAirports).Methods("GET", "OPTIONS")
	r.HandleFunc("/sensors/{airportID}", getSensors).Methods("GET", "OPTIONS")
	r.HandleFunc("/average/{airportID}/{sensorType}", getAverageBySensorType).Methods("GET", "OPTIONS")
	r.HandleFunc("/average/{airportID}", getAllAverage).Methods("GET", "OPTIONS")

	r.HandleFunc("/swaggerJson", func(w http.ResponseWriter, r *http.Request) {
		serveSwaggerJSON(w, "./docs/swagger.json")
	}).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swaggerJson"), //The url pointing to API definition
	))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
		return
	}
}

func serveSwaggerJSON(w http.ResponseWriter, swaggerJSONPath string) {
	file, err := os.Open(swaggerJSONPath)
	if err != nil {
		log.Errorf("Error opening Swagger JSON file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, file); err != nil {
		log.Errorf("Error serving Swagger JSON file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
