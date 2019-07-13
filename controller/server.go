package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	cli "gopkg.in/urfave/cli.v2"
)

func server(c *cli.Context) error {
	r := chi.NewRouter()
	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			info, err := getInfo(c)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			out, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
		})
		r.Get("/history", func(w http.ResponseWriter, r *http.Request) {
			history, err := getHistory(c)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			out, err := json.MarshalIndent(history, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
		})
		r.Get("/playlists/info", func(w http.ResponseWriter, r *http.Request) {
			info, err := getPlaylistsInfo(c)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			out, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
		})
		r.Get("/skip", func(w http.ResponseWriter, r *http.Request) {
			msg, err := skip(c)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Write([]byte(msg))
		})
	})
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "www")
	FileServer(r, "/", http.Dir(filesDir))

	fmt.Printf("Starting server on %q...\n", c.String("bind"))
	return http.ListenAndServe(c.String("bind"), r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
