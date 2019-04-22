package main

import (
	telnet "github.com/reiver/go-telnet"
	cli "gopkg.in/urfave/cli.v2"
)

func skip(c *cli.Context) (string, error) {
	connection, err := telnet.DialTo(c.String("liq-telnet-addr"))
	if err != nil {
		return "", err
	}

	if _, err := connection.Write([]byte("main(dot)harbor.skip\r\nexit")); err != nil {
		return "", err
	}

	var bytes []byte
	if _, err := connection.Read(bytes); err != nil {
		return "", err
	}

	return string(bytes), nil
}
