package csv

import (
	"github.com/biter777/countries"
	"github.com/gocarina/gocsv"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
)

type AirportCsv struct {
	Icao        string `csv:"ident"`
	Name        string `csv:"name"`
	City        string `csv:"municipality"`
	CountryIso  string `csv:"iso_country"`
	AirportType string `csv:"type"`
	Latitude    string `csv:"latitude_deg"`
	Longitude   string `csv:"longitude_deg"`
}

var airports map[string]Airport

type Airport struct {
	Icao      string
	Name      string
	City      string
	Country   string
	Longitude string
	Latitude  string
}

// Airports returns map from icao airport code to airport
func Airports() map[string]Airport {
	if airports != nil {
		return airports
	}

	// iso code to name
	countriesMap := map[string]string{}
	for _, country := range countries.All() {
		countriesMap[country.Alpha2()] = country.String()
	}
	resp, _ := http.Get(getAirportsUrl())
	var airportCsv []AirportCsv

	err := gocsv.Unmarshal(resp.Body, &airportCsv)
	if err != nil {
		log.Fatal("Error getting airportCsv: " + err.Error())
	}

	airportsMap := map[string]Airport{}
	for _, airport := range airportCsv {
		if airport.AirportType == "medium_airport" || airport.AirportType == "large_airport" {
			airportsMap[airport.Icao] = Airport{
				Icao:      airport.Icao,
				Name:      airport.Name,
				City:      airport.City,
				Country:   countriesMap[airport.CountryIso],
				Longitude: airport.Longitude,
				Latitude:  airport.Latitude,
			}
		}
	}

	airports = airportsMap
	return airportsMap
}

func getAirportsUrl() string {
	airportsCsvUrl := os.Getenv("AIRPORT_CSV_URL")
	if airportsCsvUrl != "" {
		return airportsCsvUrl
	} else {
		return "https://davidmegginson.github.io/ourairports-data/airports.csv"
	}
}

func GetAllAirports() []Airport {
	var result []Airport

	keys := make([]string, 0, len(Airports()))
	for k := range Airports() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		result = append(result, Airports()[key])
	}

	return result
}

func FindAirport(search string) []Airport {
	var searchResult []Airport
	regex, _ := regexp.Compile("(?i).*" + search + ".*")
	for _, airport := range Airports() {
		if regex.MatchString(airport.City + airport.Name) {
			searchResult = append(searchResult, airport)
		}
	}
	return searchResult
}
