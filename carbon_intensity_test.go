package main

import (
	"fmt"
	"testing"
)

type carbonIntensityClientMock struct {
	carbonIntensityServiceURL string
}

func (c *carbonIntensityClientMock) CalculateRegionIntensity(i []*region) {
	for _, r := range i {
		r.SetCarbonIntensity(c.getCarbonIntensity(r))
	}
}

func (c *carbonIntensityClientMock) getCarbonIntensity(r *region) int {
	fmt.Printf("calculating carbon intensity for region %s\n", r.name)
	return generateRandomInt(maxCarbonIntensity)
}

func Test_CalculateRegionCarbonIntensity(t *testing.T) {
	regions := getTestRegions()
	cic := &carbonIntensityClientMock{
		carbonIntensityServiceURL: "carbonintensityapi.org",
	}

	cic.CalculateRegionIntensity(regions)

	for _, r := range regions {
		if r.CarbonIntensity() == 0 || r.CarbonIntensity() > maxCarbonIntensity {
			t.Errorf("carbon intensity not set for region: %q\n", r.Name())
		}
	}
}
