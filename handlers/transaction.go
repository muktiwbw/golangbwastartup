package handlers

import (
	"bwastartup/entities/campaign"
	"bwastartup/entities/transaction"
	"bwastartup/entities/user"
	"bwastartup/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
	campaignService    campaign.Service
}

func NewTransactionHandler(transactionService transaction.Service, campaignService campaign.Service) *transactionHandler {
	return &transactionHandler{transactionService, campaignService}
}

func (h transactionHandler) CreateTransaction(c *gin.Context) {
	var campaignUri campaign.GetCampaignByIDInput
	var input transaction.TransactionInput

	err := c.ShouldBindUri(&campaignUri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input campaign ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundCampaign, err := h.campaignService.GetCampaignByID(campaignUri.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundCampaign.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Not found", http.StatusNotFound, "success", nil))

		return
	}

	err = c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input field", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	authUser := c.MustGet("authUser").(user.User)

	input.CampaignID = foundCampaign.ID
	input.UserID = authUser.ID

	createdTransaction, err := h.transactionService.CreateTransaction(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	// Calculate and update campaign stats
	currentAmount, backerCount, err := h.transactionService.GetNewCampaignStats(foundCampaign.ID)

	_, err = h.campaignService.UpdateCampaign(foundCampaign, campaign.UpdateCampaignInput{CurrentAmount: currentAmount, BackersCount: backerCount})

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	c.JSON(http.StatusCreated, helpers.APIResponse("Successfully created a transaction", http.StatusCreated, "created", transaction.FormatTransaction(createdTransaction)))

}

func (h transactionHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions()

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	formattedTransactions := []transaction.TransactionFormat{}

	for _, trx := range transactions {
		formattedTransactions = append(formattedTransactions, transaction.FormatTransaction(trx))
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", formattedTransactions))
}

func (h transactionHandler) GetOwnTransactions(c *gin.Context) {
	authUser := c.MustGet("authUser").(user.User)

	transactions, err := h.transactionService.GetAllTransactionsByRef(authUser.ID, "user_id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	formattedTransactions := []transaction.TransactionFormat{}

	for _, trx := range transactions {
		formattedTransactions = append(formattedTransactions, transaction.FormatTransaction(trx))
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", formattedTransactions))
}

func (h transactionHandler) GetTransactionByCampaignID(c *gin.Context) {
	var campaignUri campaign.GetCampaignByIDInput

	err := c.ShouldBindUri(&campaignUri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input campaign ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	transactions, err := h.transactionService.GetAllTransactionsByRef(campaignUri.ID, "campaign_id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	formattedTransactions := []transaction.TransactionFormat{}

	for _, trx := range transactions {
		formattedTransactions = append(formattedTransactions, transaction.FormatTransaction(trx))
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", formattedTransactions))
}

func (h transactionHandler) GetTransactionByID(c *gin.Context) {
	var transactionUri transaction.GetTransactionByIDInput

	err := c.ShouldBindUri(&transactionUri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input transaction ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundTransaction, err := h.transactionService.GetTransactionByID(transactionUri.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundTransaction.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Transaksi tidak ditemukan", http.StatusNotFound, "not-found", nil))

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", transaction.FormatTransaction(foundTransaction)))
}

func (h transactionHandler) VerifyTransaction(c *gin.Context) {
	var transactionUri transaction.GetTransactionByIDInput

	err := c.ShouldBindUri(&transactionUri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input transaction ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundTransaction, err := h.transactionService.GetTransactionByID(transactionUri.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundTransaction.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Transaksi tidak ditemukan", http.StatusNotFound, "not-found", nil))

		return
	}

	verifiedTransaction, err := h.transactionService.VerifyTransaction(foundTransaction)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	c.JSON(http.StatusCreated, helpers.APIResponse("Successfully verified transaction", http.StatusCreated, "updated", transaction.FormatTransaction(verifiedTransaction)))
}

func (h *transactionHandler) GetNewCampaignStats(c *gin.Context) {
	var campaignInput campaign.GetCampaignByIDInput

	err := c.ShouldBindUri(&campaignInput)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input campaign ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	_, _, err = h.transactionService.GetNewCampaignStats(campaignInput.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}
}
