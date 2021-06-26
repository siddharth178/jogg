package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
	"gitlab.com/toptal/sidd/jogg/pkg/weather"
	"go.uber.org/zap"
)

type AddActivitiesControllerDS interface {
	AddActivities(ctx context.Context, as []models.Activity, userId int) ([]*models.Activity, error)
}

type AddActivitiesController struct {
	logger *zap.SugaredLogger
	ds     AddActivitiesControllerDS
	ws     weather.WeatherService
}

func (c *AddActivitiesController) AddActivities(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	as := []models.Activity{}
	if err := ctx.Bind(&as); err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for i := 0; i < len(as); i++ {
		w, err := c.ws.GetWeather(ctx.Request().Context(), as[i].Loc, as[i].Ts)
		if err != nil {
			c.logger.Warnw("weather error", "loc", as[i].Loc, "ts", as[i].Ts)
		}
		as[i].WeatherCondition = w.Condition
		as[i].WeatherDescription = w.Description
	}

	activities, err := c.ds.AddActivities(ctx.Request().Context(), as, userId)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, activities)
}

type GetActivityControllerDS interface {
	GetActivity(ctx context.Context, activityId, userId int) (*models.Activity, error)
}

type GetActivityController struct {
	logger *zap.SugaredLogger
	ds     GetActivityControllerDS
}

func (c *GetActivityController) GetActivity(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	activityId, err := strconv.Atoi(ctx.Param("activity_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a, err := c.ds.GetActivity(ctx.Request().Context(), activityId, userId)
	switch errors.Cause(err) {
	case nil:
		// pass
	case pgx.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, a)
}

type DeleteActivityControllerDS interface {
	DeleteActivity(ctx context.Context, activityId, userId int) error
}

type DeleteActivityController struct {
	logger *zap.SugaredLogger
	ds     DeleteActivityControllerDS
}

func (c *DeleteActivityController) DeleteActivity(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	activityId, err := strconv.Atoi(ctx.Param("activity_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.ds.DeleteActivity(ctx.Request().Context(), activityId, userId)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

type UpdateActivityControllerDS interface {
	UpdateActivity(ctx context.Context, a models.Activity, userId int) (*models.Activity, error)
}

type UpdateActivityController struct {
	logger *zap.SugaredLogger
	ds     UpdateActivityControllerDS
}

func (c *UpdateActivityController) UpdateActivity(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	a := models.Activity{}
	if err := ctx.Bind(&a); err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	activity, err := c.ds.UpdateActivity(ctx.Request().Context(), a, userId)
	switch errors.Cause(err) {
	case nil:
		// pass
	case pgx.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, activity)
}

type GetAllActivitiesControllerDS interface {
	GetAllActivities(ctx context.Context, userId int, f models.Filter) ([]*models.Activity, error)
}

type GetAllActivitiesController struct {
	logger *zap.SugaredLogger
	ds     GetAllActivitiesControllerDS
}

const DefaultPageRows = 2

func (c *GetAllActivitiesController) GetAllActivities(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	f := models.Filter{}
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
	f.EntityTableName = "activities"

	a, err := c.ds.GetAllActivities(ctx.Request().Context(), userId, f)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, a)
}
