package main

type region struct {
	availableBackends []server
	carbonIntensity   int
	name              string
}

func CreateRegion(s []server, ci int, n string) *region {
	return &region{
		availableBackends: s,
		carbonIntensity:   ci,
		name:              n,
	}
}

func (r *region) AvailableBackends() []server {
	return r.availableBackends
}

func (r *region) CarbonIntensity() int {
	return r.carbonIntensity
}

func (r *region) SetCarbonIntensity(i int) {
	r.carbonIntensity = i
}

func (r *region) Name() string {
	return r.name
}

func (r *region) SetName(n string) {
	r.name = n
}
