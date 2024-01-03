package tables_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"L0/internal/caches/ordercache/ordercachemap"
	"L0/internal/models"
	"L0/internal/store/postgre/tables"
)

var (
	ID       uint = 1
	orderVar      = models.Order{
		UID:      "b563feb7b2b84b6test",
		TrackNum: "WBILMTESTTRACK",
		Entry:    "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:     9934930,
				TrackNum:   "WBILMTESTTRACK",
				Price:      453,
				RID:        "ab4219087a764ae0btest",
				Name:       "Mascaras",
				Sale:       30,
				Size:       "0",
				TotalPrice: 317,
				NmID:       2389212,
				Brand:      "Vivienne Sabo",
				Status:     202,
			},
		},
		Locale:          "en",
		InternSign:      "",
		CustomerID:      "test",
		DeliveryService: "meest",
		ShardKey:        "9",
		SmID:            99,
		CreationDate:    "2021-11-26T06:22:19Z",
		CofShard:        "1",
	}
)

func getZapLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Err create logger.")
		return nil, err
	}

	return logger, nil
}

func openDb(logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=127.0.0.1 user=admin password=root dbname=orders port=5432 sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return db, nil
}

func TestAddOrder(t *testing.T) {
	logger, err := getZapLogger()
	if err != nil {
		t.Error("Failed init logger")
		return
	}

	db, err := openDb(logger)
	if err != nil {
		t.Error("Failed open DB")
		return
	}

	cache := ordercachemap.NewCache(logger)

	orderItem, err := tables.NewOrderItemRepo(db, cache, logger)
	if err != nil {
		t.Error("Failed create order_item repos")
		return
	}

	type args struct {
		order models.Order
	}

	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "Add order",
			args: args{
				order: orderVar,
			},
			want:    nil,
			wantErr: false,
		},
	}

	t.Logf("Start AddOrder")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := orderItem.AddOrder(&tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrder error = %v, whantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v pass!", tt.name)
		})
	}

	t.Logf("Finish AddOrder")
}
