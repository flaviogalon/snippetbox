package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type appConfig struct {
	addr             string
	staticAssertsDir string
}

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	appConfig *appConfig
}

func main() {
	// Application configuration
	var appCfg appConfig

	flag.StringVar(&appCfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(
		&appCfg.staticAssertsDir,
		"static-dir",
		"./ui/static/",
		"Path to static assets",
	)
	flag.Parse()

	// Application instance
	app := &application{
		errorLog:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:   log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		appConfig: &appCfg,
	}

	// Web Server
	server := &http.Server{
		Addr:     appCfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", appCfg.addr)
	err := server.ListenAndServe()
	app.errorLog.Fatal(err)
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
