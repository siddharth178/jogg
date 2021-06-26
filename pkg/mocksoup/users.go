package mocksoup

import (
	"context"

	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

type MockGetUser struct {
	Email    string
	Password string

	Response *models.User
	Err      error
}

func (m *MockGetUser) GetUser(ctx context.Context, email, password string) (*models.User, error) {
	m.Email = email
	m.Password = password
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockGetUserByID struct {
	ID int

	Response *models.User
	Err      error
}

func (m *MockGetUserByID) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	m.ID = id
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockAddUser struct {
	User models.User

	Response *models.User
	Err      error
}

func (m *MockAddUser) AddUser(ctx context.Context, user models.User) (*models.User, error) {
	m.User = user

	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockGetAllUsers struct {
	Filter models.ListFilter

	Response []*models.User
	Err      error
}

func (m *MockGetAllUsers) GetAllUsers(ctx context.Context, f models.ListFilter) ([]*models.User, error) {
	m.Filter = f
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockUpdateUser struct {
	User   models.User
	UserId int

	Response *models.User
	Err      error
}

func (m *MockUpdateUser) UpdateUser(ctx context.Context, u models.User, userId int) (*models.User, error) {
	m.UserId = userId
	m.User = u

	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

type MockDeleteUser struct {
	UserId int

	Response *models.User
	Err      error
}

func (m *MockDeleteUser) DeleteUser(ctx context.Context, userId int) error {
	m.UserId = userId
	return m.Err
}
