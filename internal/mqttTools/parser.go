package mqttTools

import (
	"errors"
	"regexp"
)

func ParseTopic(topic string) (iata string, measure string, sensorID string, err error) {
	r := regexp.MustCompile(`^data/(?P<IATA>[A-Z]+)/(?P<Measure>Pres|Temp|Wind)/(?P<SensorID>.+)$`)
	matches := r.FindStringSubmatch(topic)
	if len(matches) == 0 {
		err = errors.New("Invalid topic : " + topic)
	} else {
		iata = matches[r.SubexpIndex("IATA")]
		measure = matches[r.SubexpIndex("Measure")]
		sensorID = matches[r.SubexpIndex("SensorID")]
	}
	return
}

func ParseData(payload string) (value string, time string, err error) {
	r := regexp.MustCompile(`(?P<Time>\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\dZ);(?P<Value>\d*.\d*)`)
	matches := r.FindStringSubmatch(payload)
	if len(matches) == 0 {
		err = errors.New("Invalid payload : " + payload)
	} else {
		value = matches[r.SubexpIndex("Value")]
		time = matches[r.SubexpIndex("Time")]
	}
	return
}
