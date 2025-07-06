package api

import (
	db "examples/SimpleBankProject/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance" binding:"required,min=0"` // Ensure balance is non-negative
	Currency string `json:"currency" binding:"required,oneof=USD Euros Rupees"`
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
