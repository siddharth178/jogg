package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/toptal/sidd/jogg/pkg/mocksoup"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

func TestLogin_HappyPath(t *testing.T) {
	loginJSON := `{"email":"admin@jogg.in", "password":"passwd"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	mock := mocksoup.MockGetUser{
		Response: &models.User{Id: 100, Email: "admin@jogg.in", Password: "passwd"},
	}
	h := &LoginController{
		ds:     &mock,
		secret: "secret",
	}

	err := h.Login(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "admin@jogg.in", mock.Email)
		assert.Equal(t, "passwd", mock.Password)

		resp := LoginResp{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotEmpty(t, resp.Token)
	}
}

func TestLogin_SadPath_DBError(t *testing.T) {
	loginJSON := `{"name":"admin", "password":"passwd"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	h := &LoginController{
		ds: &mocksoup.MockGetUser{
			Err: pgx.ErrNoRows,
		},
		secret: "secret",
	}

	err := h.Login(c)
	if assert.Error(t, err) && assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code) {
		assert.Equal(t, "", strings.TrimSpace(rec.Body.String()))
	}
}
