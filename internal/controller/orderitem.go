package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"L0/internal/store"
)

var (
	inputDataErr = "Not valid input data"
)

type OrderItem interface {
	GetOrder(w http.ResponseWriter, r *http.Request)
}

type OrderItemController struct {
	store  store.Store
	logger *zap.Logger
}

func NewOrderItemController(store store.Store, logger *zap.Logger) OrderItem {
	return &OrderItemController{
		store:  store,
		logger: logger,
	}
}

func (o *OrderItemController) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := strings.Split(r.URL.String(), "/")
	id, err := strconv.Atoi(res[len(res)-1])
	if err != nil || id <= 0 {
		o.logger.Error(inputDataErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := o.store.OrderItem().GetOrder(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytesOrder, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		o.logger.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytesOrder)
}
