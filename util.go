package main

import (
	"math/rand"
	"time"
)

func getTestBackends() []server {
	irelandServerOne := &backend{
		address: "irelandEast.com",
		latency: 1,
		proxy:   nil,
	}
	irelandServerTwo := &backend{
		address: "irelandWest.com",
		latency: 3,
		proxy:   nil,
	}
	londonServer := &backend{
		address: "london.com",
		latency: 5,
		proxy:   nil,
	}
	parisServer := &backend{
		address: "paris.com",
		latency: 7,
		proxy:   nil,
	}
	return []server{irelandServerOne, irelandServerTwo, londonServer, parisServer}
}

func getTestRegions() []*region {
	servers := getTestBackends()
	ireland := CreateRegion([]server{servers[0], servers[1]}, 23, "eu-west-1")
	london := CreateRegion([]server{servers[2]}, 67, "eu-west-2")
	paris := CreateRegion([]server{servers[3]}, 55, "eu-west-3")
	return []*region{ireland, london, paris}
}

func generateRandomInt(n int) int {
	const floor int = 10
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	return r.Intn(n-floor) + floor
}
