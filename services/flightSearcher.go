package services

import (
	"encoding/json"
	"flight2cal-backend/csv"
	"flight2cal-backend/models"
	"flight2cal-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetFlights(c *gin.Context) {
	utils.AddAccessControlAllowOriginIfSet(c)

	arrivalIcao := c.Param("arrivalIcao")
	departureIcao := c.Param("departureIcao")
	dateParam := c.Param("date")

	url := "https://airlabs.co/api/v9/" +
		"routes?api_key=" + utils.AirlabsToken() + "&" +
		"arr_icao=" + arrivalIcao + "&" +
		"dep_icao=" + departureIcao

	response, err := http.Get(url)

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var responseObject models.Airlabs
	unmarshalErr := json.Unmarshal(responseData, &responseObject)
	if unmarshalErr != nil {
		log.Fatal(err)
		return
	}

	flights := models.Flights{}
	for _, airlabsFlight := range responseObject.AirlabsFlight {
		if airlabsFlight.FlightIcao == "" {
			log.Print("Ignoring flight at " + airlabsFlight.DepTimeUtc + " because iata is empty.")
			continue
		}
		departure, err := time.Parse(time.DateTime, dateParam+" "+airlabsFlight.DepTimeUtc+":00")
		arrival, err := time.Parse(time.DateTime, dateParam+" "+airlabsFlight.ArrTimeUtc+":00")
		if err != nil {
			log.Print("Ignoring flight " + airlabsFlight.FlightIcao + " because its departure time cannot be parsed.")
			continue
		}
		if slices.Contains(airlabsFlight.Days, strings.ToLower(departure.Weekday().String()[0:3])) {
			airlineName, err := GetAirline(airlabsFlight.AirlineIata)
			if err != nil {
				log.Println("No iata code found for airline, ignoring flight.", err)
				continue
			}
			flight := models.Flight{
				ArrIcao:       airlabsFlight.ArrIcao,
				DepIcao:       airlabsFlight.DepIcao,
				FlightIcao:    airlabsFlight.FlightIcao,
				Departure:     departure,
				Arrival:       arrival,
				AirlineIata:   airlabsFlight.AirlineIata,
				AirlineName:   airlineName,
				ArrivalCity:   csv.Airports()[airlabsFlight.DepIcao].City,
				DepartureCity: csv.Airports()[airlabsFlight.ArrIcao].City,
			}
			flights.InsertFlight(flight)
		}
	}

	resultFlights := new(models.ResultFlights)
	resultFlights.FromFlights(flights)
	c.IndentedJSON(http.StatusOK, resultFlights)
}
