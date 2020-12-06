package models

type CellTransactionModel struct {
	TransactionID int
	UserID        int
	UserEmail     string
	Lecture       string
	CellColumn    string
	CellStart     int
	CellEnd       int
	Professor     string
	CreatedAt     string
	Capacity      int
}

type CellGetResponse struct {
	Cells      []CellItem `json:"cells"`
	CellsCount int        `json:"cells_count"`
}

type CellItem struct {
	Cell          string `json:"cell"`
	IsReserved    bool   `json:"is_reserved"`
	UserEmail     string `json:"user_email"`
	UserID        int    `json:"user_id"`
	Lecture       string `json:"lecture"`
	Professor     string `json:"professor"`
	Capacity      int    `json:"capacity"`
	TransactionID int    `json:"transaction_id"`
	CreatedAt     string `json:"created_at"`
}
