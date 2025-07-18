package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	driverlicense "github.com/adohong4/driving-license/internal/driver_license"
	"github.com/adohong4/driving-license/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapDriverLicenseRoutes(driverLicenseGroup *echo.Group, h driverlicense.Handlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	driverLicenseGroup.POST("/create", h.CreateDriverLicense(), mw.AuthJWTMiddleware(authUC, cfg))
	driverLicenseGroup.PUT("/:id", h.UpdateDriverLicense(), mw.AuthJWTMiddleware(authUC, cfg))
	driverLicenseGroup.DELETE("/:id", h.DeleteDriverLicense(), mw.AuthJWTMiddleware(authUC, cfg))
	driverLicenseGroup.GET("/:id", h.GetDriverLicenseById())
	driverLicenseGroup.GET("/getAll", h.GetDriverLicense())
	driverLicenseGroup.GET("/search", h.SearchByLicenseNo())
}
