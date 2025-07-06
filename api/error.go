package api

import (
	"errors"

	//"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrAccessForbidden      = errors.New("access forbidden")
	ErrUserNotFound         = errors.New("user not found")
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

/*func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}*/

func NewError(err error) *Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = err.Error()
	return &e
}

func NewValidationError(err error) *Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = err.Error()
	return &e
}

func convertToApiErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "users_username_key":
			return ErrUsernameAlreadyTaken
		case "users_email_key":
			return ErrEmailAlreadyTaken
		}
	}
	return nil
}
