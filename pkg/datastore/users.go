package datastore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func scanUser(row pgx.Row) (*models.User, error) {
	u := models.User{}
	if err := row.Scan(&u.Id, &u.Email, &u.Password, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, errors.WithStack(err)
	}
	return &u, nil
}

func scanUsers(rows pgx.Rows) ([]*models.User, error) {
	us := []*models.User{}
	for rows.Next() {
		u := models.User{}
		if err := rows.Scan(&u.Id, &u.Email, &u.Password, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, errors.WithStack(err)
		}

		us = append(us, &u)
	}
	return us, nil
}

func (ds *DS) GetUser(ctx context.Context, email, password string) (*models.User, error) {

	u, err := scanUser(ds.db.QueryRow(ctx, "SELECT * FROM users where email=$1", email))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, err
	}

	return u, nil
}

func (ds *DS) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	u, err := scanUser(ds.db.QueryRow(ctx, "SELECT * FROM users where id=$1", id))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ds *DS) AddUser(ctx context.Context, user models.User) (*models.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO users(email, password, name, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *`
	u, err := scanUser(ds.db.QueryRow(ctx, q, user.Email, string(pass), user.Name, user.Role, time.Now(), time.Now()))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ds *DS) GetAllUsers(ctx context.Context, f models.ListFilter) ([]*models.User, error) {
	q := `SELECT * FROM users
ORDER BY
	email ASC
LIMIT $1 OFFSET $2`

	l, o := f.GetLimitOffset()
	rows, err := ds.db.Query(ctx, q, l, o)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ds *DS) UpdateUser(ctx context.Context, u models.User, userId int) (*models.User, error) {
	q := `UPDATE users
SET email=$1, password=$2, name=$3, role=$4, updated_at=$5
WHERE id=$6
RETURNING *`
	user, err := scanUser(ds.db.QueryRow(ctx, q, u.Email, u.Password, u.Name, u.Role, time.Now(), userId))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ds *DS) DeleteUser(ctx context.Context, userId int) error {
	q := `DELETE FROM users WHERE
	id = $1`

	_, err := ds.db.Exec(ctx, q, userId)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
