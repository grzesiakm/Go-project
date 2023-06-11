package main

import (
	"github.com/playwright-community/playwright-go"
	"html/template"
	"log"
	"main/helper"
	"math/rand"
	"net/http"
	"time"
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
	log.Println("Starting playwright")
	useragents := helper.UserAgents(browser)

	useragent := useragents[rand.Intn(len(useragents))]
	log.Println("Using user agent", useragent)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()
	lotAirports := helper.LotAirports(page)
	ryanairAirports := helper.RyanairAirports(page)
	easyjetAirports := helper.EasyjetAirports(page)
	norwegianAirports := helper.NorwegianAirports(page)
	lufthansaAirports := helper.LufthansaAirports(page)

	log.Println("Merging airports")
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
	log.Println("Merged airports:", merged)
	return useragents, merged
}

func GetFlights(browser playwright.Browser, from, to, fromDate, toDate string, useragents []string, airports map[string][]string) helper.Flights {
	fromSymbol := helper.KeyByValue(airports, from)
	toSymbol := helper.KeyByValue(airports, to)

	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragents[rand.Intn(len(useragents))])})
	page, _ := context.NewPage()

	lotFlights := helper.Lot(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	ryanairFlights := helper.Ryanair(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	easyjetFlights := helper.Easyjet(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	norwegianFlights := helper.Norwegian(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	lufthansaFlights := helper.Lufthansa(page, fromSymbol, toSymbol, fromDate, toDate, airports)

	flights := append(lotFlights, ryanairFlights...)
	flights = append(flights, easyjetFlights...)
	flights = append(flights, norwegianFlights...)
	flights = append(flights, lufthansaFlights...)
	var airlinesFlights helper.Flights
	airlinesFlights.Flights = flights
	return airlinesFlights
}

const (
	EasyjetAirline   string = "Easyjet"
	LotAirline       string = "Lot"
	LufthansaAirline string = "Lufthansa"
	NorwegianAirline string = "Norwegian"
	RyanairAirline   string = "Ryanair"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// go run main.go

func main() {
	//Info.Println("Initializing the application")
	//rand.Seed(time.Now().Unix())
	//pw, _ := playwright.Run()
	//browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	//useragents, airports := GetStartInfo(browser)
	//
	//from := "Aalborg"
	//to := "Zagreb"
	//fromDate := "2023-06-21"
	//toDate := "2023-06-26"
	////Info.Println("Looking for flights from", from, "to", to, "date", fromDate, "to", toDate)
	//
	//flights := GetFlights(browser, from, to, fromDate, toDate, useragents, airports)
	//Info.Println("Found flights from", from, "to", to, "date", fromDate, "to", toDate, "\n", flights.ToString())
	//pw.Stop()

	var mux = http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func index(writer http.ResponseWriter, _ *http.Request) {
	tmpl.ExecuteTemplate(writer, "index.gohtml", nil)
}

func search(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	from := request.FormValue("departure")
	to := request.FormValue("arrival")
	fromDate := request.FormValue("departureTime")
	toDate := request.FormValue("arrivalTime")

	rand.Seed(time.Now().Unix())
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(helper.CustomFirefoxOptions)
	useragents, airports := GetStartInfo(browser)

	flights := GetFlights(browser, from, to, fromDate, toDate, useragents, airports)

	tmpl.ExecuteTemplate(writer, "search.gohtml", flights)
}
