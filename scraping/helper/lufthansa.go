package helper

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func LufthansaAirports(page playwright.Page) map[string]string {
	log.Println("Looking for lufthansa airports")
	airports := make(map[string]string)
	url := "https://www.lufthansa.com/us/en/flights"

	_, err := page.Goto(url)
	if err != nil {
		log.Fatal("Couldn't open the page,", err)
		return airports
	}

	err = page.Click("#cm-acceptAll")
	if err != nil {
		log.Fatal("Couldn't find the cm-acceptAll element,", err)
		return airports
	}

	err = page.Click("[placeholder='From']")
	if err != nil {
		log.Fatal("Couldn't find the From element,", err)
		return airports
	}

	err = page.Click(".autocomplete-airport .input-icon")
	if err != nil {
		log.Fatal("Couldn't find the input-icon element,", err)
		return airports
	}

	err = page.Click(".df-result-wrapper .btn-secondary")
	if err != nil {
		log.Fatal("Couldn't find the btn-secondary element,", err)
		return airports
	}

	var timeout = float64(1000)
	for {
		err := page.Click(".df-result-wrapper .btn-secondary", playwright.PageClickOptions{Timeout: &timeout})
		if err != nil {
			goto nextPart
		}
	}
nextPart:
	res, err := page.InnerHTML(".df-result-section > ol")
	if err != nil {
		log.Fatal("Couldn't find the df-result-section element,", err)
		return airports
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal("Couldn't create the goquery Document,", err)
		return airports
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		airport := s.Find(".city-name").Text()
		airportSymbol, _ := s.Find(".image-wrapper img").Attr("src")
		re := regexp.MustCompile(`destination\/(.*)-square`)
		airportSymbolMatch := re.FindStringSubmatch(airportSymbol)

		airports[strings.ToUpper(airportSymbolMatch[1])] = strings.TrimSpace(airport)
	})
	log.Println("Found lufthansa airports:", airports)
	return airports
}

func Lufthansa(page playwright.Page, fromSymbol, toSymbol, fromDate, toDate string, airports map[string][]string) []Flight {
	log.Println("Looking for lufthansa flights")
	flight := make([]Flight, 0)
	fromAirport := airports[fromSymbol]
	toAirport := airports[toSymbol]
	if !(SliceContains(fromAirport, LufthansaAirline) && SliceContains(toAirport, LufthansaAirline)) {
		log.Println("Lufthansa doesn't fly between", fromSymbol, "and", toSymbol)
		return flight
	}
	url := "https://www.lufthansa.com/us/en"

	urlQuery := url + "/flight-search?OriginCode=" + fromSymbol + "&DestinationCode=" + toSymbol + "&DepartureDate=" +
		fromDate + "T18%3A07%3A58&ReturnDate=" + toDate + "T18%3A07%3A58&Cabin=E&PaxAdults=1"
	log.Println("Opening page", urlQuery)

	_, err := page.Goto(urlQuery)
	if err != nil {
		log.Fatal("Couldn't open the page,", err)
		return flight
	}

	err = page.Click("#cm-acceptAll")
	if err != nil {
		log.Fatal("Couldn't find the cm-acceptAll element,", err)
		return flight
	}

	err = page.Click(".form-btn-section .btn-primary")
	if err != nil {
		log.Fatal("Couldn't find the btn-primary element,", err)
		return flight
	}

	err = page.Click(".sorting-filtering-area")
	if err != nil {
		log.Fatal("Couldn't find the sorting-filtering-area element,", err)
		return flight
	}

	res1, err := page.InnerHTML(".mat-accordion")
	if err != nil {
		log.Fatal("Couldn't find the mat-accordion element,", err)
		return flight
	}

	err = page.Click(".mat-accordion .flight-card-button-section > button:nth-child(1)")
	if err != nil {
		log.Fatal("Couldn't find the flight-card-button-section element,", err)
		return flight
	}

	err = page.Click(".flight-fares ul > li:nth-child(1) i")
	if err != nil {
		log.Fatal("Couldn't find the flight-fares element,", err)
		return flight
	}

	err = page.Click((".confirm-fares-button"))
	if err != nil {
		log.Fatal("Couldn't find the confirm-fares-button element,", err)
		return flight
	}

	err = page.Click(".sorting-filtering-area")
	if err != nil {
		log.Fatal("Couldn't find the sorting-filtering-area element,", err)
		return flight
	}

	res2, err := page.InnerHTML(".mat-accordion")
	if err != nil {
		log.Fatal("Couldn't find the mat-accordion element,", err)
		return flight
	}

	resSlice := []string{res1, res2}
	for _, res := range resSlice {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
		if err != nil {
			log.Fatal("Couldn't create the goquery Document,", err)
			return flight
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

			f := Flight{Airline: LufthansaAirline, Departure: airports[strings.TrimSpace(departure)][0], Arrival: airports[strings.TrimSpace(arrival)][0],
				DepartureTime: strings.TrimSpace(departureTime), ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.TrimSpace("none"),
				Duration: strings.TrimSpace(duration), Price: priceMatch[0]}

			flight = append(flight, f)
		})
	}
	return flight
}
