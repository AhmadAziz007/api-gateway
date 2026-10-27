package product

import (
	"ecommerce/api-gateway/utils/proxy"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	products := g.Group("/products")
	products.GET("/shop", forward("/products/shop"))
	products.GET("/home", forward("/products/home"))
	products.GET("/home/:id", forwardParam("/products/home/", "id"))

	categories := g.Group("/categories")
	categories.GET("/home", forward("/categories/home"))
	categories.GET("/shop", forward("/categories/shop"))
}

func RegisterProtectedRoutes(g *echo.Group) {
	products := g.Group("/admin/products")
	products.GET("", forward("/admin/products"))
	products.POST("", forward("/admin/products"))
	products.GET("/:id", forwardParam("/admin/products/", "id"))
	products.PUT("/:id", forwardParam("/admin/products/", "id"))
	products.DELETE("/:id", forwardParam("/admin/products/", "id"))
	products.POST("/upload-image", forward("/admin/image-upload"))

	admin := g.Group("/admin")
	admin.POST("/image-upload", forward("/admin/image-upload"))
	admin.PUT("/update-status-product/:id", forwardParam("/admin/update-status-product/", "id"))
	admin.GET("/categories", forward("/admin/categories"))
	admin.GET("/categories/:id", forwardParam("/admin/categories/", "id"))
	admin.GET("/categories/:slug/slug", func(c echo.Context) error {
		servicePath := "/admin/categories/" + url.PathEscape(c.Param("slug")) + "/slug"
		return proxy.ForwardRequest(c, proxy.JoinURL(productServiceURL(), servicePath))
	})
	admin.POST("/categories", forward("/admin/categories"))
	admin.PUT("/categories/:id", forwardParam("/admin/categories/", "id"))
	admin.DELETE("/categories/:id", forwardParam("/admin/categories/", "id"))

	cart := g.Group("/cart")
	cart.GET("", forward("/auth/cart"))
	cart.POST("", forward("/auth/cart"))
	cart.DELETE("", forward("/auth/cart"))
	cart.DELETE("/all", forward("/auth/cart/all"))
	cart.DELETE("/:id", forwardParam("/auth/cart/", "id"))
}

func forward(servicePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return proxy.ForwardRequest(c, proxy.JoinURL(productServiceURL(), servicePath))
	}
}

func forwardParam(prefix, parameter string) echo.HandlerFunc {
	return func(c echo.Context) error {
		servicePath := prefix + url.PathEscape(c.Param(parameter))
		return proxy.ForwardRequest(c, proxy.JoinURL(productServiceURL(), servicePath))
	}
}

func productServiceURL() string {
	if serviceURL := os.Getenv("PRODUCT_SERVICE_URL"); serviceURL != "" {
		return serviceURL
	}
	return "http://localhost:8082"
}
