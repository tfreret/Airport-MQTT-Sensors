package mqttTools

type ConfigInfluxDB struct {
	InfluxDBUsername string `mapstructure:"INFLUXDB_USERNAME"`
	InfluxDBPassword string `mapstructure:"INFLUXDB_PASSWORD"`
	InfluxDBToken    string `mapstructure:"INFLUXDB_TOKEN"`
	InfluxDBURL      string `mapstructure:"INFLUXDB_URL"`
	InfluxDBOrg      string `mapstructure:"INFLUXDB_ORG"`
	InfluxDBBucket   string `mapstructure:"INFLUXDB_BUCKET"`
}
