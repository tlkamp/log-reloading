package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	hs "github.com/mitchellh/hashstructure"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var appConfig *config

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

func main() {
	flag.Parse()
	appConfig = loadConfigFromFile(*configFile)
	serverAddress := fmt.Sprintf("%s:%d", *listenAddress, *serverPort)

	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			loadLoggingConfig(*configFile, appConfig)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/-/config", func(w http.ResponseWriter, r *http.Request) {
		content, _ := yaml.Marshal(&appConfig)
		w.Write(content)
	})

	logger.Printf("Server is starting at %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, logRequest(http.DefaultServeMux)); err != http.ErrServerClosed {
		log.Error(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

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
		l.SetFormatter(&log.TextFormatter{DisableColors: lc.Colors})
	default:
		l.SetFormatter(new(log.JSONFormatter))
	}
}

func loadConfigFromFile(path string) *config {
	newConfig := &config{}
	loadLoggingConfig(path, newConfig)
	return newConfig
}

func loadLoggingConfig(path string, current *config) {
	conf := readFile(path)

	type temp struct { // wrap logging in a temporary outer struct
		Logging logging `yaml:"logging"`
	}

	logConf := temp{}
	yaml.Unmarshal(conf, &logConf)
	h, _ := hs.Hash(logConf, nil)
	newHash := fmt.Sprintf("%d", h)

	if (logging{}) == current.Logging { // initial configuration load
		current.Logging = logConf.Logging
		current.Logging.Hash = newHash
		setLogConfig(&current.Logging, logger)
	} else if current.Logging.Hash != newHash {
		logger.Info("Change in logging configuration found! Reloading configuration")
		logger.Debugf("Hash before: %s Hash after: %s", current.Logging.Hash, newHash)
		logConf.Logging.Hash = newHash
		current.Logging = logConf.Logging
		setLogConfig(&logConf.Logging, logger)
	} else {
		logger.Info("No change detected. Skipping reload.")
	}
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Errorf("Could not read file at %s", path)
		os.Exit(1)
	}
	return content
}
