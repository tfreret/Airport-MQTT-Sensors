package randomSensor

import (
	"math/rand"
)

type NumberGenerator struct {
	lastNumber float64
}

func NewNumberGenerator(initialValue float64) *NumberGenerator {
	return &NumberGenerator{
		lastNumber: initialValue,
	}
}

func (ng *NumberGenerator) GenerateRandomNumber() float64 {
	offset := rand.Intn(3) - 1
	newNumber := ng.lastNumber + float64(offset) + rand.Float64()
	ng.lastNumber = newNumber

	return newNumber
}
