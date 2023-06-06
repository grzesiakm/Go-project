package main

import (
	"fmt"
	"math/rand"
	"time"
)

// go run .

func main() {
	rand.Seed(time.Now().Unix())
	useragents := UserAgents()

	fmt.Println("LotAirports")
	lotAirports := LotAirports(useragents[rand.Intn(len(useragents))])
	fmt.Println("Lot")
	lotFlights := Lot(lotAirports, useragents[rand.Intn(len(useragents))])

	fmt.Println("RyanairAirports")
	ryanairAirports := RyanairAirports(useragents[rand.Intn(len(useragents))])
	fmt.Println("Ryanair")
	ryanairFlights := Ryanair(ryanairAirports, useragents[rand.Intn(len(useragents))])

	fmt.Println("EasyjetAirports")
	easyjetAirports := EasyjetAirports(useragents[rand.Intn(len(useragents))])
	fmt.Println("Easyjet")
	easyjetFlights := Easyjet(easyjetAirports, useragents[rand.Intn(len(useragents))])

	fmt.Println("NorwegianAirports")
	norwegianAirports := NorwegianAirports(useragents[rand.Intn(len(useragents))])
	fmt.Println("Norwegian")
	norwegianFlights := Norwegian(norwegianAirports, useragents[rand.Intn(len(useragents))])

	fmt.Println("LufthansaAirports")
	lufthansaAirports := LufthansaAirports(useragents[rand.Intn(len(useragents))])
	fmt.Println("Lufthansa")
	lufthansaFlights := Lufthansa(lufthansaAirports, useragents[rand.Intn(len(useragents))])

	fmt.Println("Results")
	fmt.Println("Lot")
	fmt.Println(lotFlights.ToString())
	fmt.Println("Ryanair")
	fmt.Println(ryanairFlights.ToString())
	fmt.Println("Easyjet")
	fmt.Println(easyjetFlights.ToString())
	fmt.Println("Norwegian")
	fmt.Println(norwegianFlights.ToString())
	fmt.Println("Lufthansa")
	fmt.Println(lufthansaFlights.ToString())
}
