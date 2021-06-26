package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
	"go.uber.org/zap"
)

type GetUserByIDControllerDS interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

type GetUserByIDController struct {
	ds GetUserByIDControllerDS
}

func (c *GetUserByIDController) GetUserByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := c.ds.GetUserByID(ctx.Request().Context(), id)
	switch errors.Cause(err) {
	case nil:
		// pass
	case pgx.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
		// return ctx.JSON(http.StatusNotFound, err)
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

type AddUserControllerDS interface {
	AddUser(ctx context.Context, user models.User) (*models.User, error)
}

type AddUserController struct {
	ds AddUserControllerDS
}

func (c *AddUserController) AddUser(ctx echo.Context) error {
	u := models.User{}
	if err := ctx.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !HasJWTToken(ctx) && u.Role != models.RoleUser {
		return echo.NewHTTPError(http.StatusUnauthorized, "only normal users can self-register")
	}

	user, err := c.ds.AddUser(ctx.Request().Context(), u)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

type GetAllUsersControllerDS interface {
	GetAllUsers(ctx context.Context, f models.ListFilter) ([]*models.User, error)
}

type GetAllUsersController struct {
	logger *zap.SugaredLogger
	ds     GetAllUsersControllerDS
}

func (c *GetAllUsersController) GetAllUsers(ctx echo.Context) error {
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

	a, err := c.ds.GetAllUsers(ctx.Request().Context(), f)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, a)
}

type UpdateUserControllerDS interface {
	UpdateUser(ctx context.Context, u models.User, userId int) (*models.User, error)
}

type UpdateUserController struct {
	logger *zap.SugaredLogger
	ds     UpdateUserControllerDS
}

func (c *UpdateUserController) UpdateUser(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	u := models.User{}
	if err := ctx.Bind(&u); err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := c.ds.UpdateUser(ctx.Request().Context(), u, userId)
	switch errors.Cause(err) {
	case nil:
		// pass
	case pgx.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

type DeleteUserControllerDS interface {
	DeleteUser(ctx context.Context, userId int) error
}

type DeleteUserController struct {
	logger *zap.SugaredLogger
	ds     DeleteUserControllerDS
}

func (c *DeleteUserController) DeleteUser(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.logger.Error("error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.ds.DeleteUser(ctx.Request().Context(), userId)
	switch errors.Cause(err) {
	case nil:
		// pass
	default:
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}
