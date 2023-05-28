package main

import (
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
)

type Flight struct {
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
		res = res + fmt.Sprintf("#%d - departure: %s %s, arrival: %s %s, number: %s, duration: %s, price: %s\n", x, f.Flights[x].Departure, f.Flights[x].DepartureTime,
			f.Flights[x].Arrival, f.Flights[x].ArrivalTime, f.Flights[x].Number, f.Flights[x].Duration, f.Flights[x].Price)
	}
	return res
}

func KeyByValue(m map[string]string, value string) string {
	for k, v := range m {
		if strings.Contains(v, value) {
			return k
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
var NewChromeOpts = []chromedp.ExecAllocatorOption{

	chromedp.Flag("disable-background-networking", true),
	chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
	chromedp.Flag("disable-background-timer-throttling", true),
	chromedp.Flag("disable-backgrounding-occluded-windows", true),
	chromedp.Flag("disable-breakpad", true),
	chromedp.Flag("disable-client-side-phishing-detection", true),
	chromedp.Flag("disable-default-apps", true),
	chromedp.Flag("disable-dev-shm-usage", true),
	chromedp.Flag("disable-extensions", true),
	chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees,UserAgentClientHint"),
	chromedp.Flag("disable-hang-monitor", true),
	chromedp.Flag("disable-ipc-flooding-protection", true),
	chromedp.Flag("disable-popup-blocking", true),
	chromedp.Flag("disable-prompt-on-repost", true),
	chromedp.Flag("disable-renderer-backgrounding", true),
	chromedp.Flag("disable-sync", true),
	chromedp.Flag("force-color-profile", "srgb"),
	chromedp.Flag("metrics-recording-only", true),
	chromedp.Flag("safebrowsing-disable-auto-update", true),
	chromedp.Flag("password-store", "basic"),
	chromedp.Flag("use-mock-keychain", true),

	// chromedp.Flag("headless", false),
	chromedp.Flag("aggressive-cache-discard", true),
	// chromedp.Flag("disable-notifications", true),
	// chromedp.Flag("disable-remote-fonts", true),
	// chromedp.Flag("disable-reading-from-canvas", true),
	// chromedp.Flag("disable-remote-playback-api", true),
	// chromedp.Flag("disable-shared-workers", true),
	// chromedp.Flag("disable-voice-input", true),
	// chromedp.Flag("enable-aggressive-domstorage-flushing", true),
	chromedp.Flag("incognito", true),
	// chromedp.Flag("disk-cache-size", "1"),
	// chromedp.Flag("media-cache-size", "1"),
	chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
	// chromedp.Flag("blink-settings", "imagesEnabled=false"),
	chromedp.WindowSize(1600, 1200),
	// chromedp.Flag("auto-open-devtools-for-tabs", true),
	// chromedp.Flag("sec-ch-ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""),
	chromedp.Flag("enable-automation", false),
	chromedp.Flag("disable-blink-features", "AutomationControlled"),
	// chromedp.Flag("dom-automation", false),
	// chromedp.Flag("user-data", "C:\\Users\\User\\AppData\\local\\Google\\Chrome\\UserData\\Profile\\"),
	chromedp.Flag("accept-lang", "en-GB,en-US"),
}
