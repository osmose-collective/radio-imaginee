package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	cli "gopkg.in/urfave/cli.v2"
)

type Info struct {
	CurrentTime  time.Time
	Listeners    int
	CurrentTrack string
}

type Icestats struct {
	XMLName   xml.Name `xml:"icestats"`
	Listeners int      `xml:"listeners"`
}

func getInfo(c *cli.Context) (*Info, error) {
	info := &Info{
		CurrentTime:  time.Now(),
		Listeners:    42,
		CurrentTrack: "mylena - hey",
	}

	// fetch real listeners
	if xmlStr, err := getWithAuth("http://new.radio.lasuitedumonde.com:8000/admin/stats.xml"); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		var stats Icestats
		xml.Unmarshal([]byte(xmlStr), &stats)
		if stats.Listeners > info.Listeners {
			info.Listeners = stats.Listeners
		}
		(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", "secure")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	data := string(bodyText)

	return string(data), nil
}
