package main

import (
	"fmt"
	"net/http"
	"sort"
)

type loadBalancer struct {
	availableRegions  []*region
	carbonAwareClient carbonIntensityService
}

func NewLoadBalancer(regions []*region) *loadBalancer {
	return &loadBalancer{
		availableRegions:  regions,
		carbonAwareClient: &carbonIntensityClient{},
	}
}

func (lb *loadBalancer) HandleRequest(rw http.ResponseWriter, r *http.Request) {
	backend, _ := lb.GetLeastCarbonIntenseAvailableServer()
	backend.Serve(rw, r)
}

func (lb *loadBalancer) GetLeastCarbonIntenseAvailableServer() (server, *region) {
	lb.carbonAwareClient.CalculateRegionIntensity(lb.availableRegions)

	fmt.Printf("\n---carbon intensity for all available regions---\n")
	for _, region := range lb.availableRegions {
		if region != nil {
			fmt.Printf("%s: %d, ", region.Name(), region.CarbonIntensity())
		}
	}

	sort.Slice(lb.availableRegions, func(i, j int) bool {
		return lb.availableRegions[i].CarbonIntensity() < lb.availableRegions[j].CarbonIntensity()
	})

	lowestCarbonIntensityRegion := lb.availableRegions[0]

	fmt.Printf("\nlowest carbon intensity FOUND at %s=%d[gCO2eq]!\n",
		lowestCarbonIntensityRegion.Name(),
		lowestCarbonIntensityRegion.CarbonIntensity(),
	)
	fmt.Printf("\n---latencies for servers in %s region---\n", lowestCarbonIntensityRegion.Name())

	lowestCarbonIntensityBackends := lowestCarbonIntensityRegion.AvailableBackends()
	for _, server := range lowestCarbonIntensityBackends {
		if server != nil {
			fmt.Printf("server %s -> latency %d[ms]\n", server.Address(), server.Latency())
		}
	}

	sort.Slice(lowestCarbonIntensityBackends, func(i, j int) bool {
		return lowestCarbonIntensityBackends[i].Latency() < lowestCarbonIntensityBackends[j].Latency()
	})

	for _, server := range lowestCarbonIntensityBackends {
		if server != nil {
			fmt.Printf("\nsmallest server latency FOUND! | forwarding req to %s\n", server.Address())
			return server, lowestCarbonIntensityRegion
		} else {
			continue
		}
	}

	return nil, nil
}
