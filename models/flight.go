package models

import (
	"hash/fnv"
	"time"
)

func (flight *Flight) HashCode() uint32 {
	hash := fnv.New32()
	stringToHash := flight.Departure.String() + flight.Arrival.String()
	_, _ = hash.Write([]byte(stringToHash))
	return hash.Sum32()
}

type Flight struct {
	ArrIcao              string    `json:"arr_icao"`
	DepIcao              string    `json:"dep_icao"`
	FlightIcao           string    `json:"flight_icao"`
	Departure            time.Time `json:"departure"`
	Arrival              time.Time `json:"arrival"`
	DepartureAirportName string    `json:"departure_airport_name"`
	ArrivalAirportName   string    `json:"arrival_airport_name"`
}
