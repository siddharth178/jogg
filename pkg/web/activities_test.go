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
	"gitlab.com/toptal/sidd/jogg/pkg/weather"
)

func TestAddActivityController_HappyPath(t *testing.T) {
	aJSON := `[{"ts":"2021-06-02T13:43:04+00:00", "loc":"Pune", "distance":3500}]`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")
	mock := mocksoup.MockAddActivities{
		Response: []*models.Activity{
			{Id: 1, UserId: 100, Loc: "Pune"},
		},
	}
	wsMock := weather.MockWeatherService{Response: &weather.DefaultWeather}
	h := &AddActivitiesController{
		ds: &mock,
		ws: &wsMock,
	}

	err := h.AddActivities(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "Pune", mock.Activities[0].Loc)
		assert.Equal(t, 100, mock.UserId)

		assert.Equal(t, "Pune", wsMock.Loc)

		resp := []models.Activity{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Pune", resp[0].Loc)
		assert.Equal(t, 100, resp[0].UserId)
		assert.Equal(t, 1, resp[0].Id)
	}
}

func TestAddActivityController_SadPath_DBError(t *testing.T) {
	aJSON := `[{"ts":"2021-06-02T13:43:04+00:00", "loc":"Pune", "distance":3500}]`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockAddActivities{
		Err: assert.AnError,
	}
	h := &AddActivitiesController{
		ds: &mock,
		ws: &weather.MockWeatherService{Response: &weather.DefaultWeather},
	}

	err := h.AddActivities(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}

func TestGetActivityController_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities/:activity_id")
	c.SetParamNames("user_id", "activity_id")
	c.SetParamValues("100", "1")

	mock := mocksoup.MockGetActivity{
		Response: &models.Activity{Id: 1, UserId: 100, Loc: "Pune"},
	}
	h := &GetActivityController{
		ds: &mock,
	}

	err := h.GetActivity(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 1, mock.ActivityId)
		assert.Equal(t, 100, mock.UserId)

		resp := models.Activity{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Pune", resp.Loc)
		assert.Equal(t, 100, resp.UserId)
		assert.Equal(t, 1, resp.Id)
	}
}

func TestGetActivityController_SadPath_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities/:activity_id")
	c.SetParamNames("user_id", "activity_id")
	c.SetParamValues("100", "1")

	mock := mocksoup.MockGetActivity{
		Err: assert.AnError,
	}
	h := &GetActivityController{
		ds: &mock,
	}

	err := h.GetActivity(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}

func TestGetActivityController_SadPath_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities/:activity_id")
	c.SetParamNames("user_id", "activity_id")
	c.SetParamValues("100", "1")

	mock := mocksoup.MockGetActivity{
		Err: pgx.ErrNoRows,
	}
	h := &GetActivityController{
		ds: &mock,
	}

	err := h.GetActivity(c)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
	}
}

func TestDeleteActivityController_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities/:activity_id")
	c.SetParamNames("user_id", "activity_id")
	c.SetParamValues("100", "1")

	mock := mocksoup.MockDeleteActivity{}
	h := &DeleteActivityController{
		ds: &mock,
	}

	err := h.DeleteActivity(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 1, mock.ActivityId)
		assert.Equal(t, 100, mock.UserId)
	}
}

func TestDeleteActivityController_SadPath_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities/:activity_id")
	c.SetParamNames("user_id", "activity_id")
	c.SetParamValues("100", "1")

	mock := mocksoup.MockDeleteActivity{
		Err: assert.AnError,
	}
	h := &DeleteActivityController{
		ds: &mock,
	}

	err := h.DeleteActivity(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}

func TestUpdateActivityController_HappyPath(t *testing.T) {
	aJSON := `{"id":1, "ts":"2021-06-02T13:43:04+00:00", "loc":"Pune", "distance":3500}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")
	mock := mocksoup.MockUpdateActivity{
		Response: &models.Activity{Id: 1, UserId: 100, Loc: "Pune", Distance: 3500},
	}
	h := &UpdateActivityController{
		ds: &mock,
	}

	err := h.UpdateActivity(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "Pune", mock.Activity.Loc)
		assert.Equal(t, 100, mock.UserId)

		resp := models.Activity{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, 3500, resp.Distance)
		assert.Equal(t, "Pune", resp.Loc)
		assert.Equal(t, 100, resp.UserId)
		assert.Equal(t, 1, resp.Id)
	}
}

func TestUpdateActivityController_SadPath_DBError(t *testing.T) {
	aJSON := `{"ts":"2021-06-02T13:43:04+00:00", "loc":"Pune", "distance":3500}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockUpdateActivity{
		Err: assert.AnError,
	}
	h := &UpdateActivityController{
		ds: &mock,
	}

	err := h.UpdateActivity(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}

func TestUpdateActivityController_SadPath_NotFound(t *testing.T) {
	aJSON := `{"ts":"2021-06-02T13:43:04+00:00", "loc":"Pune", "distance":3500}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(aJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockUpdateActivity{
		Err: pgx.ErrNoRows,
	}
	h := &UpdateActivityController{
		ds: &mock,
	}

	err := h.UpdateActivity(c)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
	}
}

func TestGetAllActivitiesController_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockGetAllActivities{
		Response: []*models.Activity{
			{Id: 2, UserId: 100, Loc: "L1"},
			{Id: 1, UserId: 100, Loc: "L2"},
		},
	}
	h := &GetAllActivitiesController{
		ds: &mock,
	}

	err := h.GetAllActivities(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 100, mock.UserId)

		resp := []models.Activity{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		if assert.Len(t, resp, 2) {
			assert.Equal(t, "L1", resp[0].Loc)
			assert.Equal(t, 100, resp[0].UserId)
			assert.Equal(t, 2, resp[0].Id)
			assert.Equal(t, "L2", resp[1].Loc)
			assert.Equal(t, 100, resp[1].UserId)
			assert.Equal(t, 1, resp[1].Id)
		}
	}
}

func TestGetAllActivitiesController_SadPath_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/activities")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockGetAllActivities{
		Err: assert.AnError,
	}
	h := &GetAllActivitiesController{
		ds: &mock,
	}

	err := h.GetAllActivities(c)
	assert.Equal(t, errors.Cause(assert.AnError), err)
}
