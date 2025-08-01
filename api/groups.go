package api

import (
	"database/sql"
	db "examples/SimpleBankProject/db/sqlc"
	"examples/SimpleBankProject/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateGroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
	Currency  string `json:"currency" binding:"required,currency"`
	Type      string `json:"type" binding:"required,accountType"`
}

type ListGroupsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=5"`
}

type GetGroupRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type AddMemberToGroupRequest struct {
	Username string `json:"username" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
	Type     string `json:"type" binding:"required,accountType"`
}

type GetGroupMembersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=20"`
}

type UpdateGroupNameRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

type getGroupHistoryRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) createGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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

	arg, err := s.store.CreateGroupTx(c, db.CreateGroupTxParams{
		Username:  payload.Username,
		GroupName: req.GroupName,
		Currency:  req.Currency,
		Type:      req.Type,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, arg)
}

func (s *Server) listGroups(c *gin.Context) {
	var req ListGroupsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	result, err := s.store.ListGroupsByUser(c, db.ListGroupsByUserParams{
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

	c.JSON(http.StatusOK, result)
}

func (s *Server) getGroup(c *gin.Context) {
	var req GetGroupRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	account, err := s.store.GetAccountByGroupIDAndOwner(c, db.GetAccountByGroupIDAndOwnerParams{
		GroupID: sql.NullInt64{Int64: req.ID, Valid: true},
		Owner:   payload.Username,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	if account.Owner != payload.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this group"})
		return
	}

	c.JSON(http.StatusOK, account)
}

type GroupIDUri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) addMemberToGroup(c *gin.Context) {
	var req AddMemberToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var uri GroupIDUri
	if err := c.ShouldBindUri((&uri)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID := uri.ID

	args := db.GetAccountByOwnerCurrencyTypeGroupIDParams{
		Owner:    req.Username,
		Currency: req.Currency,
		Type:     req.Type,
		GroupID:  sql.NullInt64{Int64: ID, Valid: true},
	}

	newAccount, err := s.store.GetAccountByOwnerCurrencyTypeGroupID(c, args)
	if err == sql.ErrNoRows {
		newAccount, err := s.store.CreateAccountWithGroup(c, db.CreateAccountWithGroupParams{
			Owner:    req.Username,
			Currency: req.Currency,
			Type:     req.Type,
			GroupID:  sql.NullInt64{Int64: ID, Valid: true},
		})
		if err != nil {
			if apiErr := convertToApiErr(err); apiErr != nil {
				c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
				return
			}
			c.JSON(http.StatusInternalServerError, NewError(err))
			return
		}
		c.JSON(http.StatusOK, newAccount)
		return
	} else if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	/*newAccount, err := s.store.CreateAccountWithGroup(c, db.CreateAccountWithGroupParams{
		Owner:    req.Username,
		Balance:  0,
		Currency: req.Currency,
		Type:     req.Type,
		GroupID:  sql.NullInt64{Int64: ID, Valid: true},
	})*/

	c.JSON(http.StatusOK, newAccount)
}

func (s *Server) getGroupMembers(c *gin.Context) {
	var req GetGroupMembersRequest
	var ID int64
	if err := c.ShouldBindUri(&ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	arg, err := s.store.GetAccountByGroupIDAndOwner(c, db.GetAccountByGroupIDAndOwnerParams{
		GroupID: sql.NullInt64{Int64: ID, Valid: true},
		Owner:   payload.Username,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	if arg.Owner != payload.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this group's members"})
		return
	}

	members, err := s.store.GetGroupMembers(c, db.GetGroupMembersParams{
		ID:     ID,
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

	c.JSON(http.StatusOK, members)
}

type UpdateNameGroupIDUri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) updateGroupName(c *gin.Context) {
	var req UpdateGroupNameRequest
	var uri UpdateNameGroupIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID := uri.ID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	_, err := s.store.GetAccountByGroupIDAndOwner(c, db.GetAccountByGroupIDAndOwnerParams{
		GroupID: sql.NullInt64{Int64: ID, Valid: true},
		Owner:   payload.Username,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	group, err := s.store.UpdateGroupName(c, db.UpdateGroupNameParams{
		ID:        ID,
		GroupName: req.NewName,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, group)
}

type LeaveGroupIDUri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) leaveGroup(c *gin.Context) {
	var uri LeaveGroupIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID := uri.ID

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

	account, err := s.store.GetAccountByGroupIDAndOwner(c, db.GetAccountByGroupIDAndOwnerParams{
		GroupID: sql.NullInt64{Int64: ID, Valid: true},
		Owner:   payload.Username,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	updatedAccount, err := s.store.UpdateAccountGroup(c, account.ID)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, updatedAccount)
}

type DeleteGroupIDUri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) deleteGroup(c *gin.Context) {
	var uri DeleteGroupIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID := uri.ID

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

	_, err := s.store.GetAccountByGroupIDAndOwner(c, db.GetAccountByGroupIDAndOwnerParams{
		GroupID: sql.NullInt64{Int64: ID, Valid: true},
		Owner:   payload.Username,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	err = s.store.DeleteGroup(c, ID)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted group successfully"})
}

func (s *Server) getGroupHistory(c *gin.Context) {
	var req getGroupHistoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupHistory, err := s.store.GetGroupTransactionHistory(c, sql.NullInt64{Int64: req.ID, Valid: true})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, groupHistory)
}
