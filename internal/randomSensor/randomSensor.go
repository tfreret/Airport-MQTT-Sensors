package randomSensor

import (
	"math/rand"
)

type NumberGenerator struct {
	lastNumber float64
	lowerLimit float64
	upperLimit float64
}

func NewNumberGenerator(initialValue, lowerLimit, upperLimit float64) *NumberGenerator {
	return &NumberGenerator{
		lastNumber: initialValue,
		lowerLimit: lowerLimit,
		upperLimit: upperLimit,
	}
}

func (ng *NumberGenerator) GenerateRandomNumber() float64 {
	offset := rand.Intn(3) - 1
	newNumber := ng.lastNumber + float64(offset) + rand.Float64()
	if newNumber < ng.lowerLimit {
		newNumber = ng.lowerLimit
	} else if newNumber > ng.upperLimit {
		newNumber = ng.upperLimit
	}
	ng.lastNumber = newNumber
	return newNumber
}
