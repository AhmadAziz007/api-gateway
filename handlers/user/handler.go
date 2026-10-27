package user

import (
	"ecommerce/api-gateway/utils/proxy"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	users := g.Group("/users")
	users.POST("/register", forward("/signup"))
	users.POST("/signin", forward("/signin"))
	users.GET("/verify-account", forward("/verify-account"))
	users.GET("/verify-email", forward("/verify-account"))
	users.POST("/forgot-password", forward("/forgot-password"))
	users.PUT("/update-password", forward("/update-password"))
	users.PUT("/reset-password", forward("/update-password"))
}

func RegisterProtectedRoutes(g *echo.Group) {
	users := g.Group("/users")
	users.GET("/profile", forward("/auth/profile"))
	users.PUT("/profile", forward("/auth/profile"))
	users.POST("/upload-avatar", forward("/auth/profile/image-upload"))
	users.GET("/roles", forward("/admin/roles"))

	admin := g.Group("/admin")
	admin.GET("/customers", forward("/admin/customers"))
	admin.GET("/customers/:id", forwardParam("/admin/customers/", "id"))
	admin.POST("/customers", forward("/admin/customers"))
	admin.PUT("/customers/:id", forwardParam("/admin/customers/", "id"))
	admin.DELETE("/customers/:id", forwardParam("/admin/customers/", "id"))
	admin.GET("/roles", forward("/admin/roles"))
	admin.GET("/roles/:id", forwardParam("/admin/roles/", "id"))
	admin.POST("/roles", forward("/admin/roles"))
	admin.PUT("/roles/:id", forwardParam("/admin/roles/", "id"))
	admin.DELETE("/roles/:id", forwardParam("/admin/roles/", "id"))
}

func forward(servicePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return proxy.ForwardRequest(c, proxy.JoinURL(userServiceURL(), servicePath))
	}
}

func forwardParam(prefix, parameter string) echo.HandlerFunc {
	return func(c echo.Context) error {
		servicePath := prefix + url.PathEscape(c.Param(parameter))
		return proxy.ForwardRequest(c, proxy.JoinURL(userServiceURL(), servicePath))
	}
}

func userServiceURL() string {
	if serviceURL := os.Getenv("USER_SERVICE_URL"); serviceURL != "" {
		return serviceURL
	}
	return "http://localhost:8080"
}
