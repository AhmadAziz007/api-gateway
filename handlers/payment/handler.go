package payment

import (
	"ecommerce/api-gateway/utils/proxy"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	// Midtrans calls this endpoint directly; the payment-service validates its
	// signature, so it must not require a customer JWT.
	g.POST("/payments/webhook", forward("/payments/webhook"))
}

func RegisterRoutes(g *echo.Group) {
	payments := g.Group("/payments")
	payments.GET("", forward("/auth/payments"))
	payments.GET("/:id", forwardParam("/auth/payments/", "id"))
	payments.POST("", forward("/auth/payments"))
	payments.POST("/:orderID/confirm", forwardParamSuffix("/auth/payments/", "orderID", "/confirm"))

	admin := g.Group("/admin/payments")
	admin.GET("", forward("/admin/payments"))
	admin.GET("/:id", forwardParam("/admin/payments/", "id"))
}

func forward(servicePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return proxy.ForwardRequest(c, proxy.JoinURL(paymentServiceURL(), servicePath))
	}
}

func forwardParam(prefix, parameter string) echo.HandlerFunc {
	return forwardParamSuffix(prefix, parameter, "")
}

func forwardParamSuffix(prefix, parameter, suffix string) echo.HandlerFunc {
	return func(c echo.Context) error {
		servicePath := prefix + url.PathEscape(c.Param(parameter)) + suffix
		return proxy.ForwardRequest(c, proxy.JoinURL(paymentServiceURL(), servicePath))
	}
}

func paymentServiceURL() string {
	if serviceURL := os.Getenv("PAYMENT_SERVICE_URL"); serviceURL != "" {
		return serviceURL
	}
	return "http://localhost:8084"
}
