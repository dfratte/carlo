package main

import (
	"testing"
)

func Test_CalculateRegionCarbonIntensity(t *testing.T) {
	regions := getTestRegions()
	cic := &carbonIntensityClient{
		carbonIntensityServiceURL: "carbonintensityapi.org",
	}

	cic.CalculateRegionIntensity(regions)

	for _, r := range regions {
		if r.CarbonIntensity() == 0 || r.CarbonIntensity() > maxCarbonIntensity {
			t.Errorf("carbon intensity not set for region: %q\n", r.Name())
		}
	}
}
