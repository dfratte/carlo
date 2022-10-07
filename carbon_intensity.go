package main

const maxCarbonIntensity = 200

type carbonIntensityService interface {
	CalculateRegionIntensity(i []*region)
}

type carbonIntensityClient struct {
	carbonIntensityServiceURL string
}

func (c *carbonIntensityClient) CalculateRegionIntensity(i []*region) {
	// This should be solved via an external carbon intensity service using net/http e.g.
	// res, err := http.Get(c.carbonIntensityServiceURL)
	// after calling the service maybe keep a data structure mapping regions to their intensities
	// this should be kept updated, as carbon intensity may vary, though not much
	for _, r := range i {
		r.SetCarbonIntensity(c.getRandomCarbonIntensity())
	}
}

// carbon intensity dimension: gCO2eq
func (c *carbonIntensityClient) getRandomCarbonIntensity() int {
	return generateRandomInt(maxCarbonIntensity)
}
