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

func LotAirports(useragent string) map[string]string {
	fmt.Println(useragent)
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()
	url := "https://www.lot.com/us/en"
	_, _ = page.Goto(url)
	time.Sleep(time.Second)
	page.Click("#onetrust-accept-btn-handler")
	time.Sleep(time.Second)
	page.Click("#airport-select-0 > .airport-select__value")
	res, _ := page.InnerHTML(".combobox__list-wrapper")
	browser.Close()
	pw.Stop()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	airports := make(map[string]string)

	doc.Find("lot-option").Each(func(i int, s *goquery.Selection) {
		airportLabel := s.Find(".airport-select__option-label").Text()
		re := regexp.MustCompile(`\s*(.*) \((.*)\)`)
		airportLabelMatch := re.FindStringSubmatch(airportLabel)

		airports[airportLabelMatch[2]] = airportLabelMatch[1]
	})

	fmt.Println(airports)
	return airports
}

func GetDateString(inputDate string) string {
	formatDate, _ := time.Parse("2006-01-02", inputDate)
	return fmt.Sprintf("%d%02d%d", formatDate.Day(), formatDate.Month(), formatDate.Year())
}

func Lot(airports map[string]string, useragent string) Flights {
	fmt.Println(useragent)
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{UserAgent: playwright.String(useragent)})
	page, _ := context.NewPage()
	url := "https://www.lot.com/us/en"
	from := "New York"
	to := "Cairo"
	fromDate := "2023-08-20"
	toDate := "2023-08-23"

	fromSymbol := KeyByValue(airports, from)
	toSymbol := KeyByValue(airports, to)

	urlQuery := url + "?departureAirport=" + fromSymbol + "&destinationAirport=" + toSymbol + "&departureDate=" + GetDateString(fromDate) + "&class=E&adults=1&returnDate=" + GetDateString(toDate)

	fmt.Println(urlQuery)

	_, _ = page.Goto(urlQuery)
	time.Sleep(time.Second)
	page.Click("#onetrust-accept-btn-handler")
	time.Sleep(time.Second)
	page.Click(".bookerFlight__submit-button")
	time.Sleep(time.Second)
	res, _ := page.InnerHTML("#availability-content")
	browser.Close()
	pw.Stop()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	flight := make([]Flight, 0)

	doc.Find(".flights-table-panel__flight__content").Each(func(i int, s *goquery.Selection) {
		departure := s.Find(".flights-table-panel__flight__content__info__direction__text--departure .flights-table-panel__flight__content__info__direction__text__acronym").Text()
		arrival := s.Find(".flights-table-panel__flight__content__info__direction--arrival .flights-table-panel__flight__content__info__direction__text__acronym").Text()
		departureTime := s.Find(".flights-table-panel__flight__content__info__direction__text--departure").Text()
		arrivalTime := s.Find(".flights-table-panel__flight__content__info__direction--arrival").Text()
		number := s.Find(".flights-table-panel__flight__content__info__details__number").Text()
		duration := s.Find(".VAB__flight__info__time").Text()
		price := s.Find(".ratePanel__element__wrapper__link__bordered__price").Text()
		if len(departure) > 0 {
			re := regexp.MustCompile(`\d{2}:\d{2}`)
			departureTimeMatch := re.FindStringSubmatch(departureTime)
			arrivalTimeMatch := re.FindStringSubmatch(arrivalTime)

			f := Flight{Departure: airports[strings.TrimSpace(departure)], Arrival: airports[strings.TrimSpace(arrival)], DepartureTime: departureTimeMatch[0],
				ArrivalTime: arrivalTimeMatch[0], Number: strings.Join(strings.Fields(number), ", "), Duration: strings.TrimSpace(duration), Price: strings.Join(strings.Fields(price)[:2], " ")}
			flight = append(flight, f)
		}
	})

	var flights Flights
	flights.Flights = flight
	return flights
}
