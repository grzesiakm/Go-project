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
	url := "https://worldwide.easyjet.com/en"

	urlQueryFrom := url + "/search?origins=" + fromSymbol + "&destinations=" + toSymbol + "&departureDate=" + fromDate + "&returnDate=&isOneWay=true&currency=GBP&residency=GB&utm_source=easyjet_search_pod&utm_medium=&utm_campaign=&adult=18&child=&infant="

	urlQueryTo := url + "/search?origins=" + toSymbol + "&destinations=" + fromSymbol + "&departureDate=" + toDate + "&returnDate=&isOneWay=true&currency=GBP&residency=GB&utm_source=easyjet_search_pod&utm_medium=&utm_campaign=&adult=18&child=&infant="
	Info.Println("Opening page", urlQueryFrom, urlQueryTo)

	resp, err := page.Goto(urlQueryFrom)
	if err != nil {
		Error.Println("Couldn't open the page,", err, resp)
		return flight, false
	}

	err = page.Click("[data-testid='uc-accept-all-button']")
	if err != nil {
		Error.Println("Couldn't find the uc-accept-all-button element,", err)
		return flight, false
	}

	res1, err := page.InnerHTML(".css-1k6o6fq")
	if err != nil {
		Error.Println("Couldn't find the css-1k6o6fq element,", err)
		return flight, false
	}

	resp, err = page.Goto(urlQueryTo)
	if err != nil {
		Error.Println("Couldn't open the page,", err, resp)
		return flight, false
	}

	time.Sleep(time.Second)
	res2, err := page.InnerHTML(".css-1k6o6fq")
	if err != nil {
		Error.Println("Couldn't find the css-1k6o6fq element,", err)
		return flight, false
	}
	resSlice := []string{res1, res2}
	for _, res := range resSlice {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
		if err != nil {
			Error.Println("Couldn't create the goquery Document,", err)
			return flight, false
		}

		doc.Find(".css-spi7mf.eoq8clo1").Each(func(i int, s *goquery.Selection) {
			departure := ""
			arrival := ""
			departureTime := ""
			arrivalTime := ""
			s.Find(".css-12y62xy.e40v1ey3").Each(func(j int, s2 *goquery.Selection) {
				airport := s2.Text()
				re := regexp.MustCompile(`(\w*) \(`)
				airportMatch := re.FindStringSubmatch(airport)
				if len(airportMatch) > 0 {
					if j == 0 {
						departure = airportMatch[1]
					} else {
						arrival = airportMatch[1]
					}
				}
			})
			s.Find(".css-1lbgppu.e40v1ey1").Each(func(j int, s2 *goquery.Selection) {
				if j == 0 {
					departureTime = s2.Text()
				} else {
					arrivalTime = s2.Text()
				}
			})
			number := make([]string, 0)
			s.Find(".flightNumber").Each(func(j int, s2 *goquery.Selection) {
				number = append(number, s2.Text())
			})
			duration := s.Find(".css-sx8r9.edc9gsy2").Text()
			re := regexp.MustCompile(`Journey duration: (.*)`)
			durationMatch := re.FindStringSubmatch(duration)
			price := s.Find(".ejzj98p8.css-x70eh0.eevppah0").Text()
			price = strings.ReplaceAll(price, "\u00a0", " ")
			price = strings.ReplaceAll(price, "\u00a3", " ")
			priceSlice := strings.Fields(price)
			if len(priceSlice) > 1 {
				price = priceSlice[1]
			} else if len(priceSlice) == 1 {
				price = priceSlice[0]
			}
			if len(departure) > 0 {
				f := Flight{Airline: EasyjetAirline, Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival),
					DepartureTime: strings.TrimSpace(departureTime), ArrivalTime: strings.TrimSpace(arrivalTime), Number: strings.Join(number, ", "),
					Duration: GetCommonDurationFormat(strings.TrimSpace(durationMatch[1])), Price: ConvertToFloat32(price)}

				flight = append(flight, f)
			}
		})
	}
	return flight, true
}
