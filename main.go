package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var appConfig = &config{}

// We need something to start with
var logger = &log.Logger{
	Out:       os.Stdout,
	Formatter: new(log.TextFormatter),
	Level:     log.ErrorLevel,
}

var (
	listenAddress = flag.String("bind-address", "localhost", "the address to bind to.")
	serverPort    = flag.Int("port", 8080, "the port to listen on.")
	configFile    = flag.String("config-file", "config.yaml", "the location of the file to use for configuration.")
)

var reloadHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		loadLoggingConfig(*configFile, appConfig)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var configHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		content, _ := yaml.Marshal(&appConfig)
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Log the requests coming into our web server
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	// Read the file and compute its hash upon startup
	loadLoggingConfig(*configFile, appConfig)
	serverAddress := fmt.Sprintf("%s:%d", *listenAddress, *serverPort)

	// Register the HTTP handlers
	http.HandleFunc("/-/reload", reloadHandler)
	http.HandleFunc("/-/config", configHandler)

	// Start the server!
	logger.Printf("Server is starting at %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, logRequest(http.DefaultServeMux)); err != http.ErrServerClosed {
		logger.Error(err)
	}
}

// Parse the logging configuration into something useable
func setLogConfig(lc *logging, l *log.Logger) {
	level, err := log.ParseLevel(lc.Level)
	if err != nil {
		logger.Warnf("Problem parsing log level: %v, defaulting to INFO", err)
		l.SetLevel(log.InfoLevel)
	} else {
		l.SetLevel(level)
	}

	switch lc.Format {
	case "text":
		l.SetFormatter(&log.TextFormatter{DisableColors: !lc.Colors})
	case "json":
		l.SetFormatter(new(log.JSONFormatter))
	default:
		l.SetFormatter(&log.TextFormatter{DisableColors: !lc.Colors})
	}
}

func loadLoggingConfig(path string, current *config) {
	logger.Debug("Reading configuration file from ", path)
	conf := readFile(path)
	hash := hashConfig(conf)

	logger.Debug("Config file hash: ", hash)

	logConf := config{Hash: hash}
	yaml.Unmarshal(conf, &logConf)

	if current.Hash != hash {
		current.Hash = hash
		logger.Info("Change in logging configuration found! Reloading configuration")
		current.Logging = logConf.Logging
		setLogConfig(&logConf.Logging, logger)
	} else {
		logger.Info("No change detected. Skipping reload.")
	}
}

// Do not compute the hash this way if your config file is large
func hashConfig(config []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(config))
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Errorf("Could not read file at %s", path)
		os.Exit(1)
	}
	return content
}
