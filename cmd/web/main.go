package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"snippetbox.flaviogalon.github.io/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type appConfig struct {
	addr             string
	staticAssertsDir string
}

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	appConfig      *appConfig
	snippetModel   *models.SnippetModel
	userModel      *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Custom log
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Loading env variables from .env file
	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("MYSQL_USER")
	dbPwd := os.Getenv("MYSQL_PASSWORD")

	// Application configuration
	var appCfg appConfig

	flag.StringVar(&appCfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(
		&appCfg.staticAssertsDir,
		"static-dir",
		"./ui/static/",
		"Path to static assets",
	)
	dsn := flag.String(
		"dsn",
		fmt.Sprintf("%s:%s@/snippetbox?parseTime=true", dbUser, dbPwd),
		"MySQL data source name",
	)
	flag.Parse()

	// Database pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Form Decoder
	formDecoder := form.NewDecoder()

	// Session Manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Application instance
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		appConfig:      &appCfg,
		snippetModel:   &models.SnippetModel{DB: db},
		userModel:      &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Non-default TLS settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Web Server
	server := &http.Server{
		Addr:         appCfg.addr,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %s", appCfg.addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
