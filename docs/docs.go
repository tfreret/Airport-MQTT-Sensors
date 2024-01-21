// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/airports": {
            "get": {
                "description": "This endpoint retrieves a list of airports.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Airport"
                ],
                "summary": "Get a list of airports.",
                "operationId": "getAirports",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/average/{sensorType}/{airportID}": {
            "get": {
                "description": "This endpoint calculates and returns the average value for a specific sensor type at a given airport.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Average"
                ],
                "summary": "Get the average value for a specific sensor type at a given airport.",
                "operationId": "getAverageBySensorType",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of sensor",
                        "name": "sensorType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the airport",
                        "name": "airportID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.AverageResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/averages/{airportID}": {
            "get": {
                "description": "This endpoint calculates and returns the average values for temperature, pressure, and wind at a given airport.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Average"
                ],
                "summary": "Get the average values for temperature, pressure, and wind at a given airport.",
                "operationId": "getAllAverages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the airport",
                        "name": "airportID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.AverageMultipleResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/data/{sensorType}/{airportID}/{sensorID}": {
            "get": {
                "description": "This endpoint retrieves data from a specific sensor based on the sensor type and airport ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data"
                ],
                "summary": "Get data from a specific sensor for a given sensor type and airport.",
                "operationId": "getDataFromSensorTypeAirportIDSensorID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of sensor",
                        "name": "sensorType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the airport",
                        "name": "airportID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the sensor",
                        "name": "sensorID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Start date (format: 2006-01-02T15:04:05Z)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date (format: 2006-01-02T15:04:05Z)",
                        "name": "to",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.DataRecord"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/sensors/{airportID}": {
            "get": {
                "description": "This endpoint retrieves a list of sensors for a specific airport.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sensor"
                ],
                "summary": "Get a list of sensors for a specific airport.",
                "operationId": "getSensors",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the airport",
                        "name": "airportID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Sensor"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.AverageMultipleResponse": {
            "description": "Represents the response containing average values for temperature, pressure, and wind.",
            "type": "object",
            "properties": {
                "PresAverage": {
                    "description": "Average value for pressure\nExample: 1013.2",
                    "type": "number"
                },
                "TempAverage": {
                    "description": "Average value for temperature\nExample: 25.5",
                    "type": "number"
                },
                "WindAverage": {
                    "description": "Average value for wind speed\nExample: 10.2",
                    "type": "number"
                }
            }
        },
        "main.AverageResponse": {
            "description": "Represents the response containing the average value for a sensor.",
            "type": "object",
            "properties": {
                "moyenne": {
                    "description": "Average value\nExample: 25.5",
                    "type": "number"
                }
            }
        },
        "main.DataRecord": {
            "description": "Represents a data record with information about time, measurement type, airport ID, and points.",
            "type": "object",
            "properties": {
                "beginning time": {
                    "description": "beginning time\nExample: \"2022-01-01T00:00:00Z\"",
                    "type": "string"
                },
                "ending time": {
                    "description": "ending time\nExample: \"2022-01-02T00:00:00Z\"",
                    "type": "string"
                },
                "id": {
                    "description": "ID of the airport\nExample: \"JFK\"",
                    "type": "string"
                },
                "tab of points": {
                    "description": "array of points\nExample: [{\"Time\": \"2022-01-01T12:00:00Z\", \"Value\": 25.5, \"SensorID\": \"123\"}]",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Point"
                    }
                },
                "type": {
                    "description": "type of measurement\nExample: \"temperature\"",
                    "type": "string"
                }
            }
        },
        "main.Point": {
            "description": "Represents a data point with time, value, and sensor ID.",
            "type": "object",
            "properties": {
                "sensorID": {
                    "description": "Sensor ID\nExample: \"123\"",
                    "type": "string"
                },
                "time": {
                    "description": "Time of the data point\nExample: \"2022-01-01T12:00:00Z\"",
                    "type": "string"
                },
                "value": {
                    "description": "Value of the data point\nExample: 25.5"
                }
            }
        },
        "main.Sensor": {
            "description": "Represents a sensor with ID and measurement type.",
            "type": "object",
            "properties": {
                "id": {
                    "description": "Sensor ID\nExample: \"123\"",
                    "type": "string"
                },
                "measureType": {
                    "description": "Sensor category or measurement type\nExample: \"temperature\"",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Airport Data API",
	Description:      "This API provides endpoints to retrieve data from airport sensors.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
