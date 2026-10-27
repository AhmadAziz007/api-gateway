package notification

import (
	"ecommerce/api-gateway/utils/proxy"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group) {
	notifications := g.Group("/notifications")
	notifications.GET("", forward("/auth/notifications"))
	// Register the static WebSocket route before the :id route.
	notifications.GET("/ws", forwardWebSocket)
	notifications.GET("/:id", forwardParam("/auth/notifications/", "id"))
	notifications.PUT("/:id", forwardParam("/auth/notifications/", "id"))
	notifications.PUT("/:id/read", forwardParam("/auth/notifications/", "id"))
}

func forward(servicePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return proxy.ForwardRequest(c, proxy.JoinURL(notificationServiceURL(), servicePath))
	}
}

func forwardParam(prefix, parameter string) echo.HandlerFunc {
	return func(c echo.Context) error {
		servicePath := prefix + url.PathEscape(c.Param(parameter))
		return proxy.ForwardRequest(c, proxy.JoinURL(notificationServiceURL(), servicePath))
	}
}

func forwardWebSocket(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return echo.NewHTTPError(401, "Invalid token claims")
	}
	query := url.Values{"user_id": []string{strconv.FormatInt(userID, 10)}}
	target := proxy.JoinURL(notificationServiceURL(), "/ws") + "?" + query.Encode()
	return proxy.ForwardRequest(c, target)
}

func notificationServiceURL() string {
	if serviceURL := os.Getenv("NOTIFICATION_SERVICE_URL"); serviceURL != "" {
		return serviceURL
	}
	return "http://localhost:8085"
}
