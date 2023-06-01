package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"log"
	"strings"
	"time"
)

func wizzTest() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Failed to launch playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("Failed to launch Chromium: %v", err)
	}
	defer browser.Close()

	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("Failed to create browser context: %v", err)
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("Failed to create new page: %v", err)
	}

	// Navigate to the Wizz Air search page
	if _, err = page.Goto("https://wizzair.com/en-gb/flights/timetable", playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateNetworkidle}); err != nil {
		log.Fatalf("Failed to navigate to Wizz Air website: %v", err)
	}

	// Wait for the element to become visible and clickable
	_, err = page.WaitForSelector(".flight-select__flight-info", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateVisible})
	if err != nil {
		log.Fatalf("Failed to wait for the element: %v", err)
	}

	// Click on the element
	if err = page.Click(".flight-select__flight-info", playwright.PageClickOptions{}); err != nil {
		log.Fatalf("Failed to click on the element: %v", err)
	}

	// Wait for a certain duration to allow the flight details to load
	time.Sleep(5 * time.Second)

	// Extract the flight details using goquery
	htmlContent, err := page.InnerHTML(".flight-select__flight-info")
	if err != nil {
		log.Fatalf("Failed to get HTML content: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find the cheapest flight details
	var cheapestFlight struct {
		DepartureDate string
		DepartureTime string
		ArrivalDate   string
		ArrivalTime   string
		FlightNumber  string
		Duration      string
		Price         string
	}

	doc.Find(".flight-select__flight-info").Each(func(i int, s *goquery.Selection) {
		departureDate := strings.TrimSpace(s.Find(".flight-select__day").Text())
		departureTime := strings.TrimSpace(s.Find(".flight-select__time").Text())
		arrivalDate := strings.TrimSpace(s.Find(".flight-select__return-day").Text())
		arrivalTime := strings.TrimSpace(s.Find(".flight-select__return-time").Text())
		flightNumber := strings.TrimSpace(s.Find(".flight-select__flight-number").Text())
		duration := strings.TrimSpace(s.Find(".flight-select__duration").Text())
		price := strings.TrimSpace(s.Find(".fare-type__price-value").Text())

		if i == 0 || cheapestFlight.Price > price {
			cheapestFlight.DepartureDate = departureDate
			cheapestFlight.DepartureTime = departureTime
			cheapestFlight.ArrivalDate = arrivalDate
			cheapestFlight.ArrivalTime = arrivalTime
			cheapestFlight.FlightNumber = flightNumber
			cheapestFlight.Duration = duration
			cheapestFlight.Price = price
		}
	})

	// Print the cheapest flight details
	fmt.Println("Cheapest Flight Details:")
	fmt.Println("Departure Date:", cheapestFlight.DepartureDate)
	fmt.Println("Departure Time:", cheapestFlight.DepartureTime)
	fmt.Println("Arrival Date:", cheapestFlight.ArrivalDate)
	fmt.Println("Arrival Time:", cheapestFlight.ArrivalTime)
	fmt.Println("Flight Number:", cheapestFlight.FlightNumber)
	fmt.Println("Duration:", cheapestFlight.Duration)
	fmt.Println("Price:", cheapestFlight.Price)
}
