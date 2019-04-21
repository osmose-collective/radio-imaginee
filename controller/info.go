package main

import (
	"time"

	cli "gopkg.in/urfave/cli.v2"
)

type Info struct {
	CurrentTime  time.Time
	Listeners    int
	CurrentTrack string
}

func getInfo(c *cli.Context) (*Info, error) {
	info := &Info{
		CurrentTime:  time.Now(),
		Listeners:    42,
		CurrentTrack: "mylena - hey",
	}
	return info, nil
}
