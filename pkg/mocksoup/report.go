package mocksoup

import (
	"context"

	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

type MockGetWeeklyReport struct {
	UserId int
	F      models.ListFilter

	Response []*models.Weekly
	Err      error
}

func (m *MockGetWeeklyReport) GetWeeklyReport(ctx context.Context, userId int, f models.ListFilter) ([]*models.Weekly, error) {
	m.UserId = userId
	m.F = f
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}
