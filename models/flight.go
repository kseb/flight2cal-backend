package models

import (
	"hash/fnv"
	"time"
)

func (flight *Flight) HashCode() uint32 {
	hash := fnv.New32()
	stringToHash := flight.Departure.String() + flight.Arrival.String() + flight.AirlineIata
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
	AirlineIata          string    `json:"airline_iata"`
	AirlineName          string    `json:"airline_name"`
	ArrivalCity          string    `json:"arrival_city"`
	DepartureCity        string    `json:"departure_city"`
	DepartureCountry     string    `json:"departure_country"`
	ArrivalCountry       string    `json:"arrival_country"`
}
