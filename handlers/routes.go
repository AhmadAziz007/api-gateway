package handlers

import (
	"ecommerce/api-gateway/handlers/notification"
	"ecommerce/api-gateway/handlers/order"
	"ecommerce/api-gateway/handlers/payment"
	"ecommerce/api-gateway/handlers/product"
	"ecommerce/api-gateway/handlers/user"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	user.RegisterPublicRoutes(g)
	product.RegisterPublicRoutes(g)
	order.RegisterPublicRoutes(g)
	payment.RegisterPublicRoutes(g)
}

func RegisterProtectedRoutes(g *echo.Group) {
	user.RegisterProtectedRoutes(g)
	product.RegisterProtectedRoutes(g)

	order.RegisterRoutes(g)
	payment.RegisterRoutes(g)

	notification.RegisterRoutes(g)
}

func RegisterAllRoutes(g *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
	RegisterPublicRoutes(g)

	protected := g.Group("")
	protected.Use(jwtMiddleware)
	RegisterProtectedRoutes(protected)
}
