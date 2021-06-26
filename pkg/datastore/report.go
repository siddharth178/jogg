package datastore

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/models"
)

func scanWeeklies(rows pgx.Rows) ([]*models.Weekly, error) {
	as := []*models.Weekly{}
	for rows.Next() {
		a := models.Weekly{}
		if err := rows.Scan(&a.Week, &a.Days, &a.AvgDistance, &a.AvgKMs, &a.AvgHours, &a.AvgSpeedKMPH); err != nil {
			return nil, errors.WithStack(err)
		}
		as = append(as, &a)
	}
	return as, nil
}

func (ds *DS) GetWeeklyReport(ctx context.Context, userId int, f models.ListFilter) ([]*models.Weekly, error) {
	q := `SELECT
	date_trunc('week', ts) AS week,
	COUNT(1) days,
	AVG(distance) avg_distance,
	AVG(distance/1000.0) AS avg_kms,
	AVG(seconds/3600.0) AS avg_hours,
	AVG(distance/1000.0)/AVG(seconds/3600.0) avg_speed_kmph
FROM activities
WHERE user_id=$1
GROUP BY week
ORDER BY week DESC
LIMIT $2 OFFSET $3`

	l, o := f.GetLimitOffset()
	rows, err := ds.db.Query(ctx, q, userId, l, o)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	weeklies, err := scanWeeklies(rows)
	if err != nil {
		return nil, err
	}

	return weeklies, nil
}
