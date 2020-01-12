package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	hs "github.com/mitchellh/hashstructure"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var appConfig *config
var logger = &log.Logger{
	Out:       os.Stdout,
	Formatter: new(log.TextFormatter),
	Level:     log.ErrorLevel,
}

const filepath = "config.yaml"

func main() {
	appConfig = loadConfigFromFile(filepath)
	serverAddress := appConfig.Server.Address + ":" + appConfig.Server.Port

	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			loadLoggingConfig(filepath, appConfig, logger)
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
	http.ListenAndServe(serverAddress, logRequest(http.DefaultServeMux))
}

func getEnvDefault(env string, dflt string) string {
	if res := os.Getenv(env); len(res) > 0 {
		return res
	}

	return dflt
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
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
		if !lc.Colors {
			l.SetFormatter(&log.TextFormatter{DisableColors: true})
		} else {
			l.SetFormatter(new(log.TextFormatter))
		}
	default:
		l.SetFormatter(new(log.JSONFormatter))
	}
}

func loadConfigFromFile(path string) *config {
	newConfig := &config{}
	loadServerConfig(path, newConfig)
	loadLoggingConfig(path, newConfig, logger)
	return newConfig
}

func loadLoggingConfig(path string, current *config, l *log.Logger) {
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
		logger.Info("Hashes match. Skipping reload.")
	}

}

func loadServerConfig(path string, in *config) {
	conf := readFile(path)

	type temp struct { // wrap server in a temporary outer struct
		Server server `yaml:"server"`
	}

	newServer := temp{}
	yaml.Unmarshal(conf, &newServer)
	if newServer.Server.Address == "" {
		logger.Error("Address is required.")
		os.Exit(1)
	} else if newServer.Server.Port == "" {
		logger.Error("Port is required.")
		os.Exit(1)
	}
	in.Server = newServer.Server
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Errorf("Could not read file at %s", path)
		os.Exit(1)
	}
	return content
}
