package helper

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/playwright-community/playwright-go"
)

const (
	EasyjetAirline   string = "Easyjet"
	LotAirline       string = "Lot"
	LufthansaAirline string = "Lufthansa"
	NorwegianAirline string = "Norwegian"
	RyanairAirline   string = "Ryanair"
)

type Flight struct {
	Airline       string  `json:"Airline"`
	Departure     string  `json:"Departure"`
	Arrival       string  `json:"Arrival"`
	DepartureTime string  `json:"DepartureTime"`
	ArrivalTime   string  `json:"ArrivalTime"`
	Number        string  `json:"Number"`
	Duration      string  `json:"Duration"`
	Price         float32 `json:"Price"`
}

type Flights struct {
	Flights []Flight
}

type KV struct {
	Key   string
	Value string
}

func (f Flights) ToString() string {
	res := ""
	for x := 0; x < len(f.Flights); x++ {
		res = res + fmt.Sprintf("#%d %s - departure: %s %s, arrival: %s %s, number: %s, duration: %s, price: %v\n", x, f.Flights[x].Airline, f.Flights[x].Departure,
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

func SortMap(m map[string][]string) []KV {
	var sortedMap []KV
	for k, v := range m {
		sortedMap = append(sortedMap, KV{k, v[0]})
	}

	sort.Slice(sortedMap, func(i, j int) bool {
		return sortedMap[i].Value < sortedMap[j].Value
	})

	return sortedMap
}

func ConvertToFloat32(s string) float32 {
	if strings.Contains(s, ".") && strings.Contains(s, ",") {
		s = strings.ReplaceAll(s, ",", "")
	}
	value, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 32)
	if err == nil {
		return float32(math.Round(value*100) / 100)
	}
	return 0
}

func GetCommonDurationFormat(s string) string {
	if strings.Contains(s, ":") {
		sli := strings.Split(s, ":")
		for i, r := range sli[1] {
			if !unicode.IsDigit(r) {
				return fmt.Sprintf("%sh %sm", sli[0], sli[1][:i])
			}
		}
	}
	return s
}

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

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
