package api

import (
	db "examples/SimpleBankProject/db/sqlc"
	"examples/SimpleBankProject/util"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
)

type userRegisterRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type getUserRequest struct {
	username string `uri:"username" binding:"required,alphanum,min=3,max=20"`
}

type userLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func UserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(c *gin.Context) {
	var req userRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	password, err := util.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	p := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(c, p)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	res := UserResponse(user)

	c.JSON(http.StatusOK, res)
}

func (s *Server) getUser(c *gin.Context) {
	var req getUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	username := req.username

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

	if payload.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to access this user"})
		return
	}

	user, err := s.store.GetUser(c, username)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) loginUser(c *gin.Context) {
	/*err := godotenv.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}*/
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	var err error

	viper.AutomaticEnv()
	secret := viper.GetString("TOKEN_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN_SECRET is not set in the environment variables"})
		return
	}
	duration := viper.GetDuration("TOKEN_DURATION")
	if duration == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN_DURATION is not set in the environment variables"})
		return
	}

	var req userLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}
	var field string
	field = "username"

	if strings.Contains(req.Username, "@") {
		field = "email"
	}

	var user db.User
	//var err error

	if field == "username" {
		user, err = s.store.GetUser(c, req.Username)
	} else {
		user, err = s.store.GetUserByEmail(c, req.Username)
	}

	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	if !util.CheckPasswordHash(req.Password, user.HashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	//var pasetoMaker = util.PasetoMaker{}
	var maker = &util.PasetoMaker{
		Paseto:       paseto.NewV2(),
		SymmetricKey: []byte(secret),
	}

	token, err := maker.CreateToken(user.Username, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	res := loginUserResponse{
		AccessToken: token,
		User:        UserResponse(user),
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) forgotPassword(c *gin.Context) {
	fmt.Println("Inside forgotPassword handler")
	if s.redis == nil {
		log.Println(">>> Redis client is nil!")
		c.JSON(500, gin.H{"error": "internal server error: redis is nil"})
		return
	}

	var req forgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	if s.redis == nil {
		c.JSON(500, gin.H{"error": "Redis is not configured"})
		return
	}

	log.Printf("Redis client: %+v", s.redis)

	user, err := s.store.GetUserByEmail(c, req.Email)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email address"})
		return
	}

	log.Println("Before generating OTP")
	otp := util.GenerateOTP()

	log.Println("Before storing OTP in Redis")
	err = util.AddOTPToRedis(otp, user.Email, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	log.Println("Before sending OTP email")
	err = util.SendOTPEmail(otp, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to " + user.Email})

	//otp := "123456" // dummy otp for now
	/*err = s.redis.Set(c, "otp:"+user.Email, otp, 10*time.Minute).Err()
	if err != nil {
		log.Println(">>> Redis SET error:", err)
		c.JSON(500, gin.H{"error": "redis error"})
		return
	}

	c.JSON(200, gin.H{"message": "OTP sent"})*/
}
