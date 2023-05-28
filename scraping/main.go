package main

import "fmt"

// go run .

func main() {

	fmt.Println("RyanairAirports")
	ryanairAirports := RyanairAirports()
	fmt.Println("Ryanair")
	ryanairFlights := Ryanair(ryanairAirports)
	fmt.Println(ryanairFlights.ToString())

	fmt.Println("EasyjetAirports")
	easyjetAirports := EasyjetAirports()
	fmt.Println("Easyjet")
	easyjetFlights := Easyjet(easyjetAirports)

	fmt.Println("Results")
	fmt.Println(ryanairFlights.ToString())
	fmt.Println(easyjetFlights.ToString())
}
