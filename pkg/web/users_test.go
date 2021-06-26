package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/toptal/sidd/jogg/pkg/mocksoup"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

func TestGetUserControllerByID_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")
	mock := mocksoup.MockGetUserByID{
		Response: &models.User{Id: 100, Email: "admin@jogg.in", Password: "passwd"},
	}
	h := &GetUserByIDController{
		ds: &mock,
	}

	err := h.GetUserByID(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 100, mock.ID)

		resp := models.User{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "admin@jogg.in", resp.Email)
	}
}

func TestGetUserControllerByID_SadPath_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")
	mock := mocksoup.MockGetUserByID{
		Err: pgx.ErrNoRows,
	}
	h := &GetUserByIDController{
		ds: &mock,
	}

	err := h.GetUserByID(c)
	if assert.Error(t, err) && assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code) {
		assert.Equal(t, "", strings.TrimSpace(rec.Body.String()))
	}
}

func TestAddUserController_HappyPath(t *testing.T) {
	userJSON := `{"email":"s@jogg.in", "password":"s", "name":"s", "role":"user"}`
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	mock := mocksoup.MockAddUser{
		Response: &models.User{Id: 100, Email: "s@jogg.in", Password: "s", Name: "s", Role: models.RoleUser},
	}
	h := &AddUserController{
		ds: &mock,
	}

	err := h.AddUser(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "s@jogg.in", mock.User.Email)

		resp := models.User{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "s@jogg.in", resp.Email)
		assert.Equal(t, 100, resp.Id)
	}
}

func TestAddUserController_SadPath_Admin(t *testing.T) {
	userJSON := `{"email":"s@jogg.in", "password":"s", "name":"s", "role":"admin"}`
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	mock := mocksoup.MockAddUser{
		Response: &models.User{Id: 100, Email: "s@jogg.in", Password: "s", Name: "s", Role: models.RoleUser},
	}
	h := &AddUserController{
		ds: &mock,
	}

	err := h.AddUser(c)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
	}
}

func TestGetAllUsersController_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page_rows=5&page=2", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	mock := mocksoup.MockGetAllUsers{
		Response: []*models.User{
			{Id: 2, Email: "admin@jogg.com"},
			{Id: 3, Email: "s@jogg.com"},
		},
	}
	h := &GetAllUsersController{
		ds: &mock,
	}

	err := h.GetAllUsers(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 2, mock.Filter.Page)
		assert.Equal(t, 5, mock.Filter.PageRows)

		resp := []models.User{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		if assert.Len(t, resp, 2) {
			assert.Equal(t, "admin@jogg.com", resp[0].Email)
			assert.Equal(t, 2, resp[0].Id)
			assert.Equal(t, "s@jogg.com", resp[1].Email)
			assert.Equal(t, 3, resp[1].Id)
		}
	}
}

func TestGetAllUsersController_SadPath_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	mock := mocksoup.MockGetAllUsers{
		Err: assert.AnError,
	}
	h := &GetAllUsersController{
		ds: &mock,
	}

	err := h.GetAllUsers(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}

func TestUpdateUserController_HappyPath(t *testing.T) {
	aJSON := `{"email":"bbb77@jogg.in", "password":"bbb", "name":"bbb", "role":"user"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")
	mock := mocksoup.MockUpdateUser{
		Response: &models.User{Id: 1, Email: "bbb77@jogg.in"},
	}
	h := &UpdateUserController{
		ds: &mock,
	}

	err := h.UpdateUser(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "bbb77@jogg.in", mock.User.Email)
		assert.Equal(t, 100, mock.UserId)

		resp := models.User{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "bbb77@jogg.in", resp.Email)
		assert.Equal(t, 1, resp.Id)
	}
}

func TestUpdateUserController_SadPath_DBError(t *testing.T) {
	aJSON := `{"email":"bbb77@jogg.in", "password":"bbb", "name":"bbb", "role":"user"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockUpdateUser{
		Err: assert.AnError,
	}
	h := &UpdateUserController{
		ds: &mock,
	}

	err := h.UpdateUser(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}

func TestUpdateUserController_SadPath_NotFound(t *testing.T) {
	aJSON := `{"email":"bbb77@jogg.in", "password":"bbb", "name":"bbb", "role":"user"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockUpdateUser{
		Err: pgx.ErrNoRows,
	}
	h := &UpdateUserController{
		ds: &mock,
	}

	err := h.UpdateUser(c)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
	}
}

func TestDeleteUserController_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockDeleteUser{}
	h := &DeleteUserController{
		ds: &mock,
	}

	err := h.DeleteUser(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 100, mock.UserId)
	}
}

func TestDeleteUserController_SadPath_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockDeleteUser{
		Err: assert.AnError,
	}
	h := &DeleteUserController{
		ds: &mock,
	}

	err := h.DeleteUser(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}
