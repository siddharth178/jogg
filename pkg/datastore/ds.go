package datastore

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/config"
	"go.uber.org/zap"
)

type DS struct {
	config config.DBConfig
	logger *zap.SugaredLogger

	db *pgxpool.Pool
}

func NewDS(ctx context.Context, cfg config.DBConfig, logger *zap.SugaredLogger) (*DS, error) {
	config, err := pgxpool.ParseConfig(cfg.ConnURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	config.ConnConfig.Logger = zapadapter.NewLogger(logger.Desugar())
	config.ConnConfig.LogLevel = pgx.LogLevelInfo

	dbpool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &DS{
		config: cfg,
		logger: logger,
		db:     dbpool,
	}, nil
}

func (ds *DS) DBPing(ctx context.Context) error {
	if err := ds.db.Ping(ctx); err != nil {
		return errors.WithStack(err)
	}

	// simple placeholder query
	ds.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, "select count(*) from users")
		cnt := 0
		if err := row.Scan(&cnt); err != nil {
			return errors.WithStack(err)
		}
		log.Infof("users: %d", cnt)
		return nil
	})

	return nil
}
