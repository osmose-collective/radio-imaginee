package main

import (
	"encoding/json"
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := cli.App{
		Name: "controller",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "data-dir",
				Value:   "../data",
				EnvVars: []string{"CTRL_DATA_DIR"},
			},
			&cli.IntFlag{
				Name:    "history-limit",
				Value:   10,
				EnvVars: []string{"CTRL_HIST_LIMIT"},
			},
			&cli.StringFlag{
				Name:    "liq-telnet-addr",
				Value:   "localhost:5000",
				EnvVars: []string{"CTRL_LIQ_TELNET_ADDR"},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "bind",
						Value:   ":8042",
						EnvVars: []string{"CTRL_BIND"},
					},
				},
				Action: server,
			}, {
				Name: "info",
				Action: func(c *cli.Context) error {
					info, err := getInfo(c)
					if err != nil {
						return err
					}
					out, err := json.MarshalIndent(info, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(out))
					return nil
				},
			}, {
				Name: "history",
				Action: func(c *cli.Context) error {
					history, err := getHistory(c)
					if err != nil {
						return err
					}
					out, err := json.MarshalIndent(history, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(out))
					return nil
				},
			}, {
				Name: "skip",
				Action: func(c *cli.Context) error {
					ret, err := skip(c)
					if err != nil {
						return err
					}
					fmt.Print(ret)
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
