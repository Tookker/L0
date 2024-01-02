package ordercachemap_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/zap"

	"L0/internal/caches/ordercache/ordercachemap"
	"L0/internal/models"
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

func TestAddOrder(t *testing.T) {
	logger, _ := getZapLogger()
	cache := ordercachemap.NewCache(logger)

	type args struct {
		id    uint
		order models.Order
	}

	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "Add not existing order",
			args: args{
				id:    ID,
				order: orderVar,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Add existing order",
			args: args{
				id:    ID,
				order: orderVar,
			},
			want:    errors.New(ordercachemap.OrderExistsErr),
			wantErr: true,
		},
	}

	t.Logf("Start AddOrder")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cache.AddOrder(tt.args.id, &tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrder error = %v, whantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v pass!", tt.name)
		})
	}

	t.Logf("Finish AddOrder")
}

func TestGetOrder(t *testing.T) {
	logger, _ := getZapLogger()
	cache := ordercachemap.NewCache(logger)

	type args struct {
		ctx context.Context
		id  uint
	}

	type want struct {
		order models.Order
		err   error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "Get existing order",
			args: args{
				id:  ID,
				ctx: context.Background(),
			},
			want: want{
				order: orderVar,
				err:   nil,
			},
			wantErr: false,
		},
		{
			name: "Get not existing order",
			args: args{
				id:  ID,
				ctx: context.Background(),
			},
			want: want{
				order: models.Order{},
				err:   errors.New(ordercachemap.OrderNotFindErr),
			},
			wantErr: true,
		},
	}

	cache.AddOrder(ID, &orderVar)

	t.Logf("Start GetOrder")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cache.GetOrder(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrder error = %v, whantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v pass!", tt.name)
		})
	}

	t.Logf("Finish GetOrder")
}
