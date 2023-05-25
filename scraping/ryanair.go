package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func Ryanair() {
	newOpts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
		chromedp.Flag("start-maximized", true)}
	opts := append(chromedp.DefaultExecAllocatorOptions[:], newOpts...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://www.ryanair.com/en/en"
	from := "Szczecin"
	to := "Dublin"
	// fromDate := "21-06-2023"
	// toDate := "26-06-2023"
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Click(".cookie-popup-with-overlay__button", chromedp.ByQuery),

		chromedp.Click("#input-button__departure", chromedp.ByQuery),
		chromedp.SetValue("#input-button__departure", from, chromedp.ByQuery),
		chromedp.Click(".airport-item > span", chromedp.ByQuery),
		chromedp.Click(".destination-tabs__button:nth-child(3)", chromedp.ByQuery),

		chromedp.Click("#input-button__destination", chromedp.ByQuery),
		chromedp.SetValue("#input-button__destination", to, chromedp.ByQuery),
		chromedp.Click("#input-button__destination", chromedp.ByQuery),

		chromedp.Click(".ng-star-inserted > .cities__list-scrollable > div:nth-child(2)", chromedp.ByQuery),
		chromedp.Click(".m-toggle__month--after-selected", chromedp.ByQuery),
		chromedp.Click("[data-id='2023-06-11']", chromedp.ByQuery),
		chromedp.Click("[data-id='2023-06-14']", chromedp.ByQuery),

		chromedp.Click(".passengers__confirm-button", chromedp.ByQuery),

		chromedp.Click(".flight-search-widget__start-search", chromedp.ByQuery),
		chromedp.Click(".recent-search__container", chromedp.ByQuery),
		chromedp.Click(".breadcrumb-container", chromedp.ByQuery),

		chromedp.OuterHTML(".journeys-wrapper", &res, chromedp.ByQuery),
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
	fmt.Println(flights.ToString())
}
