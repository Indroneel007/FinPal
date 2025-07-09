package api

import (
	//"errors"
	db "examples/SimpleBankProject/db/sqlc"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance" binding:"required,min=0"`     // Ensure balance is non-negative
	Currency string `json:"currency" binding:"required,currency"` // Use custom validator for currency
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=5"`
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"` // Use custom validator for currency
}

func (s *Server) createAccount(c *gin.Context) {
	var req accountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	p := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  req.Balance,
		Currency: req.Currency,
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

	c.JSON(http.StatusOK, account)
}

func (s *Server) listAccounts(c *gin.Context) {
	var req listAccountsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	accounts, err := s.store.ListAccounts(c, db.ListAccountsParams{
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

	if !s.ValidAccountCurrency(c, req.FromAccountID, req.Currency) {
		return
	}

	if !s.ValidAccountCurrency(c, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
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

func (s *Server) ValidAccountCurrency(c *gin.Context, accountID int64, currency string) bool {
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

	return true
}
