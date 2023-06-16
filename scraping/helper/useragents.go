package helper

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func UserAgents(browser playwright.Browser) []string {
	log.Println("Looking for user agents")
	useragents := make([]string, 0)
	context, _ := browser.NewContext()
	page, _ := context.NewPage()
	url := "https://www.useragents.me/"

	_, err := page.Goto(url)
	if err != nil {
		log.Fatal("Couldn't open the page,", err)
		return useragents
	}

	res, err := page.InnerHTML("#most-common-desktop-useragents-json-csv > div:nth-child(1)")
	if err != nil {
		log.Fatal("Couldn't find the most-common-desktop-useragents-json-csv element,", err)
		return useragents
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal("Couldn't create the goquery Document,", err)
		return useragents
	}

	useragentsCsv := doc.Find("textarea").Text()
	re := regexp.MustCompile(`"ua": "(.*?)"`)
	useragentsMatch := re.FindAllStringSubmatch(useragentsCsv, -1)
	for _, v := range useragentsMatch {
		useragents = append(useragents, v[1])
	}
	// get the more popular agents
	useragents = useragents[:12]
	log.Println("Found user agents:", useragents)
	return useragents
}
