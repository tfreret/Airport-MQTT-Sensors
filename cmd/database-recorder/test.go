package main

import (
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func generateRandomValues() (int, int) {
	min := 10
	max := 20
	return rand.Intn(max-min+1) + min, rand.Intn(max-min+1) + min
}

func main() {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNgRTfONwVMK1f977r_g71oNQ=="
	const bucket = "users_business_events"
	const org = "iot"

	client := influxdb2.NewClient("http://localhost:8086", token)

	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)

	for {
		val1, val2 := generateRandomValues()
		fmt.Printf("Valeur 1: %d, Valeur 2: %d\n", val1, val2)
		writeAPI.WriteRecord(fmt.Sprintf("stat,unit=temperature avg=%d,max=%d", val1, val2))
		writeAPI.Flush()
		time.Sleep(2 * time.Second)
	}

	// always close client at the end
	defer client.Close()
}
