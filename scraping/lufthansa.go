package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func LufthansaAirports(useragent string) map[string]string {
	fmt.Println(useragent)
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()

	url := "https://www.lufthansa.com/us/en/flights"
	_, _ = page.Goto(url)
	page.Click("#cm-acceptAll")
	page.Click("[placeholder='From']")
	page.Click(".autocomplete-airport .input-icon")
	page.Click(".df-result-wrapper .btn-secondary")

	var timeout = float64(1000)
	for {
		err := page.Click(".df-result-wrapper .btn-secondary", playwright.PageClickOptions{Timeout: &timeout})
		if err != nil {
			goto nextPart
		}
	}
nextPart:
	res, _ := page.InnerHTML(".df-result-section > ol")
	browser.Close()
	pw.Stop()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	airports := make(map[string]string)

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		airport := s.Find(".city-name").Text()
		airportSymbol, _ := s.Find(".image-wrapper img").Attr("src")
		re := regexp.MustCompile(`destination\/(.*)-square`)
		airportSymbolMatch := re.FindStringSubmatch(airportSymbol)

		airports[strings.ToUpper(airportSymbolMatch[1])] = strings.TrimSpace(airport)
	})

	fmt.Println(airports)
	return airports
}

func Lufthansa(airports map[string]string, useragent string) Flights {
	fmt.Println(useragent)
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()

	url := "https://www.lufthansa.com/us/en"
	from := "Frankfurt"
	to := "Aarhus"
	fromDate := "2023-06-07"
	toDate := "2023-06-20"

	fromSymbol := KeyByValue(airports, from)
	toSymbol := KeyByValue(airports, to)

	urlQuery := url + "/flight-search?OriginCode=" + fromSymbol + "&DestinationCode=" + toSymbol + "&DepartureDate=" + fromDate + "T18%3A07%3A58&ReturnDate=" + toDate + "T18%3A07%3A58&Cabin=E&PaxAdults=1"
	fmt.Println(urlQuery)

	_, _ = page.Goto(urlQuery)
	page.Click("#cm-acceptAll")
	page.Click(".form-btn-section .btn-primary")
	page.Click(".sorting-filtering-area")
	res1, _ := page.InnerHTML(".mat-accordion")
	page.Click(".mat-accordion .flight-card-button-section > button:nth-child(1)")
	page.Click(".flight-fares ul > li:nth-child(1) i")
	page.Click((".confirm-fares-button"))
	page.Click(".sorting-filtering-area")
	res2, _ := page.InnerHTML(".mat-accordion")
	browser.Close()
	pw.Stop()
	flight := make([]Flight, 0)
	resSlice := []string{res1, res2}
	for _, res := range resSlice {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".upsell-premium-row-pres-container").Each(func(i int, s *goquery.Selection) {
			departure := s.Find(".bound-departure-airport-code").Text()
			departureTime := s.Find(".bound-departure-datetime").Text()
			arrival := s.Find(".bound-arrival-airport-code").Text()
			arrivalTime := s.Find(".bound-arrival-datetime").Text()
			// number := s.Find(".flight-select__flight-number").Text()
			duration := s.Find(".duration-value").Text()
			price := s.Find(".price-amount").Text()
			re := regexp.MustCompile(`\d*.\d{2}`)
			priceMatch := re.FindStringSubmatch(price)

			f := Flight{Departure: airports[strings.TrimSpace(departure)], Arrival: airports[strings.TrimSpace(arrival)], DepartureTime: strings.TrimSpace(departureTime),
				ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.TrimSpace("none"), Duration: strings.TrimSpace(duration), Price: priceMatch[0]}

			flight = append(flight, f)
		})
	}

	var flights Flights
	flights.Flights = flight
	return flights
}
