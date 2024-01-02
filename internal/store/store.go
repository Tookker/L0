package store

import "L0/internal/store/tables"

type Store interface {
	OrderItem() tables.OrderItem
}
