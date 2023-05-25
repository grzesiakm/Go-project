package main

import "fmt"

type Flight struct {
	Departure     string `json:"Departure"`
	Arrival       string `json:"Arrival"`
	DepartureTime string `json:"DepartureTime"`
	ArrivalTime   string `json:"ArrivalTime"`
	Number        string `json:"Number"`
	Duration      string `json:"Duration"`
	Price         string `json:"Price"`
}

type Flights struct {
	Flights []Flight
}

func (f Flights) ToString() string {
	res := ""
	for x := 0; x < len(f.Flights); x++ {
		res = res + fmt.Sprintf("#%d - departure: %s %s, arrival: %s %s, number: %s, duration: %s, price: %s\n", x, f.Flights[x].Departure, f.Flights[x].DepartureTime,
			f.Flights[x].Arrival, f.Flights[x].ArrivalTime, f.Flights[x].Number, f.Flights[x].Duration, f.Flights[x].Price)
	}
	return res
}
