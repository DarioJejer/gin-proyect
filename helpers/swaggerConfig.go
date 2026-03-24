package helpers

import (
	"net/url"
	"os"
	"strings"

	docs "Gin/docs"
	"log"
)

// configureSwaggerForEnvironment sets Swagger host/scheme.
// Priority: PUBLIC_BASE_URL → HEROKU_APP_NAME → localhost.
func ConfigureSwaggerForEnvironment(port string) {
	if base := strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL")); base != "" {
		u, err := url.Parse(base)
		if err != nil || u.Host == "" {
			log.Printf("swagger: invalid PUBLIC_BASE_URL %q, using localhost", base)
			setSwaggerLocalhost(port)
			return
		}
		docs.SwaggerInfo.Host = u.Host
		scheme := strings.ToLower(u.Scheme)
		if scheme == "" {
			scheme = "https"
		}
		docs.SwaggerInfo.Schemes = []string{scheme}
		return
	}

	if app := strings.TrimSpace(os.Getenv("HEROKU_APP_NAME")); app != "" {
		docs.SwaggerInfo.Host = app + ".herokuapp.com"
		docs.SwaggerInfo.Schemes = []string{"https"}
		return
	}

	setSwaggerLocalhost(port)
}

func setSwaggerLocalhost(port string) {
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.Schemes = []string{"http"}
}
