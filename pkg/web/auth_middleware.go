package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func AdminOrUserManager(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if HasJWTToken(c) && (IsAdmin(c) || IsUserManager(c)) {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "access allowed only for admin or usermanager")
		}
	}
}

func ValidUser(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if HasJWTToken(c) {
				if IsAdmin(c) ||
					IsUserManager(c) ||
					(c.Param("user_id") == GetJWTTokenClaims(c)[ClaimUserId]) {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access to user data")
		}
	}
}

func AdminOrValidUser(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if HasJWTToken(c) {
				if IsAdmin(c) ||
					(c.Param("user_id") == GetJWTTokenClaims(c)[ClaimUserId]) {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access to user data")
		}
	}
}
