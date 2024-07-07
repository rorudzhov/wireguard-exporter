package src

import (
	"log/slog"
	"main/src/routes"
	"net/http"
)

// Router routes handler
func Router(writer http.ResponseWriter, request *http.Request, logger *slog.Logger) {
	switch request.RequestURI {
	case "/metrics":
		routes.Metrics(writer, request, logger)
	default:
		routes.NotFound(writer, request, logger)
	}

}
