package transaction

import (
	"time"
)

type Service interface {
	CreateTransaction(transactionInput TransactionInput) (Transaction, error)
	GetAllTransactions() ([]Transaction, error)
	GetAllTransactionsByRef(id int, field string) ([]Transaction, error)
	GetTransactionByID(id int) (Transaction, error)
	VerifyTransaction(transaction Transaction) (Transaction, error)
	GetNewCampaignStats(campaignID int) (currentAmount int, backerCount int64, err error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateTransaction(transactionInput TransactionInput) (Transaction, error) {
	transaction := Transaction{
		CampaignID: transactionInput.CampaignID,
		UserID:     transactionInput.UserID,
		Amount:     transactionInput.Amount,
		Status:     "pending",
		Code:       "randomString",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return createdTransaction, err
	}

	return createdTransaction, nil
}

func (s *service) GetAllTransactions() ([]Transaction, error) {
	transactions, err := s.repository.All()

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetAllTransactionsByRef(id int, field string) ([]Transaction, error) {
	transactions, err := s.repository.AllByRef(id, field)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByID(id int) (Transaction, error) {
	foundTransaction, err := s.repository.Get(id)

	if err != nil {
		return foundTransaction, err
	}

	return foundTransaction, nil
}

func (s *service) VerifyTransaction(transaction Transaction) (Transaction, error) {
	transaction.Status = "paid"
	transaction.UpdatedAt = time.Now()

	verifiedTransaction, err := s.repository.Verify(transaction)

	if err != nil {
		return verifiedTransaction, err
	}

	return verifiedTransaction, nil
}

func (s *service) GetNewCampaignStats(campaignID int) (currentAmount int, backerCount int64, err error) {
	currentAmount, backerCount, err = s.repository.CalculateCampaignStats(campaignID)

	if err != nil {
		return currentAmount, backerCount, err
	}

	return currentAmount, backerCount, nil
}
