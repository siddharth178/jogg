package web

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthControllerDS interface {
	DBPing(ctx context.Context) error
}

type HealthController struct {
	ds HealthControllerDS
}

func (c *HealthController) Health(ctx echo.Context) error {
	if err := c.ds.DBPing(ctx.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "ok")
}
