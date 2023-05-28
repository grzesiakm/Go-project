package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func EasyjetAirports() map[string]string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:3], NewChromeOpts...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://www.easyjet.com/en/routemap"
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("[data-title='Flights']", &res, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	frameSrc, _ := doc.Find("iframe").Attr("src")

	err = chromedp.Run(ctx,
		chromedp.Navigate(frameSrc),
		chromedp.OuterHTML("#acOriginAirport_ddl", &res, chromedp.ByQuery),
	)

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	airports := make(map[string]string)

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		airport := s.Text()

		re := regexp.MustCompile(`(.*) ([A-Z]{3})`)
		airportSymbolMatch := re.FindStringSubmatch(airport)

		if len(airportSymbolMatch) == 3 {
			airports[airportSymbolMatch[2]] = airportSymbolMatch[1]
		}
	})

	fmt.Println(airports)
	return airports
}

func Easyjet(airports map[string]string) Flights {
	opts := append(chromedp.DefaultExecAllocatorOptions[:3], NewChromeOpts...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://www.easyjet.com"
	from := "Bari"
	to := "Basel"
	fromDate := "2023-07-01"
	toDate := "2023-07-08"

	fromSymbol := KeyByValue(airports, from)
	toSymbol := KeyByValue(airports, to)

	urlQuery := url + "/deeplink?lang=EN&dep=" + fromSymbol + "&dest=" + toSymbol + "&dd=" + fromDate + "&rd=" + toDate + "&apax=1&cpax=0&ipax=0&SearchFrom=SearchPod2_/en/&isOneWay=off"

	fmt.Println(urlQuery)

	var res string
	err := chromedp.Run(ctx,
		network.EnableReportingAPI(true),
		chromedp.Navigate(urlQuery),

		chromedp.Click("#ensCloseBanner", chromedp.ByQuery),
		chromedp.Click(".drawer-button > button", chromedp.ByQuery),
		chromedp.Click(".return .flight-grid-slider > div:nth-child(2) .flight-grid-day div", chromedp.ByQuery),
		chromedp.Sleep(time.Second),
		chromedp.Click(".outbound .flight-grid-slider > div:nth-child(2) .flight-grid-day div", chromedp.ByQuery),
		chromedp.OuterHTML("body", &res, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	flight := make([]Flight, 0)

	if len(doc.Find(".basket-wrapper").Text()) > 0 {
		doc.Find(".funnel-basket-flight").Each(func(i int, s *goquery.Selection) {
			route := s.Find(".route-text").Text()
			re := regexp.MustCompile(`(.*) to (.*)`)
			routeMatch := re.FindStringSubmatch(route)
			departure := routeMatch[1]
			arrival := routeMatch[2]
			departureTime := s.Find("[ej-date='Flight.LocalDepartureTime']").Text()
			arrivalTime := s.Find("[ej-date='Flight.LocalArrivalTime']").Text()
			number := s.Find(".flight-number").Text()
			duration := "placeholder"
			price := s.Find(".price-eur").Text()

			f := Flight{Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival), DepartureTime: strings.TrimSpace(departureTime),
				ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.TrimSpace(number), Duration: strings.TrimSpace(duration), Price: strings.TrimSpace(price)}

			flight = append(flight, f)
		})
	}

	var flights Flights
	flights.Flights = flight
	return flights
}
