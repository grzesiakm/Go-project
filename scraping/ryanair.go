package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func RyanairAirports() map[string]string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], NewChromeOpts...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://www.ryanair.com/us/en"
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Click(".cookie-popup-with-overlay__button", chromedp.ByQuery),

		chromedp.Click("#input-button__destination", chromedp.ByQuery),
		chromedp.OuterHTML(".list__airports-scrollable-container", &res, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	airports := make(map[string]string)

	doc.Find("fsw-airport-item").Each(func(i int, s *goquery.Selection) {
		airportElement := s.Find("[data-ref='airport-item__name']")
		airport := airportElement.Text()
		airportSymbol, _ := airportElement.Attr("data-id")

		airports[airportSymbol] = airport
	})

	fmt.Println(airports)
	return airports
}

func Ryanair(airports map[string]string) Flights {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], NewChromeOpts...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://www.ryanair.com/us/en"
	from := "Szczecin"
	to := "Dublin"
	fromDate := "2023-06-21"
	toDate := "2023-06-26"

	fromSymbol := KeyByValue(airports, from)
	toSymbol := KeyByValue(airports, to)

	urlQuery := url + "/trip/flights/select?adults=1&teens=0&children=0&infants=0&dateOut=" + fromDate + "&dateIn=" + toDate + "&isConnectedFlight=false&isReturn=true&discount=0&promoCode=&originIata=" + fromSymbol + "&destinationIata=" + toSymbol + "&tpAdults=1&tpTeens=0&tpChildren=0&tpInfants=0&tpStartDate=" + fromDate + "&tpEndDate=" + toDate + "&tpDiscount=0&tpPromoCode=&tpOriginIata=" + fromSymbol + "&tpDestinationIata=" + toSymbol

	fmt.Println(urlQuery)

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(urlQuery),

		chromedp.Click(".cookie-popup-with-overlay__button", chromedp.ByQuery),
		chromedp.OuterHTML(".journeys-wrapper", &res, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	flight := make([]Flight, 0)

	doc.Find(".flight-card__header").Each(func(i int, s *goquery.Selection) {
		departure := s.Find("[data-ref='flight-segment.departure'] .flight-info__city").Text()
		arrival := s.Find("[data-ref='flight-segment.arrival'] .flight-info__city").Text()
		departureTime := s.Find("[data-ref='flight-segment.departure'] .flight-info__hour").Text()
		arrivalTime := s.Find("[data-ref='flight-segment.arrival'] .flight-info__hour").Text()
		number := s.Find(".flight-info__middle-block .card-flight-num__content").Text()
		duration := s.Find("[data-ref='flight_duration']").Text()
		price := s.Find(".flight-card-summary__new-value flights-price-simple").Text()

		f := Flight{Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival), DepartureTime: strings.TrimSpace(departureTime),
			ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.TrimSpace(number), Duration: strings.TrimSpace(duration), Price: strings.TrimSpace(price)}

		flight = append(flight, f)
	})

	var flights Flights
	flights.Flights = flight
	return flights
}
