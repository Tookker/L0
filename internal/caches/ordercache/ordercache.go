package ordercache

import (
	"context"

	"L0/internal/models"
)

type OrderCache interface {
	/*
		@brief - Получить заказ по его ID
		@param orderID - контекст
		@param orderID - ID заказа
		@return models.Order - заказ, в случае ошибки возвращается nil
		@return error - описание ошибки, в случае успешного выполенния метода возвращается nil
	*/
	GetOrder(ctx context.Context, orderID uint) (models.Order, error)
	/*
		@brief - Добавить заказ в кэш
		@param order - Описание заказа
		@param orderID - ID заказа
		@return - nil в случае успешного добавления, в случае ошибки её описание
	*/
	AddOrder(orderID uint, order *models.Order) error
}
