package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
		Listeners:    0,
		CurrentTrack: "mylena - hey",
	}

	currentTrackFile, err := ioutil.ReadFile(c.String("data-dir") + "/latest.txt")
	if err != nil {
		return nil, err
	}
	info.CurrentTrack = strings.TrimSpace(string(currentTrackFile))

	// fetch real listeners
	if xmlStr, err := getWithAuth("http://new.radio.lasuitedumonde.com:8000/admin/stats.xml"); err != nil {
		return nil, fmt.Errorf("failed to get XML: %v", err)
	} else {
		var stats Icestats
		xml.Unmarshal([]byte(xmlStr), &stats)
		if stats.Listeners > info.Listeners {
			info.Listeners = stats.Listeners
		}
	}
	return info, nil
}

func getWithAuth(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", "secure")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyText), nil
}
