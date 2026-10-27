package cors

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Middleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		// Keep the browser origin explicit. This also remains valid if
		// credentials are added to API requests later.
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.HEAD, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
			"X-Idempotency-Key",
		},
		AllowCredentials: true,
	})
}
