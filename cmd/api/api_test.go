package main

import (
	"airport/internal/mqttTools"
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/magiconair/properties/assert"
	"github.com/ory/dockertest/v3"
	_ "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var db *influxdb2.Client
var testConfig = mqttTools.ConfigInfluxDB{
	InfluxDBURL:      "http://localhost:8085",
	InfluxDBBucket:   "airport",
	InfluxDBToken:    "token",
	InfluxDBUsername: "admin",
	InfluxDBPassword: "admin1234",
	InfluxDBOrg:      "iot",
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "influxdb",
		Name:       "influxdb-test",
		Tag:        "2.0.7",
		Env: []string{
			"DOCKER_INFLUXDB_INIT_MODE=setup",
			"DOCKER_INFLUXDB_INIT_USERNAME=" + testConfig.InfluxDBUsername,
			"DOCKER_INFLUXDB_INIT_PASSWORD=" + testConfig.InfluxDBPassword,
			"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=" + testConfig.InfluxDBToken,
			"DOCKER_INFLUXDB_INIT_BUCKET=" + testConfig.InfluxDBBucket,
			"DOCKER_INFLUXDB_INIT_ORG=" + testConfig.InfluxDBOrg,
		},
		ExposedPorts: []string{"8085/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8086/tcp": {{HostIP: "127.0.0.1", HostPort: "8085/tcp"}},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		client := influxdb2.NewClient(testConfig.InfluxDBURL, testConfig.InfluxDBToken)
		db = &client
		if db == nil {
			return err
		}
		_, err := client.Ping(context.TODO())
		return err
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	InsertSampleData()
	dbClient = *db
	ConfigInflux = testConfig
	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}(pool, resource)
	m.Run()
}

func InsertSampleData() {
	tags := map[string]string{
		"airport_id":      "AAA",
		"sensor_category": "Wind",
		"sensor_id":       "capteur1",
	}
	fields := map[string]interface{}{
		"value": 12,
	}
	timestamp, _ := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")
	writeAPI := (*db).WriteAPI(testConfig.InfluxDBOrg, testConfig.InfluxDBBucket)
	point := influxdb2.NewPoint("sensor_data", tags, fields, timestamp)

	writeAPI.WritePoint(point)
	writeAPI.Flush()

	var builder strings.Builder
	builder.WriteString("from(bucket:\"" + testConfig.InfluxDBBucket + "\") ")
}

func TestGetAirports(t *testing.T) {
	r, _ := http.NewRequest("GET", "/airports", nil)
	w := httptest.NewRecorder()
	getAirports(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, strings.Contains(w.Body.String(), "\"AAA\""), true)
}

/*
func TestAppendFilter(t *testing.T) {
	strBuilder := strings.Builder{}
	appendFilter(&strBuilder, "testField", "0")

	expectedResult := "and r.testField == \"0\""

	fmt.Println(expectedResult)
	fmt.Println(strBuilder.String())
	if expectedResult != strBuilder.String() {
		t.Errorf("la chaîne ne correpond pas")
	}
}*/

func TestCheckDatesEmptyFrom(t *testing.T) {
	from := ""
	to := "2024-02-16T12:00:00Z"

	parsedFrom, parsedTo, _ := checkDates(from, to)
	expectedFrom := time.Time{}
	expectedTo := time.Time{}

	if parsedFrom != expectedFrom {
		t.Errorf("from ne correspond pas")
	}

	if parsedTo != expectedTo {
		t.Errorf("to ne correspond pas")
	}
}

func TestCheckDatesEmptyTo(t *testing.T) {
	from := "2024-02-16T12:00:00Z"
	to := ""

	parsedFrom, parsedTo, _ := checkDates(from, to)
	expectedFrom := time.Time{}
	expectedTo := time.Time{}
	if parsedFrom != expectedFrom {
		t.Errorf("from ne correspond pas")
	}

	if parsedTo != expectedTo {
		t.Errorf("to ne correspond pas")
	}
}

func TestCheckDatesInvalidFormat(t *testing.T) {
	from := "12/01/2024"
	to := "2024-02-16T12:00:00Z"

	parsedFrom, parsedTo, _ := checkDates(from, to)
	expectedFrom := time.Time{}
	expectedTo := time.Time{}

	if parsedFrom != expectedFrom {
		t.Errorf("from ne correspond pas")
	}

	if parsedTo != expectedTo {
		t.Errorf("to ne correspond pas")
	}
}

func TestCheckDatesToBeforeFrom(t *testing.T) {
	from := "2024-01-16T12:00:00Z"
	to := "2024-01-02T12:00:00Z"

	parsedFrom, parsedTo, _ := checkDates(from, to)
	expectedFrom := time.Time{}
	expectedTo := time.Time{}

	if parsedFrom != expectedFrom {
		t.Errorf("from ne correspond pas")
	}

	if parsedTo != expectedTo {
		t.Errorf("to ne correspond pas")
	}
}

func TestCheckDatesSuccess(t *testing.T) {
	from := "2024-01-01T12:00:00Z"
	to := "2024-01-02T12:00:00Z"

	parsedFrom, parsedTo, _ := checkDates(from, to)
	expectedFrom := time.Date(2024, time.January, 01, 12, 0, 0, 0, time.UTC)
	expectedTo := time.Date(2024, time.January, 02, 12, 0, 0, 0, time.UTC)

	if parsedFrom != expectedFrom {
		t.Errorf("from ne correspond pas")
	}

	if parsedTo != expectedTo {
		t.Errorf("to ne correspond pas")
	}
}
