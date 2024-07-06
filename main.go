package main

import (
	"flag"
	"main/src"
	"net/http"
	"os"
)

var (
	webListenAddress = flag.String("web.listem.address", "0.0.0.0:9172", "Address to listen")
	logLevel         = flag.String("log.level", "info", "Only log messages with the given severity or above. One of: [debug, info, warn, error]")
	logFormat        = flag.String("log.format", "text", "Output format of log messages. One of: [text, json]")
)

func main() {
	// Parse flags
	flag.Parse()

	// Init logger
	src.Init(logLevel, logFormat)
	logger := src.GetLogger()

	logger.Info("Exporter is started")

	// Start HTTP server
	server := http.NewServeMux()
	server.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		src.Router(writer, request, logger)
	})
	logger.Info("Starting HTTP server on " + *webListenAddress + "")
	err := http.ListenAndServe(*webListenAddress, server)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
