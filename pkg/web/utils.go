package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

func HasJWTToken(c echo.Context) bool {
	_, ok := c.Get("user").(*jwt.Token)
	return ok
}

func GetJWTTokenClaims(c echo.Context) jwt.MapClaims {
	t, ok := c.Get("user").(*jwt.Token)
	if ok {
		return t.Claims.(jwt.MapClaims)
	}
	return nil
}

func IsAdmin(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims[ClaimRole].(string) == models.RoleAdmin
}

func IsUserManager(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims[ClaimRole].(string) == models.RoleUserManager
}

func IsUser(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims[ClaimRole].(string) == models.RoleUser
}
