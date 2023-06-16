package helper

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func EasyjetAirports(page playwright.Page) map[string]string {
	Info.Println("Looking for easyjet airports")
	airports := make(map[string]string)
	url := "https://www.easyjet.com/en/routemap"

	_, err := page.Goto(url)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return airports
	}

	res, err := page.InnerHTML("[data-title='Flights']")
	if err != nil {
		Error.Println("Couldn't find the Flights element,", err)
		return airports
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return airports
	}

	frameSrc, exists := doc.Find("iframe").Attr("src")
	if !exists {
		Error.Println("Couldn't find the iframe element,", err)
		return airports
	}

	_, err = page.Goto(frameSrc)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return airports
	}
	res, err = page.InnerHTML("#acOriginAirport_ddl")
	if err != nil {
		Error.Println("Couldn't find the acOriginAirport_ddl element,", err)
		return airports
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return airports
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		airport := s.Text()

		re := regexp.MustCompile(`(.*) ([A-Z]{3})`)
		airportSymbolMatch := re.FindStringSubmatch(airport)

		if len(airportSymbolMatch) == 3 {
			airports[airportSymbolMatch[2]] = airportSymbolMatch[1]
		}
	})
	Info.Println("Found easyjet airports:", airports)
	return airports
}

func Easyjet(page playwright.Page, fromSymbol, toSymbol, fromDate, toDate string, airports map[string][]string) ([]Flight, bool) {
	Info.Println("Looking for easyjet flights")
	flight := make([]Flight, 0)
	fromAirport := airports[fromSymbol]
	toAirport := airports[toSymbol]
	if !(SliceContains(fromAirport, EasyjetAirline) && SliceContains(toAirport, EasyjetAirline)) {
		Warning.Println("Easyjet doesn't fly between", fromSymbol, "and", toSymbol)
		return flight, false
	}
	url := "https://www.easyjet.com"

	urlQuery := url + "/deeplink?lang=EN&dep=" + fromSymbol + "&dest=" + toSymbol + "&dd=" + fromDate + "&rd=" +
		toDate + "&apax=1&cpax=0&ipax=0&SearchFrom=SearchPod2_/en/&isOneWay=off"
	Info.Println("Opening page", urlQuery)

	_, err := page.Goto(urlQuery)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return flight, false
	}

	err = page.Click("#ensCloseBanner")
	if err != nil {
		Error.Println("Couldn't find the ensCloseBanner element,", err)
		return flight, false
	}

	page.Click(".drawer-button > button")
	if err != nil {
		Error.Println("Couldn't find the drawer-button button element,", err)
		return flight, false
	}

	page.Click(".outbound .flight-grid-slider > div:nth-child(2) .flight-grid-day div ul")
	if err != nil {
		Error.Println("Couldn't find the outbound flight-grid-day element,", err)
		return flight, false
	}

	time.Sleep(time.Second)
	page.Click(".return .flight-grid-slider > div:nth-child(2) .flight-grid-day div ul")
	if err != nil {
		Error.Println("Couldn't find the return flight-grid-day element,", err)
		return flight, false
	}

	res, err := page.InnerHTML("body")
	if err != nil {
		Error.Println("Couldn't find the body element,", err)
		return flight, false
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return flight, false
	}

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
			duration := "-"
			price := s.Find(".price-eur").Text()

			if len(departureTime) > 0 {
				f := Flight{Airline: EasyjetAirline, Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival),
					DepartureTime: strings.TrimSpace(departureTime), ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.TrimSpace(number),
					Duration: strings.TrimSpace(duration), Price: strings.TrimSpace(price)}

				flight = append(flight, f)
			}
		})
	}
	return flight, true
}
