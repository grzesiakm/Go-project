package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

func MergeMaps(m1 map[string]string, m2 map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}

func GetStartInfo(browser playwright.Browser) ([]string, map[string][]string) {
	Info.Println("Starting playwright")
	useragents := UserAgents(browser)

	useragent := useragents[rand.Intn(len(useragents))]
	Info.Println("Using user agent", useragent)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()
	lotAirports := LotAirports(page)
	ryanairAirports := RyanairAirports(page)
	easyjetAirports := EasyjetAirports(page)
	norwegianAirports := NorwegianAirports(page)
	lufthansaAirports := LufthansaAirports(page)

	Info.Println("Merging airports")
	lr := MergeMaps(lotAirports, ryanairAirports)
	lre := MergeMaps(lr, easyjetAirports)
	lren := MergeMaps(lre, norwegianAirports)
	lrenl := MergeMaps(lren, lufthansaAirports)
	merged := make(map[string][]string)
	for k, v := range lrenl {
		merged[k] = []string{v}
	}
	for k := range lotAirports {
		merged[k] = append(merged[k], LotAirline)
	}
	for k := range ryanairAirports {
		merged[k] = append(merged[k], RyanairAirline)
	}
	for k := range easyjetAirports {
		merged[k] = append(merged[k], EasyjetAirline)
	}
	for k := range norwegianAirports {
		merged[k] = append(merged[k], NorwegianAirline)
	}
	for k := range lufthansaAirports {
		merged[k] = append(merged[k], LufthansaAirline)
	}
	Info.Println("Merged airports:", merged)
	return useragents, merged
}

func GetFlights(browser playwright.Browser, from, to, fromDate, toDate string, useragents []string, airports map[string][]string) Flights {
	fromSymbol := KeyByValue(airports, from)
	toSymbol := KeyByValue(airports, to)

	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragents[rand.Intn(len(useragents))])})
	page, _ := context.NewPage()

	lotFlights := Lot(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	ryanairFlights := Ryanair(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	easyjetFlights := Easyjet(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	norwegianFlights := Norwegian(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	lufthansaFlights := Lufthansa(page, fromSymbol, toSymbol, fromDate, toDate, airports)

	flights := append(lotFlights, ryanairFlights...)
	flights = append(flights, easyjetFlights...)
	flights = append(flights, norwegianFlights...)
	flights = append(flights, lufthansaFlights...)
	var airlinesFlights Flights
	airlinesFlights.Flights = flights
	return airlinesFlights
}

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

const (
	EasyjetAirline   string = "Easyjet"
	LotAirline       string = "Lot"
	LufthansaAirline string = "Lufthansa"
	NorwegianAirline string = "Norwegian"
	RyanairAirline   string = "Ryanair"
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// go run .

func main() {
	Info.Println("Initializing the application")
	rand.Seed(time.Now().Unix())
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	useragents, airports := GetStartInfo(browser)

	from := "Aalborg"
	to := "Zagreb"
	fromDate := "2023-06-21"
	toDate := "2023-06-26"
	Info.Println("Looking for flights from", from, "to", to, "date", fromDate, "to", toDate)

	flights := GetFlights(browser, from, to, fromDate, toDate, useragents, airports)
	Info.Println("Found flights from", from, "to", to, "date", fromDate, "to", toDate, "\n", flights.ToString())
	pw.Stop()
}
