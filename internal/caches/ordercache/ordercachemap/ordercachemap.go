package ordercachemap

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"L0/internal/caches/ordercache"
	"L0/internal/models"
)

type OrderCacheMap struct {
	logger   *zap.Logger
	cacheMap map[uint]models.Order
}

var (
	OrderNotFindErr = "Order not found!"
	OrderExistsErr  = "Order already exists!"
)

func NewCache(logger *zap.Logger) ordercache.OrderCache {
	return &OrderCacheMap{
		logger:   logger,
		cacheMap: make(map[uint]models.Order),
	}
}

func (o *OrderCacheMap) GetOrder(ctx context.Context, orderID uint) (models.Order, error) {
	order, ok := o.cacheMap[orderID]
	orderIDStr := strconv.FormatUint(uint64(orderID), 10)

	if !ok {
		o.logger.Info(strings.Join([]string{"Order", orderIDStr, "not found!"}, " "))
		return models.Order{}, errors.New(OrderNotFindErr)
	}

	o.logger.Info(strings.Join([]string{"The order", orderIDStr, "was successfully found."}, " "))
	return order, nil
}

func (o *OrderCacheMap) AddOrder(orderID uint, order *models.Order) error {
	orderIDStr := strconv.FormatUint(uint64(orderID), 10)

	_, ok := o.cacheMap[orderID]
	if ok {
		o.logger.Info(strings.Join([]string{"Order", orderIDStr, "already exists!"}, " "))
		return errors.New(OrderExistsErr)
	}

	o.cacheMap[orderID] = *order
	o.logger.Info(strings.Join([]string{"The order", orderIDStr, "was successfully added."}, " "))
	return nil
}
