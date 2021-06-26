package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
	"go.uber.org/zap"
)

type WeeklyReportControllerDS interface {
	GetWeeklyReport(ctx context.Context, userId int, f models.ListFilter) ([]*models.Weekly, error)
}

type WeeklyReportController struct {
	logger *zap.SugaredLogger
	ds     WeeklyReportControllerDS
}

func (c *WeeklyReportController) WeeklyReport(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	f := models.ListFilter{}
	if err := ctx.Bind(&f); err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageRows == 0 {
		f.PageRows = DefaultPageRows
	}

	a, err := c.ds.GetWeeklyReport(ctx.Request().Context(), userId, f)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, a)
}
