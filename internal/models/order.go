package models

type Order struct {
	UID             string   `json:"order_uid"`
	TrackNum        string   `json:"track_number"`
	Entry           string   `json:"entry"`
	Delivery        Delivery `json:"delivery"`
	Payment         Payment  `json:"payment"`
	Items           []Item   `json:"items"`
	Locale          string   `json:"locale"`
	InternSign      string   `json:"internal_signature"`
	CustomerID      string   `json:"customer_id"`
	DeliveryService string   `json:"delivery_service"`
	ShardKey        string   `json:"shardkey"`
	SmID            int      `json:"sm_id"`
	CreationDate    string   `json:"date_created"`
	CofShard        string   `json:"oof_shard"`
}
