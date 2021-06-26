package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

const (
	ClaimUserId = "userId"
	ClaimName   = "name"
	ClaimEmail  = "email"
	ClaimRole   = "role"
	ClaimExp    = "exp"
)

type LoginControllerDS interface {
	GetUser(ctx context.Context, name, password string) (*models.User, error)
}

type LoginController struct {
	ds     LoginControllerDS
	secret string
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

func (c *LoginController) Login(ctx echo.Context) error {
	lr := &LoginReq{}
	if err := ctx.Bind(lr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u, err := c.ds.GetUser(ctx.Request().Context(), lr.Email, lr.Password)
	switch errors.Cause(err) {
	case nil:
		//pass
	case pgx.ErrNoRows:
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	default:
		return err
	}

	// generate JWT token and return it
	token, err := generateToken(*u, c.secret)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, LoginResp{token})
}

func generateToken(user models.User, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.SigningString()

	claims := token.Claims.(jwt.MapClaims)
	claims[ClaimUserId] = fmt.Sprintf("%d", user.Id)
	claims[ClaimName] = user.Name
	claims[ClaimEmail] = user.Email
	claims[ClaimRole] = user.Role
	claims[ClaimExp] = time.Now().Add(time.Minute * 60).Unix()

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return t, nil
}
