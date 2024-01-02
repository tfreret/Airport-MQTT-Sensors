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

func doHttpRequest(method string, url string, header string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Add("X-API-Key", header)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res
}

func getBody(res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return body
}

func translateApiResponse(body []byte) (ApiResponse, error) {
	var apiResponse ApiResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return ApiResponse{}, err
	}

	return apiResponse, nil
}

func GetApiResponse(apiUrl string, apiKey string) (ApiResponse, error) {
	method := "GET"

	res := doHttpRequest(method, apiUrl, apiKey)
	defer res.Body.Close()

	body := getBody(res)

	apiResponse, err := translateApiResponse(body)
	if err != nil {
		return ApiResponse{}, err
	}

	return apiResponse, nil
}
