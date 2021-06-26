package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/toptal/sidd/jogg/pkg/mocksoup"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

func TestWeeklyReportController_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=3&page_rows=5", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/reports/:user_id/weekly")
	c.SetParamNames("user_id")
	c.SetParamValues("100")

	mock := mocksoup.MockGetWeeklyReport{
		Response: []*models.Weekly{
			{Week: time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC), AvgDistance: 100},
			{Week: time.Date(2021, 5, 24, 0, 0, 0, 0, time.UTC), AvgDistance: 120},
		},
	}
	h := &WeeklyReportController{
		ds: &mock,
	}

	err := h.WeeklyReport(c)
	if assert.NoError(t, err) && assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, 100, mock.UserId)
		assert.Equal(t, models.ListFilter{Page: 3, PageRows: 5}, mock.F)

		resp := []models.Weekly{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		if assert.Len(t, resp, 2) {
			assert.Equal(t, float32(100), resp[0].AvgDistance)
			assert.Equal(t, float32(120), resp[1].AvgDistance)
		}
	}
}
