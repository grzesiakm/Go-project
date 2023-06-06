package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func NorwegianAirports(useragent string) map[string]string {
	fmt.Println(useragent)
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()

	url := "https://www.norwegian.com/uk/"
	_, _ = page.Goto(url)
	page.Click("#nas-cookie-consent-accept-all")
	page.Click("#nas-airport-select-dropdown-input-0")
	res, _ := page.InnerHTML("#nas-airport-select-dropdown-results-0")
	browser.Close()
	pw.Stop()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	airports := make(map[string]string)

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		airportElement := s.Find(".nas-airport-select__name")
		airport := airportElement.Text()
		re := regexp.MustCompile(`(.*) \((.*)\)`)
		airportMatch := re.FindStringSubmatch(airport)
		if len(airportMatch) == 3 {
			airports[airportMatch[2]] = airportMatch[1]
		}
	})

	fmt.Println(airports)
	return airports
}

func GetMonthDayDateString(inputDate string) (string, string) {
	formatDate, _ := time.Parse("2006-01-02", inputDate)
	return fmt.Sprintf("%d%02d", formatDate.Year(), formatDate.Month()), fmt.Sprintf("%02d", formatDate.Day())
}

func Norwegian(airports map[string]string, useragent string) Flights {
	fmt.Println(useragent)
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()

	url := "https://www.norwegian.com/uk"
	from := "Aalborg"
	to := "Algarve"
	fromDate := "2023-06-21"
	toDate := "2023-06-26"

	fromSymbol := KeyByValue(airports, from)
	toSymbol := KeyByValue(airports, to)

	fromYearMonth, fromDay := GetMonthDayDateString(fromDate)
	toYearMonth, toDay := GetMonthDayDateString(toDate)

	urlQuery := url + "/ipc/availability/avaday?AdultCount=1&A_City=" + toSymbol + "&D_City=" + fromSymbol + "&D_Month=" + fromYearMonth + "&D_Day=" + fromDay + "&R_Month=" + toYearMonth + "&R_Day=" + toDay + "&IncludeTransit=true&TripType=2"

	fmt.Println(urlQuery)

	_, _ = page.Goto(urlQuery)
	page.Click(".cookie-consent__accept-all-button")
	res, _ := page.InnerHTML(".sectioncontainer")
	browser.Close()
	pw.Stop()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	flight := make([]Flight, 0)

	doc.Find(".rowinfo1").Each(func(i int, s *goquery.Selection) {
		rowinfo2 := s.Next()
		departure := rowinfo2.Find(".depdest").Text()
		arrival := rowinfo2.Find(".arrdest").Text()
		departureTime := s.Find(".depdest").Text()
		arrivalTime := s.Find(".arrdest").Text()
		number, _ := s.Find(".inputselect .standardlowfare input .hidden").Attr("value")
		re := regexp.MustCompile(`(\w[0-9]{5})`)
		numberMatch := re.FindStringSubmatch(number)
		duration := rowinfo2.Find(".duration").Text()
		re = regexp.MustCompile(`Duration: (.*)`)
		durationMatch := re.FindStringSubmatch(duration)
		price := s.Find(".standardlowfare [title='GBP']").Text()

		f := Flight{Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival), DepartureTime: strings.TrimSpace(departureTime),
			ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.Join(numberMatch[:], ", "), Duration: durationMatch[1], Price: strings.TrimSpace(price)}

		flight = append(flight, f)
	})

	var flights Flights
	flights.Flights = flight
	return flights
}
