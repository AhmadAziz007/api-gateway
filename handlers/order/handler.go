package order

import (
	"ecommerce/api-gateway/utils/proxy"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	g.GET("/public/orders/:orderCode/code", forwardOrderCode("/public/orders/", "/code"))
}

func RegisterRoutes(g *echo.Group) {
	orders := g.Group("/orders")
	orders.GET("", forward("/auth/orders"))
	orders.POST("", forward("/auth/orders"))
	orders.GET("/:orderID", forwardParam("/auth/orders/", "orderID"))
	orders.GET("/code/:orderCode/code", forwardOrderCode("/auth/orders/", "/code"))

	admin := g.Group("/admin/orders")
	admin.GET("", forward("/admin/orders"))
	admin.GET("/:orderID", forwardParam("/admin/orders/", "orderID"))
	admin.PUT("/:orderID/status", forwardParamSuffix("/admin/orders/", "orderID", "/status"))
	admin.DELETE("/:orderID", forwardParam("/admin/orders/", "orderID"))
}

func forward(servicePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return proxy.ForwardRequest(c, proxy.JoinURL(orderServiceURL(), servicePath))
	}
}

func forwardParam(prefix, parameter string) echo.HandlerFunc {
	return forwardParamSuffix(prefix, parameter, "")
}

func forwardParamSuffix(prefix, parameter, suffix string) echo.HandlerFunc {
	return func(c echo.Context) error {
		servicePath := prefix + url.PathEscape(c.Param(parameter)) + suffix
		return proxy.ForwardRequest(c, proxy.JoinURL(orderServiceURL(), servicePath))
	}
}

func forwardOrderCode(prefix, suffix string) echo.HandlerFunc {
	return forwardParamSuffix(prefix, "orderCode", suffix)
}

func orderServiceURL() string {
	if serviceURL := os.Getenv("ORDER_SERVICE_URL"); serviceURL != "" {
		return serviceURL
	}
	return "http://localhost:8083"
}
