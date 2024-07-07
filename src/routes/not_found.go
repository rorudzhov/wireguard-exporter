package routes

import (
	"log/slog"
	"net/http"
)

func NotFound(writer http.ResponseWriter, request *http.Request, logger *slog.Logger) {

	body := `
<!DOCTYPE html>
<html lang="en">
<head>
   <meta charset="UTF-8">
   <title>404</title>
</head>
   <body>
       <h4>Wireguard-exporter</h4>
       <p>This page not found</p>
       <p>Please see metrics page <a href="/metrics">/metrics</a></p>
   </body>
</html>
`
	writer.WriteHeader(404)
	logger.Info("Page not found", "url", request.RequestURI, "method", request.Method, "remote", request.RemoteAddr, "code", 404)
	writer.Write([]byte(body))
}
