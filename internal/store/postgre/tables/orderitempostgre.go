package tables

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"L0/internal/caches/ordercache"
	"L0/internal/models"
	"L0/internal/store/tables"
)

var (
	FillCacheErr = "Filling cahce error!"
)

type OrderItemRepo struct {
	db         *sql.DB
	orderCache ordercache.OrderCache
	logger     *zap.Logger
}

/*
@brief - Заполнить кэш данным из БД
@param db - указатель на иницилиазированный объект БД
@param orderCache - кэш, который нужно заполнить данными
@param logger - логгер
@return error - описание возникшей ошибке, при успешном выполнении вернется nil
*/
func fillCache(db *sql.DB, orderCache ordercache.OrderCache, logger *zap.Logger) error {
	request := `
	SELECT 
		orders_items.id_order,
		orders.order_uid, 
		orders.track_number, 
		orders.entry,
		deliverys.name, 
		deliverys.phone, 
		deliverys.zip, 
		deliverys.city, 
		deliverys.address, 
		deliverys.region,
		deliverys.email,
		payments.transaction, 
		payments.request_id, 
		payments.currency, 
		payments.provider, 
		payments.amount, 
		payments.payment_dt,
		payments.bank, 
		payments.delivery_cost, 
		payments.goods_total, 
		payments.custom_fee, 
		orders.locale,
		orders.internal_signature, 
		orders.customer_id, 
		orders.delivery_service, 
		orders.shardkey, 
		orders.sm_id, 
		orders.date_created, 
		orders.oof_shard,
		items.chrt_id, 
		items.track_number,
		items.price, 
		items.rid, 
		items.name, 
		items.sale, 
		items.size, 
		items.total_price, 
		items.nm_id, 
		items.brand, 
		items.status
	FROM 
		orders_items
	JOIN 
		orders ON orders.id = orders_items.id_order
	JOIN 
		deliverys ON deliverys.id = orders.id_delivery
	JOIN 
		payments ON payments.id = orders.id_payment
	JOIN 
		items ON items.id = orders_items.id_item;
	`

	res, err := db.Query(request)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Close()

	bufMap := make(map[uint]models.Order)

	for res.Next() {
		var item models.Item
		var order models.Order
		var idOrder uint
		var requestID sql.NullString
		var internSign sql.NullString

		err = res.Scan(&idOrder, &order.UID, &order.TrackNum, &order.Entry,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &requestID, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank,
			&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&order.Locale, &internSign, &order.CustomerID, &order.DeliveryService, &order.ShardKey,
			&order.SmID, &order.CreationDate, &order.CofShard,
			&item.ChrtID, &item.TrackNum, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		if requestID.Valid {
			order.Payment.RequestID = requestID.String
		}

		if internSign.Valid {
			order.InternSign = internSign.String
		}

		val, ok := bufMap[idOrder]
		if ok {
			val.Items = append(val.Items, item)
			bufMap[idOrder] = val
		} else {
			order.Items = append(order.Items, item)
			bufMap[idOrder] = order
		}
	}

	for id, o := range bufMap {
		err := orderCache.AddOrder(id, &o)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	return nil
}

/*
@brief - создание объекта таблицы OrderItem
@param db - инициализированная БД
@param orderCache - объект кэша
@param logger - логгер
@return tables.OrderItem - объект OrderItem, при ошибке возвращает nil
@return error - описание возникшей ошибке. при успещном выполнение возвращает nil
*/
func NewOrderItemRepo(db *sql.DB, orderCache ordercache.OrderCache, logger *zap.Logger) (tables.OrderItem, error) {

	err := fillCache(db, orderCache, logger)
	if err != nil {
		logger.Error(FillCacheErr)
		return nil, errors.New(FillCacheErr)
	}

	return &OrderItemRepo{
		db:         db,
		orderCache: orderCache,
		logger:     logger,
	}, nil
}

func (o *OrderItemRepo) GetOrder(ctx context.Context, id uint) (models.Order, error) {
	order, err := o.orderCache.GetOrder(ctx, id)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (o *OrderItemRepo) AddOrder(order *models.Order) error {
	//TODO реализовать добавление заказа в БД и КЭШ
	return nil
}
