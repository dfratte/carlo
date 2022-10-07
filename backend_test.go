package main

import (
	"testing"
)

func Test_BackendLatency(t *testing.T) {
	b := CreateBackend("backend.test")
	b.latency = generateRandomInt(maxLatency)
	res := b.Latency()
	if res > maxLatency {
		t.Errorf("Latency=%d, can't be larger than 100", res)
	} else {
		t.Logf("Latency=%d within limits", res)
	}
}
