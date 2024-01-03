package apiClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApiResponse struct {
	Data []struct {
		Barometer struct {
			Hpa float64 `json:"hpa"`
		} `json:"barometer"`
		Temperature struct {
			Celsius float64 `json:"celsius"`
		} `json:"temperature"`
		Wind struct {
			SpeedKph float64 `json:"speed_kph"`
		} `json:"wind"`
	} `json:"data"`
}

func doHttpRequest(method string, url string, header string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la requête HTTP : %w", err)
	}

	req.Header.Add("X-API-Key", header)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête HTTP : %w", err)
	}

	return res, nil
}

func getBody(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du corps de la réponse : %w", err)
	}
	return body, nil
}

func translateApiResponse(body []byte) (ApiResponse, error) {
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return ApiResponse{}, fmt.Errorf("erreur lors du désérialisation de la réponse de l'API : %w", err)
	}
	return apiResponse, nil
}

func GetApiResponse(apiUrl string, apiKey string) (ApiResponse, error) {
	method := "GET"

	res, err := doHttpRequest(method, apiUrl, apiKey)
	if err != nil {
		return ApiResponse{}, err
	}
	defer res.Body.Close()

	body, err := getBody(res)
	if err != nil {
		return ApiResponse{}, err
	}

	apiResponse, err := translateApiResponse(body)
	if err != nil {
		return ApiResponse{}, err
	}

	return apiResponse, nil
}
