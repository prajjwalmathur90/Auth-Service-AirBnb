package utils

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ProxyToService(targetBaseUrl string, pathPrefix string) http.HandlerFunc {
	target, err := url.Parse(targetBaseUrl)

	if err != nil {
		fmt.Println("URL Parsing Error: ", err)
		return nil;
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// for modification in request (adding sth in headers or removing some part etc)
	originalDirector := proxy.Director

	proxy.Director = func(r *http.Request) {
		originalDirector(r)

		originalPath := r.URL.Path
		strippedPath := strings.TrimPrefix(originalPath, pathPrefix)

		r.URL.Host = target.Host
		r.URL.Path = target.Path + strippedPath
		
		r.Host = target.Host

		if userId, ok := r.Context().Value("userId").(string); ok {
			r.Header.Set("X-User-Id", userId)
		}
	}

	return proxy.ServeHTTP
}