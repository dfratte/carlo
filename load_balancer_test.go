package main

import (
	"fmt"
	"testing"
)

func Test_GetLeastCarbonIntenseAvailableServer(t *testing.T) {
	regions := getTestRegions()
	lb := NewLoadBalancer(regions)
	server, lir := lb.GetLeastCarbonIntenseAvailableServer()
	fmt.Printf("server:%s, region:%s\n", server.Address(), lir.Name())
	for _, r := range regions {
		if r.CarbonIntensity() < lir.CarbonIntensity() {
			t.Errorf("lowest carbon intensity region mismatch!")
		}
	}
	for _, s := range lir.AvailableBackends() {
		if s.Latency() < server.Latency() {
			t.Errorf("server latency mismatch!")
		}
	}
}
