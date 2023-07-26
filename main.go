package main

import (
	"encoding/json"
	"flight2cal-backend/csv"
	. "flight2cal-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var AirlabsToken = os.Getenv("AIRLABS_TOKEN")
var AllowOriginHost = os.Getenv("ALLOW_ORIGIN_HOST")

func main() {
	router := startServer()

	if AirlabsToken == "" {
		log.Fatal("AIRLABS_TOKEN not set")
	}

	const BindAddress = "0.0.0.0:8080"
	err := router.Run(BindAddress)
	if err != nil {
		log.Fatal("Cannot run on " + BindAddress)
	}
}

func startServer() *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	router.GET("/ics/:arrivalIcao/:departureIcao/:date", getIcs)
	router.GET("/airports/search/:search", searchAirport)
	router.GET("/airports/all", getAllAirports)

	return router
}

func searchAirport(context *gin.Context) {
	addAccessControlAllowOriginIfSet(context)
	context.IndentedJSON(http.StatusOK, csv.FindAirport(context.Param("search")))
}

func addAccessControlAllowOriginIfSet(context *gin.Context) {
	if AllowOriginHost != "" {
		context.Header("Access-Control-Allow-Origin", AllowOriginHost)
	}
}

func getAllAirports(context *gin.Context) {
	addAccessControlAllowOriginIfSet(context)
	context.Header("Cache-Control", "max-age=3600")
	context.IndentedJSON(http.StatusOK, csv.GetAllAirports())
}

func getIcs(c *gin.Context) {
	addAccessControlAllowOriginIfSet(c)

	arrivalIcao := c.Param("arrivalIcao")
	departureIcao := c.Param("departureIcao")
	dateParam := c.Param("date")

	url := "https://airlabs.co/api/v9/" +
		"routes?api_key=" + AirlabsToken + "&" +
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

	var responseObject Airlabs
	unmarshalErr := json.Unmarshal(responseData, &responseObject)
	if unmarshalErr != nil {
		log.Fatal(err)
		return
	}

	flights := Flights{}
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
			flight := Flight{
				ArrIcao:    airlabsFlight.ArrIcao,
				DepIcao:    airlabsFlight.DepIcao,
				FlightIcao: airlabsFlight.FlightIcao,
				Departure:  departure,
				Arrival:    arrival,
			}
			flights.InsertFlight(flight)
		}
	}

	resultFlights := new(ResultFlights)
	resultFlights.FromFlights(flights)
	c.IndentedJSON(http.StatusOK, resultFlights)
}
