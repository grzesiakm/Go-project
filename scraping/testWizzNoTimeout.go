package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func test() {
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	context, _ := browser.NewContext()
	page, _ := context.NewPage()
	url := "https://wizzair.com/en-gb/flights/timetable"

	if err, _ := page.Goto(url); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	// Wait for the departure airport input field to appear
	if err, _ := page.WaitForSelector("#search-departure-station"); err != nil {
		log.Fatalf("Failed to wait for departure airport input: %v", err)
	}

	// Enter the departure airport
	if err := page.Fill("#search-departure-station", "Krakow"); err != nil {
		log.Fatalf("Failed to enter departure airport: %v", err)
	}

	// Wait for the suggestion list to appear and click the first suggestion
	if err, _ := page.WaitForSelector(".flight-search__panel__flight-selector .flight-search__panel__flight-selector__item:nth-child(1)"); err != nil {
		log.Fatalf("Failed to wait for suggestion list: %v", err)
	}
	if err := page.Click(".flight-search__panel__flight-selector .flight-search__panel__flight-selector__item:nth-child(1)"); err != nil {
		log.Fatalf("Failed to click on the first suggestion: %v", err)
	}

	// Wait for the arrival airport input field to appear
	if err, _ := page.WaitForSelector("#search-arrival-station"); err != nil {
		log.Fatalf("Failed to wait for arrival airport input: %v", err)
	}

	// Enter the arrival airport
	if err := page.Fill("#search-arrival-station", "Budapest"); err != nil {
		log.Fatalf("Failed to enter arrival airport: %v", err)
	}

	// Wait for the suggestion list to appear and click the first suggestion
	if err, _ := page.WaitForSelector(".flight-search__panel__flight-selector .flight-search__panel__flight-selector__item:nth-child(1)"); err != nil {
		log.Fatalf("Failed to wait for suggestion list: %v", err)
	}
	if err := page.Click(".flight-search__panel__flight-selector .flight-search__panel__flight-selector__item:nth-child(1)"); err != nil {
		log.Fatalf("Failed to click on the first suggestion: %v", err)
	}

	// Click the search button
	if err := page.Click(".flight-search__panel__cta .flight-search__panel__cta__button"); err != nil {
		log.Fatalf("Failed to click on the search button: %v", err)
	}

	// Wait for the flight list to load
	if err, _ := page.WaitForSelector(".flight-select__flight-info"); err != nil {
		log.Fatalf("Failed to wait for flight list: %v", err)
	}

	// Get the flight details
	flights := make([]Flight, 0)

	res, _ := page.InnerHTML(".flight-select__flight-info")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	doc.Find(".flight-select__flight-info").Each(func(i int, s *goquery.Selection) {
		departureDate := s.Find(".flight-select__date").Text()
		departureTime := s.Find(".flight-select__time").Text()
		arrivalDate := s.Find(".flight-select__date--return").Text()
		arrivalTime := s.Find(".flight-select__time--return").Text()
		flightNumber := s.Find(".flight-select__flight-number").Text()
		duration := s.Find(".flight-select__duration").Text()
		price := s.Find(".flight-select__price .price").Text()

		flight := Flight{
			Departure:     departureDate,
			DepartureTime: departureTime,
			Arrival:       arrivalDate,
			ArrivalTime:   arrivalTime,
			Number:        flightNumber,
			Duration:      duration,
			Price:         price,
		}
		flights = append(flights, flight)
	})

	// Print the flight details
	for _, flight := range flights {
		fmt.Printf("Departure Date: %s\n", flight.Departure)
		fmt.Printf("Departure Time: %s\n", flight.DepartureTime)
		fmt.Printf("Arrival Date: %s\n", flight.Arrival)
		fmt.Printf("Arrival Time: %s\n", flight.ArrivalTime)
		fmt.Printf("Flight Number: %s\n", flight.Number)
		fmt.Printf("Duration: %s\n", flight.Duration)
		fmt.Printf("Price: %s\n", flight.Price)
		fmt.Println("--------------------")
	}

	// Close the browser
	browser.Close()
	pw.Stop()
}
