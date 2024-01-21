package main

import (
	"testing"
	"time"
)

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
