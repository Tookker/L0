package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/stan.go"

	"L0/internal/config"
	"L0/internal/models"
)

func main() {
	sendingOrders := []models.Order{
		{
			UID:      "n5634327b2b84b6test",
			TrackNum: "WBILMTESTTRACK",
			Entry:    "WBIL",
			Delivery: models.Delivery{
				Name:    "Alim Testov",
				Phone:   "+9740000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 11",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: models.Payment{
				Transaction:  "b563f2357b2b84b6test",
				RequestID:    "",
				Currency:     "Bel RUB",
				Provider:     "wbpay",
				Amount:       1827,
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
					RID:        "ab421ads87a764ae0btest",
					Name:       "Wheel",
					Sale:       30,
					Size:       "0",
					TotalPrice: 317,
					NmID:       2389212,
					Brand:      "Kama",
					Status:     202,
				},
			},
			Locale:          "en",
			InternSign:      "",
			CustomerID:      "admin",
			DeliveryService: "meest",
			ShardKey:        "9",
			SmID:            9,
			CreationDate:    "2023-11-26T06:22:19Z",
			CofShard:        "1",
		},
		{
			UID:      "543534327b2b84b6test",
			TrackNum: "WBILMTESTTRACK",
			Entry:    "WBIL",
			Delivery: models.Delivery{
				Name:    "Root Testov",
				Phone:   "+9740003200",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 17",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: models.Payment{
				Transaction:  "b563f2357b2b84b6test",
				RequestID:    "",
				Currency:     "Tenge",
				Provider:     "wbpay",
				Amount:       1827,
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
					RID:        "ab421ads87a764ae0btest",
					Name:       "Monitor",
					Sale:       30,
					Size:       "0",
					TotalPrice: 317,
					NmID:       2389212,
					Brand:      "BENQ",
					Status:     202,
				},
			},
			Locale:          "en",
			InternSign:      "",
			CustomerID:      "admin",
			DeliveryService: "meest",
			ShardKey:        "9",
			SmID:            9,
			CreationDate:    "2023-11-22T06:22:19Z",
			CofShard:        "1",
		},
	}

	config, err := config.LoadConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	sc, err := stan.Connect(config.ClusterID, config.PublisherID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatalln(err)
	}
	defer sc.Close()

	for i := 0; i < len(sendingOrders); i++ {
		data, err := json.Marshal(sendingOrders[i])
		if err != nil {
			log.Fatalln(err)
		}

		err = sc.Publish(config.Subject, data)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Published order with order_uid=%s\n", sendingOrders[i].UID)
	}

	log.Printf("Send error msg")
	err = sc.Publish(config.Subject, []byte("error msg sending"))
	if err != nil {
		log.Fatalln(err)
	}
}
