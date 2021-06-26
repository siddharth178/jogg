package mocksoup

import (
	"context"

	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

type MockAddActivities struct {
	Activities []models.Activity
	UserId     int

	Response []*models.Activity
	Err      error
}

func (m *MockAddActivities) AddActivities(ctx context.Context, as []models.Activity, userId int) ([]*models.Activity, error) {
	m.UserId = userId
	m.Activities = as

	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockGetActivity struct {
	UserId     int
	ActivityId int

	Response *models.Activity
	Err      error
}

func (m *MockGetActivity) GetActivity(ctx context.Context, activityId, userId int) (*models.Activity, error) {
	m.ActivityId = activityId
	m.UserId = userId
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Response, nil
}

type MockDeleteActivity struct {
	UserId     int
	ActivityId int

	Response *models.Activity
	Err      error
}

func (m *MockDeleteActivity) DeleteActivity(ctx context.Context, activityId, userId int) error {
	m.ActivityId = activityId
	m.UserId = userId
	return m.Err
}

type MockUpdateActivity struct {
	Activity models.Activity
	UserId   int

	Response *models.Activity
	Err      error
}

func (m *MockUpdateActivity) UpdateActivity(ctx context.Context, a models.Activity, userId int) (*models.Activity, error) {
	m.UserId = userId
	m.Activity = a

	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockGetAllActivities struct {
	UserId int
	Filter models.Filter

	Response []*models.Activity
	Err      error
}

func (m *MockGetAllActivities) GetAllActivities(ctx context.Context, userId int, f models.Filter) ([]*models.Activity, error) {
	m.UserId = userId
	m.Filter = f
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}
