package tables

import (
	"context"

	"L0/internal/models"
)

type OrderItem interface {
	/*
		@brief Получить информацию о заказе
		@param ctx - контекст
		@param id - id заказа
		@return models.Order - заказ, при ошибке вернется пустая структура models.Order
		@return error - ошибка, при успещном поиске вернется  nil
	*/
	GetOrder(ctx context.Context, id uint) (models.Order, error)
	/*
		@brief Добавить информацию о заказе
		@param order - описание заказа
		@return error - ошибка, при успещном поиске вернется nil
	*/
	AddOrder(id uint, order *models.Order) error
}
