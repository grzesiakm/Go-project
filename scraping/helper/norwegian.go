package helper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func NorwegianAirports(page playwright.Page) map[string]string {
	Info.Println("Looking for norwegian airports")
	airports := make(map[string]string)
	url := "https://www.norwegian.com/uk/"

	_, err := page.Goto(url)
	if err != nil {
		Error.Println("Couldn't open the page,", err)
		return airports
	}

	err = page.Click("#nas-cookie-consent-accept-all")
	if err != nil {
		Error.Println("Couldn't find the nas-cookie-consent-accept-all element,", err)
		return airports
	}

	err = page.Click("#nas-airport-select-dropdown-input-0")
	if err != nil {
		Error.Println("Couldn't find the nas-airport-select-dropdown-input-0 element,", err)
		return airports
	}

	res, err := page.InnerHTML("#nas-airport-select-dropdown-results-0")
	if err != nil {
		Error.Println("Couldn't find the nas-airport-select-dropdown-results-0 element,", err)
		return airports
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return airports
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		airportElement := s.Find(".nas-airport-select__name")
		airport := airportElement.Text()
		re := regexp.MustCompile(`(.*) \((.*)\)`)
		airportMatch := re.FindStringSubmatch(airport)
		if len(airportMatch) == 3 {
			airports[airportMatch[2]] = airportMatch[1]
		}
	})
	Info.Println("Found norwegian airports:", airports)
	return airports
}

func GetMonthDayDateString(inputDate string) (string, string) {
	formatDate, _ := time.Parse("2006-01-02", inputDate)
	return fmt.Sprintf("%d%02d", formatDate.Year(), formatDate.Month()), fmt.Sprintf("%02d", formatDate.Day())
}

func Norwegian(page playwright.Page, fromSymbol, toSymbol, fromDate, toDate string, airports map[string][]string) ([]Flight, bool) {
	Info.Println("Looking for norwegian flights")
	flight := make([]Flight, 0)
	fromAirport := airports[fromSymbol]
	toAirport := airports[toSymbol]
	if !(SliceContains(fromAirport, NorwegianAirline) && SliceContains(toAirport, NorwegianAirline)) {
		Warning.Println("Norwegian doesn't fly between", fromSymbol, "and", toSymbol)
		return flight, false
	}
	url := "https://www.norwegian.com/uk"

	fromYearMonth, fromDay := GetMonthDayDateString(fromDate)
	toYearMonth, toDay := GetMonthDayDateString(toDate)

	urlQuery := url + "/ipc/availability/avaday?AdultCount=1&A_City=" + toSymbol + "&D_City=" +
		fromSymbol + "&D_Month=" + fromYearMonth + "&D_Day=" + fromDay + "&R_Month=" + toYearMonth + "&R_Day=" + toDay + "&IncludeTransit=true&TripType=2"
	Info.Println("Opening page", urlQuery)

	resp, err := page.Goto(urlQuery)
	if err != nil {
		Error.Println("Couldn't open the page,", err, resp)
		return flight, false
	}

	res, err := page.InnerHTML("body")
	if err != nil {
		Error.Println("Couldn't find the body element,", err)
		return flight, false
	}

	if strings.Contains(res, "nas-element-cookie-consent__accept-all-button") {
		err = page.Click("#nas-cookie-consent-accept-all")
		if err != nil {
			Error.Println("Couldn't find the nas-cookie-consent-accept-all element,", err)
			return flight, false
		}
	} else {
		err = page.Click(".cookie-consent__accept-all-button")
		if err != nil {
			Error.Println("Couldn't find the cookie-consent__accept-all-button element,", err)
			return flight, false
		}
	}

	res, err = page.InnerHTML("main")
	if err != nil {
		Error.Println("Couldn't find the main element,", err)
		return flight, false
	}

	if !strings.Contains(res, "return trip") {
		Warning.Println("No flights for the input")
		return flight, false
	}

	res, err = page.InnerHTML(".sectioncontainer")
	if err != nil {
		Error.Println("Couldn't find the sectioncontainer element,", err)
		return flight, false
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		Error.Println("Couldn't create the goquery Document,", err)
		return flight, false
	}

	doc.Find(".rowinfo1").Each(func(i int, s *goquery.Selection) {
		rowinfo2 := s.Next()
		departure := rowinfo2.Find(".depdest").Text()
		arrival := rowinfo2.Find(".arrdest").Text()
		departureTime := s.Find(".depdest").Text()
		arrivalTime := s.Find(".arrdest").Text()
		number, _ := s.Find(".inputselect.standardlowfare [type='radio']").Attr("value")
		re := regexp.MustCompile(`([A-Z]{1,2}[0-9]{3,5})`)
		numberMatch := re.FindStringSubmatch(number)
		duration := rowinfo2.Find(".duration").Text()
		re = regexp.MustCompile(`Duration: (.*)`)
		durationMatch := re.FindStringSubmatch(duration)
		price := s.Find(".standardlowfare [title='GBP']").Text()

		if len(price) > 0 {
			f := Flight{Airline: NorwegianAirline, Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival),
				DepartureTime: strings.TrimSpace(departureTime), ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.Join(numberMatch[:], ", "),
				Duration: GetCommonDurationFormat(durationMatch[1]), Price: ConvertToFloat32(strings.TrimSpace(price))}

			flight = append(flight, f)
		}
	})
	return flight, true
}
