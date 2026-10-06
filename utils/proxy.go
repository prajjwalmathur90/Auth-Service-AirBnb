package utils

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ProxyToService builds a reverse proxy that forwards a request to targetBaseUrl
// after replacing publicPrefix with upstreamPrefix in the incoming path.
//
//	ProxyToService("http://localhost:3001", "/api/bookings", "/api/v1/bookings")
//	POST /api/bookings/7?x=1  ->  POST http://localhost:3001/api/v1/bookings/7?x=1
//
// Method, body, query string and headers (including Authorization) are preserved.
// It must be mounted behind JWTAuthMiddleware so the user identity is in the context.
func ProxyToService(targetBaseUrl string, publicPrefix string, upstreamPrefix string) http.HandlerFunc {
	target, err := url.Parse(targetBaseUrl)

	if err != nil || target.Scheme == "" || target.Host == "" {
		// a bad service URL is a startup misconfiguration, fail fast instead of
		// returning a nil handler that would panic on the first request
		log.Fatalf("invalid upstream URL %q: %v", targetBaseUrl, err)
	}

	proxy := &httputil.ReverseProxy{
		// Rewrite is used instead of Director: it runs after hop-by-hop headers are
		// removed, so a client can't send "Connection: X-User-Id" to make the
		// proxy drop the identity headers we add below.
		Rewrite: func(pr *httputil.ProxyRequest) {
			// for modification in request (adding sth in headers or removing some part etc)

			// 1. map the public prefix to the service's real prefix:
			// /api/bookings/7 -> /api/v1/bookings/7
			rewrittenPath := upstreamPrefix + strings.TrimPrefix(pr.In.URL.Path, publicPrefix)
			if rewrittenPath == "" {
				rewrittenPath = "/"
			}
			pr.Out.URL.Path = rewrittenPath
			pr.Out.URL.RawPath = ""

			// 2. point the request at the upstream (scheme, host, base path, Host header)
			pr.SetURL(target)

			// 3. X-Forwarded-For / -Host / -Proto, so the service knows the real client
			pr.SetXForwarded()

			// 4. never trust identity headers coming from the client
			pr.Out.Header.Del("X-User-Id")
			pr.Out.Header.Del("X-User-Email")

			if userId, ok := pr.In.Context().Value("userId").(string); ok {
				pr.Out.Header.Set("X-User-Id", userId)
			}
			if email, ok := pr.In.Context().Value("email").(string); ok {
				pr.Out.Header.Set("X-User-Email", email)
			}

			// 5. forward the JWT so the downstream service can verify it itself.
			// Browser clients send it as a cookie, so turn it into a Bearer header.
			if pr.Out.Header.Get("Authorization") == "" {
				if cookie, err := pr.In.Cookie(AuthCookieName); err == nil && cookie.Value != "" {
					pr.Out.Header.Set("Authorization", "Bearer "+cookie.Value)
				}
			}
		},

		// upstream down / unreachable -> clean JSON 502 instead of an empty response
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error %s %s -> %s: %v", r.Method, r.URL.Path, target.Host, err)
			WriteJsonErrorResponse(w, http.StatusBadGateway, "Upstream service unavailable", nil)
		},
	}

	return proxy.ServeHTTP
}