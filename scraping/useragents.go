package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func UserAgents() []string {
	pw, _ := playwright.Run()
	browser, _ := pw.Firefox.Launch(CustomFirefoxOptions)
	context, _ := browser.NewContext()
	page, _ := context.NewPage()

	url := "https://www.useragents.me/"
	_, _ = page.Goto(url)
	res, _ := page.InnerHTML("#most-common-desktop-useragents-json-csv > div:nth-child(1)")
	browser.Close()
	pw.Stop()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		log.Fatal(err)
	}

	useragents := make([]string, 0)
	useragentsCsv := doc.Find("textarea").Text()
	re := regexp.MustCompile(`"ua": "(.*?)"`)
	useragentsMatch := re.FindAllStringSubmatch(useragentsCsv, -1)
	for _, v := range useragentsMatch {
		useragents = append(useragents, v[1])
	}
	// get the more popular agents
	useragents = useragents[:12]
	fmt.Println(useragents)
	return useragents
}
