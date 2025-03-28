package models

type Stock struct {
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string `json:"company"`
	StockID string `json:"stockid"`
}
