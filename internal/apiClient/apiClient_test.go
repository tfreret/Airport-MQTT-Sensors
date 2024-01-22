package apiClient

import (
	"reflect"
	"testing"
)

func TestTranslateApiResponse(t *testing.T) {
	// Cas de test 1: JSON valide
	jsonData := []byte(`{
		"data": [
			{
				"barometer": {"hpa": 1012.5},
				"temperature": {"celsius": 22.5},
				"wind": {"speed_kph": 15.5}
			}
		]
	}`)
	expectedApiResponse := ApiResponse{
		Data: []struct {
			Barometer struct {
				Hpa float64 `json:"hpa"`
			} `json:"barometer"`
			Temperature struct {
				Celsius float64 `json:"celsius"`
			} `json:"temperature"`
			Wind struct {
				SpeedKph float64 `json:"speed_kph"`
			} `json:"wind"`
		}{
			{
				Barometer: struct {
					Hpa float64 `json:"hpa"`
				}{Hpa: 1012.5},
				Temperature: struct {
					Celsius float64 `json:"celsius"`
				}{Celsius: 22.5},
				Wind: struct {
					SpeedKph float64 `json:"speed_kph"`
				}{SpeedKph: 15.5},
			},
		},
	}

	apiResponse, err := translateApiResponse(jsonData)

	if err != nil {
		t.Errorf("Erreur inattendue : %v", err)
	}

	if !reflect.DeepEqual(apiResponse, expectedApiResponse) {
		t.Errorf("La réponse de l'API n'est pas celle attendue. Attendu %v, obtenu %v", expectedApiResponse, apiResponse)
	}

	// Cas de test 2: JSON invalide
	jsonDataInvalid := []byte(`{"invalid_json":}`)
	_, errInvalid := translateApiResponse(jsonDataInvalid)

	if errInvalid == nil {
		t.Error("Une erreur était attendue pour un JSON invalide, mais aucune erreur n'a été retournée")
	}
}
