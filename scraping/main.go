package main

import "fmt"

// go run .

func main() {
	fmt.Println("WizzAirports")
	//wizzTest()
	test()
	//wizzAirports := WizzairAirports()
	//fmt.Println("Wizz")
	//wizzFlights := Wizzair(wizzAirports)

	fmt.Println("LotAirports")
	lotAirports := LotAirports()
	fmt.Println("Lot")
	lotFlights := Lot(lotAirports)

	fmt.Println("RyanairAirports")
	ryanairAirports := RyanairAirports()
	fmt.Println("Ryanair")
	ryanairFlights := Ryanair(ryanairAirports)

	fmt.Println("EasyjetAirports")
	easyjetAirports := EasyjetAirports()
	fmt.Println("Easyjet")
	easyjetFlights := Easyjet(easyjetAirports)

	fmt.Println("Results")
	fmt.Println(lotFlights.ToString())
	fmt.Println(ryanairFlights.ToString())
	fmt.Println(easyjetFlights.ToString())
	//fmt.Println(wizzFlights.ToString())
}
