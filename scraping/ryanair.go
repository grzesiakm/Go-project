package main

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func RyanairAirports(page playwright.Page) map[string]string {
	Info.Println("Looking for ryanair airports")
	airports := make(map[string]string)
	url := "https://www.ryanair.com/us/en"

	_, err := page.Goto(url)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return airports
	}

	err = page.Click(".cookie-popup-with-overlay__button")
	if err != nil {
		Error.Println("Couldn't find the cookie-popup-with-overlay__button element,", err)
		return airports
	}

	err = page.Click("#input-button__destination")
	if err != nil {
		Error.Println("Couldn't find the input-button__destination element,", err)
		return airports
	}

	res, err := page.InnerHTML(".list__airports-scrollable-container")
	if err != nil {
		Error.Println("Couldn't find the list__airports-scrollable-container element,", err)
		return airports
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return airports
	}

	doc.Find("fsw-airport-item").Each(func(i int, s *goquery.Selection) {
		airportElement := s.Find("[data-ref='airport-item__name']")
		airport := airportElement.Text()
		re := regexp.MustCompile(`^(\s+)|(\s+)$`)
		airportMatch := re.ReplaceAllLiteralString(airport, "")
		airportSymbol, _ := airportElement.Attr("data-id")

		airports[airportSymbol] = airportMatch
	})
	Info.Println("Found ryanair airports:", airports)
	return airports
}

func Ryanair(page playwright.Page, fromSymbol, toSymbol, fromDate, toDate string, airports map[string][]string) []Flight {
	Info.Println("Looking for ryanair flights")
	flight := make([]Flight, 0)
	fromAirport := airports[fromSymbol]
	toAirport := airports[toSymbol]
	if !(SliceContains(fromAirport, RyanairAirline) && SliceContains(toAirport, RyanairAirline)) {
		Warning.Println("Ryanair doesn't fly between", fromSymbol, "and", toSymbol)
		return flight
	}
	url := "https://www.ryanair.com/us/en"

	urlQuery := url + "/trip/flights/select?adults=1&teens=0&children=0&infants=0&dateOut=" +
		fromDate + "&dateIn=" + toDate + "&isConnectedFlight=false&isReturn=true&discount=0&promoCode=&originIata=" +
		fromSymbol + "&destinationIata=" + toSymbol + "&tpAdults=1&tpTeens=0&tpChildren=0&tpInfants=0&tpStartDate=" +
		fromDate + "&tpEndDate=" + toDate + "&tpDiscount=0&tpPromoCode=&tpOriginIata=" + fromSymbol + "&tpDestinationIata=" + toSymbol
	Info.Println("Opening page", urlQuery)

	_, err := page.Goto(urlQuery)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return flight
	}

	err = page.Click(".cookie-popup-with-overlay__button")
	if err != nil {
		Error.Println("Couldn't find the cookie-popup-with-overlay__button element,", err)
		return flight
	}

	res, err := page.InnerHTML(".journeys-wrapper")
	if err != nil {
		Error.Println("Couldn't find the journeys-wrapper element,", err)
		return flight
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return flight
	}

	doc.Find(".flight-card__header").Each(func(i int, s *goquery.Selection) {
		departure := s.Find("[data-ref='flight-segment.departure'] .flight-info__city").Text()
		arrival := s.Find("[data-ref='flight-segment.arrival'] .flight-info__city").Text()
		departureTime := s.Find("[data-ref='flight-segment.departure'] .flight-info__hour").Text()
		arrivalTime := s.Find("[data-ref='flight-segment.arrival'] .flight-info__hour").Text()
		number := s.Find(".flight-info__middle-block .card-flight-num__content").Text()
		duration := s.Find("[data-ref='flight_duration']").Text()
		price := s.Find(".flight-card-summary__new-value flights-price-simple").Text()

		f := Flight{Airline: RyanairAirline, Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival),
			DepartureTime: strings.TrimSpace(departureTime), ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.TrimSpace(number),
			Duration: strings.TrimSpace(duration), Price: strings.TrimSpace(price)}

		flight = append(flight, f)
	})
	return flight
}
