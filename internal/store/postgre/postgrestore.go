package postgre

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"L0/internal/caches/ordercache/ordercachemap"
	"L0/internal/config"
	"L0/internal/store"
	"L0/internal/store/postgre/tables"
	iTables "L0/internal/store/tables"
)

type Postgre struct {
	db        *sql.DB
	logger    *zap.Logger
	orderItem iTables.OrderItem
}

var (
	OpenPosgreErr = "Open postgresql db error!"
)

func open(config *config.Config, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DataBase)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(OpenPosgreErr)
	}

	return db, nil
}

func NewPostgre(config *config.Config, logger *zap.Logger) (store.Store, error) {
	db, err := open(config, logger)
	if err != nil {
		return nil, err
	}

	orderItem, err := tables.NewOrderItemRepo(db, ordercachemap.NewCache(logger), logger)
	if err != nil {
		return nil, err
	}

	return &Postgre{
		db:        db,
		logger:    logger,
		orderItem: orderItem,
	}, nil
}

func (p *Postgre) OrderItem() iTables.OrderItem {
	return p.orderItem
}

func Close(db *Postgre) {
	db.db.Close()
}
