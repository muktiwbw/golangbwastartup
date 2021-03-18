package transaction

type GetTransactionByIDInput struct {
	ID int `uri:"transaction_id" binding:"required"`
}

type TransactionInput struct {
	CampaignID int
	UserID     int
	Amount     int `json:"amount" binding:"required"`
}
