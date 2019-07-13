package main

import (
	"encoding/xml"
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
	Errors       []string
}

type Icestats struct {
	XMLName   xml.Name `xml:"icestats"`
	Listeners int      `xml:"listeners"`
}

func getHistory(c *cli.Context) ([]string, error) {
	f, err := ioutil.ReadFile(c.String("data-dir") + "/history.txt")
	if err != nil {
		return nil, err
	}
	history := strings.Split(string(f), "\n")
	last := c.Int("history-limit")
	if last < len(history) {
		history = history[len(history)-last:]
	}
	return history, nil
}

func getInfo(c *cli.Context) (*Info, error) {
	info := &Info{
		CurrentTime:  time.Now(),
		Listeners:    0, // default (if different from 0, it will be the minimum number displayed)
		CurrentTrack: "",
		Errors:       []string{},
	}

	currentTrackFile, err := ioutil.ReadFile(c.String("data-dir") + "/latest.txt")
	if err != nil {
		info.Errors = append(info.Errors, err.Error())
		info.CurrentTrack = "Unknown song"
	} else {
		info.CurrentTrack = strings.TrimSpace(string(currentTrackFile))
	}

	// fetch real listeners
	xmlStr, err := getWithAuth("https://stream.osmose.world/admin/stats.xml")
	if err != nil {
		info.Errors = append(info.Errors, err.Error())
		info.Listeners = 42 // fake value when cannot get the real one :)
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
