package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type Flight struct {
	Airline       string `json:"Airline"`
	Departure     string `json:"Departure"`
	Arrival       string `json:"Arrival"`
	DepartureTime string `json:"DepartureTime"`
	ArrivalTime   string `json:"ArrivalTime"`
	Number        string `json:"Number"`
	Duration      string `json:"Duration"`
	Price         string `json:"Price"`
}

type Flights struct {
	Flights []Flight
}

func (f Flights) ToString() string {
	res := ""
	for x := 0; x < len(f.Flights); x++ {
		res = res + fmt.Sprintf("#%d %s - departure: %s %s, arrival: %s %s, number: %s, duration: %s, price: %s\n", x, f.Flights[x].Airline, f.Flights[x].Departure,
			f.Flights[x].DepartureTime, f.Flights[x].Arrival, f.Flights[x].ArrivalTime, f.Flights[x].Number, f.Flights[x].Duration, f.Flights[x].Price)
	}
	return res
}

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func KeyByValue(m map[string][]string, value string) string {
	for k, v := range m {
		for _, vals := range v {
			if strings.Contains(vals, value) {
				return k
			}
		}
	}
	return ""
}

var RemoveAnimationCss = `
	* {
		transition-duration: 0s !important;
	}`

// animation-delay: -0.0001s !important;
// animation-duration: 0s !important;
// animation-play-state: paused !important;
// caret-color: transparent !important;

var AddCssScript = `
	(css) => {
		const style = document.createElement('style');
		style.type = 'text/css';
		style.appendChild(document.createTextNode(css));
		document.head.appendChild(style);
		return true;
	}`
var RemoveElementScript = `
	(id) => {
		const element = document.getElementById(id);
		element.remove();
		return true;
	}`
var HideElementScript = `
	(id) => {
		const element = document.getElementById(id);
		element.style.display = 'none';
		element.style.visibility = 'hidden';
		return true;
	}`

func SaveToFile(filename, content string) {

	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

var CustomFirefoxOptions = playwright.BrowserTypeLaunchOptions{
	Headless: playwright.Bool(true),
	FirefoxUserPrefs: map[string]interface{}{"security.insecure_field_warning.contextual.enabled": false,
		"security.certerrors.permanentOverride":       false,
		"network.stricttransportsecurity.preloadlist": false,
		"security.enterprise_roots.enabled":           true},
}
