package domain

import (
	"time"
)

type Merch struct {
	Id    int    `json:"-" db:"id"`
	Name  string `json:"name" binding:"required"`
	Price string `json:"price" binding:"required"`
}

type Transactions struct {
	Id          int        `json:"-" db:"id"`
	Source      *int       `json:"source"`
	Destination string     `json:"destination" binding:"required"`
	Amount      int        `json:"amount" binding:"required"`
	Timestamp   *time.Time `json:"timestamp"`
}

type UserSummary struct {
	UserName            string              `json:"username"`
	Coins               int                 `json:"coins"`
	PurchasedItems      []PurchasedItem     `json:"purchased_items"`
	TransactionsSummary TransactionsSummary `json:"transactions_summary"`
}

type PurchasedItem struct {
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
}
type TransactionsSummary struct {
	ReceivedCoins []Transactions `json:"received_coins"`
	SentCoins     []Transactions `json:"sent_coins"`
}
