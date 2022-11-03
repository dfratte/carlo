package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const maxCarbonIntensity = 200

type CarbonIntensityResponse struct {
	Location        string    `json:"location"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	CarbonIntensity float64   `json:"carbonIntensity"`
}

type carbonIntensityService interface {
	CalculateRegionIntensity(i []*region)
}

type carbonIntensityClient struct {
	carbonIntensityServiceURL string
}

func (c *carbonIntensityClient) CalculateRegionIntensity(i []*region) {
	for _, r := range i {
		r.SetCarbonIntensity(c.getCarbonIntensity(r))
	}
}

func (c *carbonIntensityClient) getCarbonIntensity(r *region) int {
	requestURL := "https://carbon-aware-api.azurewebsites.net/emissions/average-carbon-intensity"

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	current := time.Now().AddDate(0, 0, -1)
	aFewDaysAgo := current.AddDate(0, 0, -7)

	q := req.URL.Query()
	q.Add("location", r.name)
	q.Add("startTime", aFewDaysAgo.Format(time.RFC3339))
	q.Add("endTime", current.Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	res, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var result CarbonIntensityResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return int(result.CarbonIntensity)
}
