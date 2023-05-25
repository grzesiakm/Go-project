package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func Lot() {
	newOpts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
		chromedp.Flag("start-maximized", true),
		// chromedp.Flag("auto-open-devtools-for-tabs", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("accept-lang", "en-GB,en-US")}
	opts := append(chromedp.DefaultExecAllocatorOptions[:], newOpts...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://www.lot.com/us/en"
	from := "Szczecin"
	to := "Tirana"
	fromDate := "06-21-2023"
	toDate := "06-26-2023"
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#onetrust-accept-btn-handler", chromedp.ByQuery),
		chromedp.Sleep(time.Millisecond*20),
		chromedp.DoubleClick("#onetrust-accept-btn-handler", chromedp.ByQuery),
		chromedp.WaitNotVisible("#onetrust-policy-text", chromedp.ByQuery),
		chromedp.Sleep(time.Millisecond*10),

		chromedp.DoubleClick("#airport-select-0 > .airport-select__value", chromedp.ByQuery),
		chromedp.SetValue("#airport-select-0_combobox", from, chromedp.ByQuery),
		chromedp.Click("mark:nth-child(2)", chromedp.ByQuery),

		chromedp.Click("#airport-select-1 > .airport-select__value", chromedp.ByQuery),
		chromedp.SetValue("#airport-select-1_combobox", to, chromedp.ByQuery),
		chromedp.Click("mark:nth-child(2)", chromedp.ByQuery),

		chromedp.Click(".booker-form-field--left", chromedp.ByQuery),
		chromedp.Click("#from-date-input", chromedp.ByQuery),
		chromedp.SetValue("#from-date-input", fromDate, chromedp.ByQuery),

		chromedp.Click("#to-date-input", chromedp.ByQuery),
		chromedp.SetValue("#to-date-input", toDate, chromedp.ByQuery),

		chromedp.Click(".datepicker-popover__footer .mat-button-wrapper", chromedp.ByQuery),

		chromedp.Click(".bookerFlight__submit-button", chromedp.ByQuery),

		chromedp.OuterHTML("#availability-content", &res, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(res)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	flight := make([]Flight, 0)

	doc.Find(".flights-table-panel__flight__content").Each(func(i int, s *goquery.Selection) {
		departure := s.Find(".flights-table-panel__flight__content__info__direction__text--departure").Text()
		// arrival := s.Find("[data-ref='flight-segment.arrival'] .flight-info__city").Text()
		arrival := "placeholder"
		departureTime := s.Find(".flights-table-panel__flight__content__info__direction__text--departure").Text()
		arrivalTime := s.Find(".flights-table-panel__flight__content__info__direction--arrival").Text()
		number := s.Find(".flights-table-panel__flight__content__info__details__number").Text()
		duration := s.Find(".VAB__flight__info__time").Text()
		price := s.Find(".ratePanel__element__wrapper__link__bordered__price").Text()
		if len(departure) > 0 {
			departure := "placeholder"
			re := regexp.MustCompile(`\d{2}:\d{2}`)
			departureTimeMatch := re.FindStringSubmatch(departureTime)
			arrivalTimeMatch := re.FindStringSubmatch(arrivalTime)

			f := Flight{Departure: strings.TrimSpace(departure), Arrival: strings.TrimSpace(arrival), DepartureTime: departureTimeMatch[0],
				ArrivalTime: arrivalTimeMatch[0], Number: strings.Join(strings.Fields(number), ", "), Duration: strings.TrimSpace(duration), Price: strings.Join(strings.Fields(price)[:2], " ")}
			flight = append(flight, f)
		}
	})

	var flights Flights
	flights.Flights = flight
	fmt.Println(flights.ToString())
}
