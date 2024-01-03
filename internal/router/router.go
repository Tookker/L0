package router

import (
	_ "net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"L0/internal/controller"
	"L0/internal/store"
)

type ChiRouter struct {
	router *chi.Mux
	logger *zap.Logger
	store  store.Store
}

func NewChiRouter(store store.Store, logger *zap.Logger) *ChiRouter {
	return &ChiRouter{
		router: chi.NewRouter(),
		store:  store,
		logger: logger,
	}
}

func (r *ChiRouter) SetHandlers() {
	r.setUIHandlers()
	r.setOrderItemHandlers()
}

func (r *ChiRouter) setUIHandlers() {
	var handlers = controller.NewUIController(r.logger)

	r.router.Route("/", func(rout chi.Router) {
		rout.Get("", handlers.GetMainUi)
	})
}

func (r *ChiRouter) setOrderItemHandlers() {
	var handlers = controller.NewOrderItemController(r.store, r.logger)

	r.router.Route("/order", func(rout chi.Router) {
		rout.Get("/{id}", handlers.GetOrder)
	})
}

func (r *ChiRouter) GetRouter() *chi.Mux {
	return r.router
}
