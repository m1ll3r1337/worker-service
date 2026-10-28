package entity

import "github.com/rs/zerolog"

////////////////////////////////////////////////////////////////////////////////
///// EVENT MODEL //////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type EventOrderCreated struct {
	OrderGUID  string                  `json:"order_guid"`
	UserGUID   *string                 `json:"user_guid,omitempty"`
	Currency   string                  `json:"currency"`
	TotalPrice int64                   `json:"total_price"`
	Items      []EventOrderCreatedItem `json:"items"`
	CreatedAt  string                  `json:"created_at"`
}

type EventOrderCreatedItem struct {
	ProductGUID string `json:"product_guid"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
}

func (e EventOrderCreated) MarshalZerologObject(ev *zerolog.Event) {
	ev.Str("order_guid", e.OrderGUID)
	ev.Str("currency", e.Currency)
	ev.Int64("total_price", e.TotalPrice)
	ev.Int("items_count", len(e.Items))
}
