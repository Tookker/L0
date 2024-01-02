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
	const request = `
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

	tx, err := o.db.Begin()
	if err != nil {
		o.logger.Error(err.Error())
		return err
	}

	const insertDeliverys = `
	INSERT INTO deliverys (name, phone, zip, city, address, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;
	`
	var idDelivery uint

	res := tx.QueryRow(insertDeliverys, &order.Delivery.Name, &order.Delivery.Phone,
		&order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
		&order.Delivery.Region, &order.Delivery.Email)

	err = res.Scan(&idDelivery)
	if err != nil {
		o.logger.Error(err.Error())
		tx.Rollback()
		return err
	}

	const insertPayments = `
	INSERT INTO payments (transaction, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;
	`
	var idPayments uint

	res = tx.QueryRow(insertPayments, &order.Payment.Transaction, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)

	err = res.Scan(&idPayments)
	if err != nil {
		o.logger.Error(err.Error())
		tx.Rollback()
		return err
	}

	const insertOrder = `
	INSERT INTO orders (order_uid, track_number, entry, id_delivery, id_payment,
	locale,  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;
	`
	var idOrder uint
	res = tx.QueryRow(insertOrder, &order.UID, &order.TrackNum, &order.Entry, &idDelivery, &idPayments,
		&order.Locale, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID,
		&order.CreationDate, &order.CofShard)

	err = res.Scan(&idOrder)
	if err != nil {
		o.logger.Error(err.Error())
		tx.Rollback()
		return err
	}

	const insertItem = `
	INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
	`
	sizeItems := len(order.Items)
	arrIdItem := make([]uint, 0, sizeItems)

	for i := 0; i < sizeItems; i++ {
		var id uint
		res = tx.QueryRow(insertItem, &order.Items[i].ChrtID, &order.Items[i].TrackNum, &order.Items[i].Price, &order.Items[i].RID,
			&order.Items[i].Name, &order.Items[i].Sale, &order.Items[i].Size, &order.Items[i].TotalPrice, &order.Items[i].NmID,
			&order.Items[i].Brand, &order.Items[i].Status)

		err = res.Scan(&id)
		if err != nil {
			o.logger.Error(err.Error())
			tx.Rollback()
			return err
		}

		arrIdItem = append(arrIdItem, id)
	}

	const insertOrderItem = `
	INSERT INTO orders_items (id_order, id_item)
	VALUES ($1, $2);
	`
	for i := 0; i < sizeItems; i++ {
		_, err = tx.Exec(insertOrderItem, &idOrder, &arrIdItem[i])
		if err != nil {
			o.logger.Error(err.Error())
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		o.logger.Error(err.Error())
		return err
	}

	err = o.orderCache.AddOrder(idOrder, order)
	if err != nil {
		return err
	}

	return nil
}
