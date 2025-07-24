package api

import (
	//"errors"
	db "examples/SimpleBankProject/db/sqlc"
	"examples/SimpleBankProject/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountRequest struct {
	Currency string `json:"currency" binding:"required,currency"` // Use custom validator for currency
	Type     string `json:"type" binding:"required,accountType"`  // Use custom validator for type
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=5"`
}

type transferRequest struct {
	ToUsername string `json:"to_username" binding:"required"`
	Amount     int64  `json:"amount" binding:"required,min=1"`
	Currency   string `json:"currency" binding:"required,currency"` // Use custom validator for currency
	Type       string `json:"type" binding:"required,accountType"`  // Use custom validator for type
}

type getAccountListByOwnerAndTypeRequest struct {
	Type     string `form:"type" binding:"required,accountType"` // Use custom validator for type
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=5"`
}

type getTransferBetweenUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

type getTransferBetweenUserResponse struct {
	Paid     []db.ListTransfersBetweenAccountsRow `json:"paid"`
	Received []db.ListTransfersBetweenAccountsRow `json:"received"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req accountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
		return
	}

	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
		return
	}

	p := db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
		Type:     req.Type,
	}

	account, err := s.store.CreateAccount(c, p)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) TestRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Test route is working!"})
}

func (s *Server) getAccount(c *gin.Context) {
	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	id := req.ID

	account, err := s.store.GetAccount(c, id)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
		return
	}

	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
		return
	}

	if account.Owner != payload.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) listAccounts(c *gin.Context) {
	var req listAccountsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
		return
	}

	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
		return
	}

	accounts, err := s.store.ListAccountsByOwner(c, db.ListAccountsByOwnerParams{
		Owner:  payload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (s *Server) createTransfer(c *gin.Context) {
	var req transferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	payloadInterface, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	payload, ok := payloadInterface.(*util.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization payload"})
		return
	}

	arg := db.TransferTxParams{
		FromUsername: payload.Username,
		ToUsername:   req.ToUsername,
		Currency:     req.Currency,
		Type:         req.Type,
		Amount:       req.Amount,
	}

	result, err := s.store.TransferTx(c, arg)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) ValidAccountCurrencyAndType(c *gin.Context, accountID int64, currency string, type1 string) bool {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return false
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("currency mismatch: account currency is %s, but request currency is %s", account.Currency, currency)
		c.JSON(http.StatusBadRequest, NewError(err))
		return false
	}

	if account.Type != type1 {
		err = fmt.Errorf("type mismatch: account type is %s, but request type is %s", account.Type, type1)
		c.JSON(http.StatusBadRequest, NewError(err))
		return false
	}

	return true
}

func (s *Server) getAccountListByOwnerAndType(c *gin.Context) {
	var req getAccountListByOwnerAndTypeRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
		return
	}

	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
		return
	}

	accounts, err := s.store.GetAccountListByOwnerAndType(c, db.GetAccountListByOwnerAndTypeParams{
		Owner:  payload.Username,
		Type:   req.Type,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (s *Server) getTransferBetweenUser(c *gin.Context) {
	var req getTransferBetweenUserRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
		return
	}

	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
		return
	}

	myAccounts, err := s.store.GetAccountsByUser(c, payload.Username)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	otherAccounts, err := s.store.GetAccountsByUser(c, req.Username)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	if len(myAccounts) == 0 || len(otherAccounts) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	Paid, err := s.store.ListTransfersBetweenAccounts(c, db.ListTransfersBetweenAccountsParams{
		FromAccountID: myAccounts,
		ToAccountID:   otherAccounts,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	Received, err := s.store.ListTransfersBetweenAccounts(c, db.ListTransfersBetweenAccountsParams{
		FromAccountID: otherAccounts,
		ToAccountID:   myAccounts,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	response := getTransferBetweenUserResponse{
		Paid:     Paid,
		Received: Received,
	}

	c.JSON(http.StatusOK, response)
}
