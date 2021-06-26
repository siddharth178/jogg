package datastore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

func scanActivity(row pgx.Row) (*models.Activity, error) {
	a := models.Activity{}
	if err := row.Scan(&a.Id, &a.UserId, &a.Ts, &a.Loc, &a.Distance, &a.CreatedAt, &a.UpdatedAt, &a.WeatherCondition, &a.WeatherDescription, &a.Seconds); err != nil {
		return nil, errors.WithStack(err)
	}
	a.Ts = a.Ts.Local()
	return &a, nil
}

func scanActivities(rows pgx.Rows) ([]*models.Activity, error) {
	as := []*models.Activity{}
	for rows.Next() {
		a := models.Activity{}
		if err := rows.Scan(&a.Id, &a.UserId, &a.Ts, &a.Loc, &a.Distance, &a.CreatedAt, &a.UpdatedAt, &a.WeatherCondition, &a.WeatherDescription, &a.Seconds); err != nil {
			return nil, errors.WithStack(err)
		}
		a.Ts = a.Ts.Local()
		as = append(as, &a)
	}
	return as, nil
}

func (ds *DS) AddActivities(ctx context.Context, as []models.Activity, userId int) ([]*models.Activity, error) {
	q := `INSERT INTO activities(user_id, ts, loc, distance, seconds, w_cond, w_desc, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *`
	activities := []*models.Activity{}

	tx, err := ds.db.Begin(context.Background())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer tx.Rollback(context.Background())

	for _, a := range as {
		activity, err := scanActivity(tx.QueryRow(ctx, q, userId, a.Ts.UTC(), a.Loc, a.Distance, a.Seconds, a.WeatherCondition, a.WeatherDescription, time.Now(), time.Now()))
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return activities, nil
}

func (ds *DS) GetActivity(ctx context.Context, activityId, userId int) (*models.Activity, error) {
	q := `SELECT * FROM activities WHERE
	id = $1
	AND user_id = $2`

	activity, err := scanActivity(ds.db.QueryRow(ctx, q, activityId, userId))
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (ds *DS) DeleteActivity(ctx context.Context, activityId, userId int) error {
	q := `DELETE FROM activities WHERE
	id = $1
	AND user_id = $2`

	_, err := ds.db.Exec(ctx, q, activityId, userId)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ds *DS) UpdateActivity(ctx context.Context, a models.Activity, userId int) (*models.Activity, error) {
	q := `UPDATE activities
SET user_id=$1, ts=$2, loc=$3, distance=$4, seconds=$5 w_cond=$6, w_desc=$7, updated_at=$8
WHERE id=$9
RETURNING *`
	activity, err := scanActivity(ds.db.QueryRow(ctx, q, userId, a.Ts.UTC().Unix(), a.Loc, a.Distance, a.Seconds, a.WeatherCondition, a.WeatherDescription, time.Now(), a.Id))
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (ds *DS) GetAllActivities(ctx context.Context, userId int, f models.Filter) ([]*models.Activity, error) {
	ds.logger.Infof("filter: %+v", f)
	q1 := `SELECT * FROM activities
WHERE
	true
	AND user_id = $1`
	q2 := `
ORDER BY
	ts DESC
LIMIT $2 OFFSET $3`

	l, o := f.GetLimitOffset()
	condStr, err := f.BuildQuery()
	if err != nil {
		return nil, err
	}
	if condStr != "" {
		condStr = " AND " + condStr
		q1 = q1 + condStr + q2
	} else {
		q1 = q1 + q2
	}

	rows, err := ds.db.Query(ctx, q1, userId, l, o)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	activities, err := scanActivities(rows)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
