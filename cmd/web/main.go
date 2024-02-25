package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

type appConfig struct {
	addr             string
	staticAssertsDir string
}

func main() {
	var appCfg appConfig

	flag.StringVar(&appCfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(
		&appCfg.staticAssertsDir,
		"static-dir",
		"./ui/static/",
		"Path to static assets",
	)
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(appCfg.staticAssertsDir)})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", appCfg.addr)
	err := http.ListenAndServe(appCfg.addr, mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
