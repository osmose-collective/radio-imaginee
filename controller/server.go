package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	cli "gopkg.in/urfave/cli.v2"
)

func server(c *cli.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<body bgcolor=black>
  <audio autoload="true" autoplay="true" controls="true" id="player-background">
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-160.mp3" type="audio/mpeg" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-128.aac" type="audio/aac" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-192.mp3" type="audio/mpeg" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee.ogg" type="audio/ogg" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-128.mp3" type="audio/mpeg" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-160.aac" type="audio/aac" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-192.aac" type="audio/aac" />
    <source src="http://new.radio.lasuitedumonde.com:8000/radio-imaginee-64.mp3" type="audio/mpeg" />
  </audio>
</body>
`))
	})
	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		info, err := getInfo(c)
		out, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			// FIXME: display error
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	})
	fmt.Printf("Starting server on %q...\n", c.String("bind"))
	return http.ListenAndServe(c.String("bind"), r)
}
