package main

import (
	"flight2cal-backend/csv"
	"flight2cal-backend/services"
	"flight2cal-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load the env vars: %v", err)
	}

	router := startServer()

	if utils.AirlabsToken() == "" {
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
	router.GET("/flights/:arrivalIcao/:departureIcao/:date", services.GetFlights)
	router.GET("/airports/search/:search", searchAirport)
	router.GET("/airports/all", getAllAirports)

	return router
}

func searchAirport(context *gin.Context) {
	utils.AddAccessControlAllowOriginIfSet(context)
	context.IndentedJSON(http.StatusOK, csv.FindAirport(context.Param("search")))
}

func getAllAirports(context *gin.Context) {
	utils.AddAccessControlAllowOriginIfSet(context)
	context.Header("Cache-Control", "max-age=3600")
	context.IndentedJSON(http.StatusOK, csv.GetAllAirports())
}
