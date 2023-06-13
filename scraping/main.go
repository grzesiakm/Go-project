package main

import (
	"encoding/gob"
	"html/template"
	"log"
	. "main/helper"
	"math/rand"
	"net/http"
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
	log.Println("Starting playwright")
	useragents := make([]string, 0)
	decodeFile, err := os.Open("useragents.gob")
	if err == nil {
		defer decodeFile.Close()
		decoder := gob.NewDecoder(decodeFile)
		decoder.Decode(&useragents)
		Info.Println("Read useragents from file")
	} else {
		useragents = UserAgents(browser)
		encodeFile, err := os.Create("useragents.gob")
		if err != nil {
			Error.Println("Error while writing to useragents.gob", err)
			goto nextPart
		}
		encoder := gob.NewEncoder(encodeFile)
		if err := encoder.Encode(useragents); err != nil {
			Error.Println("Error while encoding useragents", err)
			goto nextPart
		}
		encodeFile.Close()
	}
nextPart:
	merged := make(map[string][]string)
	decodeFile, err = os.Open("airports.gob")
	if err == nil {
		defer decodeFile.Close()
		decoder := gob.NewDecoder(decodeFile)
		decoder.Decode(&merged)
		Info.Println("Read airports from file")
	} else {
		useragent := useragents[rand.Intn(len(useragents))]
		log.Println("Using user agent", useragent)
		context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
		page, _ := context.NewPage()
		lotAirports := LotAirports(page)
		ryanairAirports := RyanairAirports(page)
		easyjetAirports := EasyjetAirports(page)
		norwegianAirports := NorwegianAirports(page)
		lufthansaAirports := LufthansaAirports(page)

		log.Println("Merging airports")
		lr := MergeMaps(lotAirports, ryanairAirports)
		lre := MergeMaps(lr, easyjetAirports)
		lren := MergeMaps(lre, norwegianAirports)
		lrenl := MergeMaps(lren, lufthansaAirports)
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
		encodeFile, err := os.Create("airports.gob")
		if err != nil {
			Error.Println("Error while writing to airports.gob", err)
			goto nextNextPart
		}
		encoder := gob.NewEncoder(encodeFile)
		if err := encoder.Encode(merged); err != nil {
			Error.Println("Error while encoding airports", err)
			goto nextNextPart
		}
		encodeFile.Close()
	}
nextNextPart:
	Info.Println("User agents:", useragents)
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

var tmpl *template.Template
var pw *playwright.Playwright
var browser playwright.Browser
var useragents []string
var airports map[string][]string

func init() {
	rand.Seed(time.Now().Unix())
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	pw, _ = playwright.Run()
	browser, _ = pw.Firefox.Launch(CustomFirefoxOptions)
	tmpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	useragents, airports = GetStartInfo(browser)
}

// go run main.go

func main() {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":8080", mux))
	pw.Stop()
}

func index(writer http.ResponseWriter, _ *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "index.gohtml", nil)
	if err != nil {
		Error.Println(err)
		return
	}
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
	flights := GetFlights(browser, from, to, fromDate, toDate, useragents, airports)
	Info.Println(flights.ToString())
	tmpl.ExecuteTemplate(writer, "search.gohtml", flights)
}
