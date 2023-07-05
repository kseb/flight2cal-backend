package models

import (
	"golang.org/x/exp/slices"
)

type ResultFlights struct {
	Size    int      `json:"size"`
	Flights []Flight `json:"flights"`
}

func before(a Flight, b Flight) int {
	if a == b {
		return 0
	}
	if a.Departure.Before(b.Departure) {
		return -1
	}
	return 1
}

func less(a Flight, b Flight) bool {
	return before(a, b) == -1
}

func (resultFlights *ResultFlights) FromFlights(flights Flights) {
	resultFlights.Size = len(flights.Flight)
	for _, flight := range flights.Flight {
		resultFlights.Flights = append(resultFlights.Flights, flight)
	}

	slices.SortFunc(resultFlights.Flights, less)
}
