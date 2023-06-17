package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	. "main/helper"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
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

func GetStartInfo(browser playwright.Browser) ([]string, map[string][]string, map[string]float64) {
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
		page.Close()
		page, _ = context.NewPage()
		ryanairAirports := RyanairAirports(page)
		page.Close()
		page, _ = context.NewPage()
		easyjetAirports := EasyjetAirports(page)
		page.Close()
		page, _ = context.NewPage()
		norwegianAirports := NorwegianAirports(page)
		page.Close()
		page, _ = context.NewPage()
		lufthansaAirports := LufthansaAirports(page)
		page.Close()
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
	currencies := Currency(browser)
	Info.Println("User agents:", useragents)
	// cleaning
	delete(merged, "ANY")
	delete(merged, "")
	for key, value := range merged {
		if strings.Contains(value[0], "- all") || strings.Contains(value[0], "-All") {
			delete(merged, key)
		}
	}
	Info.Println("Merged airports:", merged)
	return useragents, merged, currencies
}

func GetFlights(browser playwright.Browser, fromSymbol, toSymbol, fromDate, toDate string, useragents []string, airports map[string][]string) (Flights, bool) {
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragents[rand.Intn(len(useragents))])})
	page, _ := context.NewPage()
	lotFlights, lotOk := Lot(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	page.Close()
	page, _ = context.NewPage()
	ryanairFlights, ryanairOk := Ryanair(page, fromSymbol, toSymbol, fromDate, toDate, airports, currencies)
	page.Close()
	page, _ = context.NewPage()
	easyjetFlights, easyjetOk := Easyjet(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	page.Close()
	page, _ = context.NewPage()
	norwegianFlights, norwegianOk := Norwegian(page, fromSymbol, toSymbol, fromDate, toDate, airports)
	page.Close()
	page, _ = context.NewPage()
	lufthansaFlights, lufthansaOk := Lufthansa(page, fromSymbol, toSymbol, fromDate, toDate, airports, currencies)
	page.Close()
	var airlinesFlights Flights
	if lotOk || ryanairOk || easyjetOk || norwegianOk || lufthansaOk {
		flights := append(lotFlights, ryanairFlights...)
		flights = append(flights, easyjetFlights...)
		flights = append(flights, norwegianFlights...)
		flights = append(flights, lufthansaFlights...)
		airlinesFlights.Flights = flights
		return airlinesFlights, true
	}
	return airlinesFlights, false
}

var tmpl *template.Template
var pw *playwright.Playwright
var browser playwright.Browser
var useragents []string
var airports map[string][]string
var currencies map[string]float64

func init() {
	rand.Seed(time.Now().Unix())
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	pw, _ = playwright.Run()
	browser, _ = pw.Firefox.Launch(CustomFirefoxOptions)
	tmpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	useragents, airports, currencies = GetStartInfo(browser)
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
	sortedAirports := SortMap(airports)
	err := tmpl.ExecuteTemplate(writer, "index.gohtml", sortedAirports)
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
	fromDate := request.FormValue("start")
	toDate := request.FormValue("end")
	Info.Println("Looking for flights between", from, to, "date", fromDate, toDate)
	flights, ok := GetFlights(browser, from, to, fromDate, toDate, useragents, airports)
	inputData := fmt.Sprintf("%s - %s, %s - %s", from, to, fromDate, toDate)
	if ok {
		sort.Slice(flights.Flights, func(i, j int) bool {
			return flights.Flights[i].Price < flights.Flights[j].Price
		})
		Info.Println(flights.ToString())
		tmpl.ExecuteTemplate(writer, "search.gohtml", struct {
			Results   Flights
			InputData string
		}{flights, inputData})
	} else {
		tmpl.ExecuteTemplate(writer, "noresults.gohtml", inputData)
	}
}
