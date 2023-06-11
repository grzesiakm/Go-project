package helper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func LotAirports(page playwright.Page) map[string]string {
	Info.Println("Looking for lot airports")
	airports := make(map[string]string)
	url := "https://www.lot.com/us/en"

	_, err := page.Goto(url)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return airports
	}

	time.Sleep(time.Second)
	err = page.Click("#onetrust-accept-btn-handler")
	if err != nil {
		Error.Println("Couldn't find the onetrust-accept-btn-handler element,", err)
		return airports
	}

	time.Sleep(time.Second)
	err = page.Click("#airport-select-0 > .airport-select__value")
	if err != nil {
		Error.Println("Couldn't find the airport-select__value element,", err)
		return airports
	}

	res, err := page.InnerHTML(".combobox__list-wrapper")
	if err != nil {
		Error.Println("Couldn't find the combobox__list-wrapper element,", err)
		return airports
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return airports
	}

	doc.Find("lot-option").Each(func(i int, s *goquery.Selection) {
		airportLabel := s.Find(".airport-select__option-label").Text()
		re := regexp.MustCompile(`\s*(.*) \((.*)\)`)
		airportLabelMatch := re.FindStringSubmatch(airportLabel)

		airports[airportLabelMatch[2]] = airportLabelMatch[1]
	})
	Info.Println("Found lot airports:", airports)
	return airports
}

func GetDateString(inputDate string) string {
	formatDate, _ := time.Parse("2006-01-02", inputDate)
	return fmt.Sprintf("%d%02d%d", formatDate.Day(), formatDate.Month(), formatDate.Year())
}

func Lot(page playwright.Page, fromSymbol, toSymbol, fromDate, toDate string, airports map[string][]string) []Flight {
	Info.Println("Looking for lot flights")
	flight := make([]Flight, 0)
	fromAirport := airports[fromSymbol]
	toAirport := airports[toSymbol]
	if !(SliceContains(fromAirport, LotAirline) && SliceContains(toAirport, LotAirline)) {
		Warning.Println("Lot doesn't fly between", fromSymbol, "and", toSymbol)
		return flight
	}
	url := "https://www.lot.com/us/en"

	urlQuery := url + "?departureAirport=" + fromSymbol + "&destinationAirport=" + toSymbol + "&departureDate=" +
		GetDateString(fromDate) + "&class=E&adults=1&returnDate=" + GetDateString(toDate)
	Info.Println("Opening page", urlQuery)

	_, err := page.Goto(urlQuery)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return flight
	}

	time.Sleep(time.Second)
	err = page.Click("#onetrust-accept-btn-handler")
	if err != nil {
		Error.Println("Couldn't find the onetrust-accept-btn-handler element,", err)
		return flight
	}

	time.Sleep(time.Second)
	err = page.Click(".bookerFlight__submit-button")
	if err != nil {
		Error.Println("Couldn't find the bookerFlight__submit-button element,", err)
		return flight
	}

	time.Sleep(time.Second)
	res, err := page.InnerHTML("#availability-content")
	if err != nil {
		Error.Println("Couldn't find the availability-content element,", err)
		return flight
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return flight
	}

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

			f := Flight{Airline: LotAirline, Departure: airports[strings.TrimSpace(departure)][0], Arrival: airports[strings.TrimSpace(arrival)][0],
				DepartureTime: departureTimeMatch[0], ArrivalTime: arrivalTimeMatch[0], Number: strings.Join(strings.Fields(number), ", "),
				Duration: strings.TrimSpace(duration), Price: strings.Join(strings.Fields(price)[:2], " ")}
			flight = append(flight, f)
		}
	})
	return flight
}
