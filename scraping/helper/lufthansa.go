package helper

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func LufthansaAirports(page playwright.Page) map[string]string {
	Info.Println("Looking for lufthansa airports")
	airports := make(map[string]string)
	url := "https://www.flightconnections.com/route-map-lufthansa-lh"

	_, err := page.Goto(url)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return airports
	}

	err = page.Click(".qc-cmp2-summary-buttons [mode='primary']")
	if err != nil {
		Error.Println("Couldn't find the qc-cmp2-summary-buttons element,", err)
		return airports
	}

	res, err := page.InnerHTML(".airline-info")
	if err != nil {
		Error.Println("Couldn't find the airline-info-list element,", err)
		return airports
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return airports
	}

	doc.Find(".airline-destination").Each(func(i int, s *goquery.Selection) {
		airport, _ := s.Attr("data-a")
		re := regexp.MustCompile(`(.*) \((.*)\)`)
		airportSymbolMatch := re.FindStringSubmatch(airport)

		if len(airportSymbolMatch) == 3 {
			airports[airportSymbolMatch[2]] = airportSymbolMatch[1]
		}
	})
	Info.Println("Found lufthansa airports:", airports)
	return airports
}

func Lufthansa(page playwright.Page, fromSymbol, toSymbol, fromDate, toDate string, airports map[string][]string) ([]Flight, bool) {
	Info.Println("Looking for lufthansa flights")
	flight := make([]Flight, 0)
	fromAirport := airports[fromSymbol]
	toAirport := airports[toSymbol]
	if !(SliceContains(fromAirport, LufthansaAirline) && SliceContains(toAirport, LufthansaAirline)) {
		Warning.Println("Lufthansa doesn't fly between", fromSymbol, "and", toSymbol)
		return flight, false
	}
	url := "https://www.lufthansa.com/gb/en"

	urlQuery := url + "/flight-search?OriginCode=" + fromSymbol + "&DestinationCode=" + toSymbol + "&DepartureDate=" +
		fromDate + "T18%3A07%3A58&ReturnDate=" + toDate + "T18%3A07%3A58&Cabin=E&PaxAdults=1"
	Info.Println("Opening page", urlQuery)

	_, err := page.Goto(urlQuery)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return flight, false
	}

	err = page.Click("#cm-acceptAll")
	if err != nil {
		Error.Println("Couldn't find the cm-acceptAll element,", err)
		return flight, false
	}

	err = page.Click(".form-btn-section .btn-primary")
	if err != nil {
		Error.Println("Couldn't find the btn-primary element,", err)
		return flight, false
	}

	res, err := page.InnerHTML(".main-content")
	if err != nil {
		Error.Println("Couldn't find the main-content element,", err)
		return flight, false
	}

	if !strings.Contains(res, "No flights found") {
		Warning.Println("No flights for the input")
		return flight, false
	}

	err = page.Click(".sorting-filtering-area")
	if err != nil {
		Error.Println("Couldn't find the sorting-filtering-area element,", err)
		return flight, false
	}

	res1, err := page.InnerHTML(".mat-accordion")
	if err != nil {
		Error.Println("Couldn't find the mat-accordion element,", err)
		return flight, false
	}

	err = page.Click(".mat-accordion .flight-card-button-section > button:nth-child(1)")
	if err != nil {
		Error.Println("Couldn't find the flight-card-button-section element,", err)
		return flight, false
	}

	err = page.Click(".flight-fares ul > li:nth-child(1) i")
	if err != nil {
		Error.Println("Couldn't find the flight-fares element,", err)
		return flight, false
	}

	err = page.Click((".confirm-fares-button"))
	if err != nil {
		Error.Println("Couldn't find the confirm-fares-button element,", err)
		return flight, false
	}

	err = page.Click(".sorting-filtering-area")
	if err != nil {
		Error.Println("Couldn't find the sorting-filtering-area element,", err)
		return flight, false
	}

	res2, err := page.InnerHTML(".mat-accordion")
	if err != nil {
		Error.Println("Couldn't find the mat-accordion element,", err)
		return flight, false
	}

	resSlice := []string{res1, res2}
	for _, res := range resSlice {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
		if err != nil {
			Error.Println("Couldn't create the goquery Document,", err)
			return flight, false
		}

		doc.Find(".upsell-premium-row-pres-container").Each(func(i int, s *goquery.Selection) {
			departure := s.Find(".bound-departure-airport-code").Text()
			departureTime := s.Find(".bound-departure-datetime").Text()
			arrival := s.Find(".bound-arrival-airport-code").Text()
			arrivalTime := s.Find(".bound-arrival-datetime").Text()
			duration := s.Find(".duration-value").Text()
			price := s.Find(".price-amount").Text()
			re := regexp.MustCompile(`\d*.\d{2}`)
			priceMatch := re.FindStringSubmatch(price)

			f := Flight{Airline: LufthansaAirline, Departure: airports[strings.TrimSpace(departure)][0], Arrival: airports[strings.TrimSpace(arrival)][0],
				DepartureTime: strings.TrimSpace(departureTime), ArrivalTime: strings.TrimSpace(arrivalTime), Number: "-",
				Duration: strings.TrimSpace(duration), Price: priceMatch[0]}

			flight = append(flight, f)
		})
	}
	return flight, true
}
