package service

type StockOrder struct {
	ID         uint32 `json:"id"`
	Sid        uint32 `json:"sid"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
}
