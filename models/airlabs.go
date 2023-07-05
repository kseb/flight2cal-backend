package models

type Airlabs struct {
	AirlabsFlight []struct {
		ArrIcao    string   `json:"arr_icao"`
		DepIcao    string   `json:"dep_icao"`
		FlightIcao string   `json:"flight_icao"`
		DepTimeUtc string   `json:"dep_time_utc"`
		ArrTimeUtc string   `json:"arr_time_utc"`
		Days       []string `json:"days"`
	} `json:"response"`
}
