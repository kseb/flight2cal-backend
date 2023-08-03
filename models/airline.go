package models

type AirlabsAirlines struct {
	Airlines []struct {
		Iata string `json:"iata_code"`
		Name string `json:"name"`
	} `json:"response"`
}
