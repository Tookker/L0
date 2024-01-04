package subscriber

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"L0/internal/config"
	"L0/internal/models"
	"L0/internal/store"
)

type Broker interface {
	Subscribe() error
}

type BrokerNUTS struct {
	connect stan.Conn
	store   store.Store
	logger  *zap.Logger
	subject string
}

func NewBroker(config *config.Config, store store.Store, logger *zap.Logger) (Broker, error) {
	connect, err := stan.Connect(config.ClusterID, config.SubscriberID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		return nil, err
	}

	return &BrokerNUTS{
		connect: connect,
		store:   store,
		logger:  logger,
		subject: config.Subject,
	}, nil
}

func (b *BrokerNUTS) Subscribe() error {
	_, err := b.connect.Subscribe(b.subject, func(msg *stan.Msg) {
		err := msg.Ack()
		if err != nil {
			b.logger.Error(err.Error())
			return
		}

		var order models.Order
		if err = json.Unmarshal(msg.Data, &order); err != nil {
			b.logger.Error(err.Error())
			return
		}

		err = b.store.OrderItem().AddOrder(&order)
		if err != nil {
			return
		}

		b.logger.Info(strings.Join([]string{"order with  order_uid=", order.UID, "stored to database"}, " "))

	}, stan.SetManualAckMode())

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
