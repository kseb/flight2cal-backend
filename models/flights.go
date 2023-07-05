package models

import "flight2cal-backend/csv"

type Flights struct {
	Flight map[uint32]Flight `json:"flights"`
}

func (flights *Flights) InsertFlight(flightToInsert Flight) {
	if flights.Flight == nil {
		flights.Flight = make(map[uint32]Flight)
	}
	flightToInsert.ArrivalAirportName = csv.Airports()[flightToInsert.ArrIcao].Name
	flightToInsert.DepartureAirportName = csv.Airports()[flightToInsert.DepIcao].Name
	flights.Flight[flightToInsert.HashCode()] = flightToInsert
}
