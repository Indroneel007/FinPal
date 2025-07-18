package api

import (
	"database/sql"
	db "examples/SimpleBankProject/db/sqlc"
	"examples/SimpleBankProject/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateLocationRequest struct {
	Address   string  `json:"address" binding:"required"`
	Lattitude float64 `json:"lattitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

func (s *Server) createLocation(c *gin.Context) {
	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, NewError(err))
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization payload not found"})
		return
	}
	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid authorization payload"})
		return
	}

	err := s.store.CreateLocation(c, db.CreateLocationParams{
		Username:  payload.Username,
		Address:   req.Address,
		Latitude:  sql.NullFloat64{Float64: req.Lattitude, Valid: true},
		Longitude: sql.NullFloat64{Float64: req.Longitude, Valid: true},
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location created successfully"})
}

func (s *Server) getLocation(c *gin.Context) {
	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization payload not found"})
		return
	}
	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid authorization payload"})
		return
	}

	location, err := s.store.GetLocationByUsername(c, payload.Username)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, location)
}
