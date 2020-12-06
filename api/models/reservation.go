package models

type ReservationPostResponse struct {
	IsSuccess     bool   `json:"is_success"`
	TransactionID int64  `json:"transaction_id"`
	CellColumn    string `json:"cell_column"`
	CellStart     int    `json:"cell_start"`
	CellEnd       int    `json:"cell_end"`
	Lecture       string `json:"lecture"`
	Professor     string `json:"professor"`
}
