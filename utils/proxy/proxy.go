package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// ForwardRequest proxies a request to the exact upstream URL supplied by the
// route handler. It preserves request bodies and supports WebSocket upgrades.
func ForwardRequest(c echo.Context, targetURL string) error {
	target, err := url.Parse(targetURL)
	if err != nil || target.Scheme == "" || target.Host == "" {
		return echo.NewHTTPError(http.StatusBadGateway, "Invalid target URL")
	}

	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			query := req.URL.Query()
			targetQuery := target.Query()
			for key, values := range targetQuery {
				query[key] = values
			}

			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = target.Path
			req.URL.RawPath = target.RawPath
			req.URL.RawQuery = query.Encode()
			req.Host = target.Host

			req.Header.Set("X-API-Gateway", "true")
			req.Header.Set("X-API-Gateway-Version", "1.0")
			if requestID := c.Response().Header().Get("X-Request-ID"); requestID != "" {
				req.Header.Set("X-Request-ID", requestID)
			}
			clientIP := c.RealIP()
			req.Header.Set("X-Forwarded-For", clientIP)
			req.Header.Set("X-Real-IP", clientIP)
			if c.Request().TLS != nil {
				req.Header.Set("X-Forwarded-Proto", "https")
			} else {
				req.Header.Set("X-Forwarded-Proto", "http")
			}

			if gatewaySecret := os.Getenv("GATEWAY_SECRET_KEY"); gatewaySecret != "" {
				req.Header.Set("X-Gateway-Secret", gatewaySecret)
			} else {
				req.Header.Del("X-Gateway-Secret")
			}
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, proxyErr error) {
			c.Logger().Errorf("gateway upstream request failed: %v", proxyErr)
			http.Error(w, "Upstream service unavailable", http.StatusBadGateway)
		},
		ModifyResponse: func(response *http.Response) error {
			// The API gateway owns CORS. Forwarding upstream CORS headers would
			// produce duplicate values such as "*, *", which browsers reject.
			response.Header.Del("Access-Control-Allow-Origin")
			response.Header.Del("Access-Control-Allow-Methods")
			response.Header.Del("Access-Control-Allow-Headers")
			response.Header.Del("Access-Control-Allow-Credentials")
			response.Header.Del("Access-Control-Expose-Headers")
			response.Header.Del("Access-Control-Max-Age")
			return nil
		},
	}

	reverseProxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

func JoinURL(baseURL, servicePath string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(servicePath, "/")
}
