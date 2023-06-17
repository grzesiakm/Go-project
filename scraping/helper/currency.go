package helper

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func Currency(browser playwright.Browser) map[string]float64 {
	log.Println("Looking for currency rates")
	currencies := make(map[string]float64)
	context, _ := browser.NewContext()
	page, _ := context.NewPage()
	url := "https://www.x-rates.com/table/?from=GBP&amount=1"

	_, err := page.Goto(url)
	if err != nil {
		log.Fatal("Couldn't open the page,", err)
		return currencies
	}

	res, err := page.InnerHTML(".tablesorter.ratesTable tbody")
	if err != nil {
		log.Fatal("Couldn't find the ratesTable element,", err)
		return currencies
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal("Couldn't create the goquery Document,", err)
		return currencies
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		rate := s.Text()
		href, _ := s.Attr("href")
		if strings.Contains(href, "to=GBP") {
			re := regexp.MustCompile(`from=(.*)&`)
			currencyMatch := re.FindStringSubmatch(href)
			value, err := strconv.ParseFloat(strings.ReplaceAll(rate, ",", "."), 32)
			if err == nil {
				currencies[currencyMatch[1]] = value
			}
		}
	})
	log.Println("Found currencies:", currencies)
	return currencies
}
