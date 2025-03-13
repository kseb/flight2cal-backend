package services

import (
	"encoding/json"
	"errors"
	"flight2cal-backend/models"
	"flight2cal-backend/utils"
	"io"
	"log"
	"net/http"
	"regexp"
)

var airlineMap = make(map[string]string)

func GetAirline(airlineIcao string) (string, error) {
	match, _ := regexp.MatchString("[A-Z]{2,4}", airlineIcao)

	if !match {
		message := "airlineIcao does not match regular expression \"[A-Z]{2,4}\""
		log.Println(message)
		return "", errors.New(message)
	}

	value, found := airlineMap[airlineIcao]
	if found {
		return value, nil
	}

	url := "https://airlabs.co/api/v9/" +
		"airlines?api_key=" + utils.AirlabsToken() + "&" +
		"icao_code=" + airlineIcao

	response, err := http.Get(url)
	if err != nil {
		return "Error retrieving data for getting airline name.", err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "Error reading response data", err
	}

	err = response.Body.Close()
	if err != nil {
		return "Internal server error", err
	}

	var responseObject models.AirlabsAirlines
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return "Error parsing response", err
	}

	for _, airline := range responseObject.Airlines {
		airlineMap[airlineIcao] = airline.Name
		return airline.Name, nil
	}

	return "", errors.New("no airline found")
}
