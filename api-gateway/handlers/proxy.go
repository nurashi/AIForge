package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(targetURL string) http.Handler {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("ERROR: Invalid target URL for reverse proxy '%s': %v", targetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)


	return proxy
}

