package transaction

import (
	"time"
)

type TransactionFormat struct {
	ID        int                       `json:"id"`
	Amount    int                       `json:"amount"`
	Status    string                    `json:"status"`
	Code      string                    `json:"code"`
	Campaign  CampaignTransactionFormat `json:"campaign"`
	User      UserTransactionFormat     `json:"user"`
	CreatedAt time.Time                 `json:"created_at"`
}

type CampaignTransactionFormat struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Highlight string `json:"highlight"`
}

type UserTransactionFormat struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func FormatTransaction(transaction Transaction) TransactionFormat {
	return TransactionFormat{
		ID:     transaction.ID,
		Amount: transaction.Amount,
		Status: transaction.Status,
		Code:   transaction.Code,
		Campaign: CampaignTransactionFormat{
			ID:        transaction.Campaign.ID,
			Name:      transaction.Campaign.Name,
			Highlight: transaction.Campaign.Highlight,
		},
		User: UserTransactionFormat{
			ID:    transaction.User.ID,
			Name:  transaction.User.Name,
			Email: transaction.User.Email,
		},
		CreatedAt: transaction.CreatedAt,
	}
}
