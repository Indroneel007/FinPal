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
	Salary   int64  `json:"salary" binding:"required"`
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
	Salary            int64     `json:"salary"`
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

type checkOtpRequest struct {
	OTP string `json:"otp" binding:"required"`
}

type resetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

type getOtherUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum,min=3,max=20"`
}

type getOtherUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func OtherUserResponse(user db.User) getOtherUserResponse {
	return getOtherUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

func UserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		Salary:            user.Salary,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(c *gin.Context) {
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
		Salary:         req.Salary,
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

func (s *Server) getOtherUser(c *gin.Context) {
	var req getOtherUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		log.Printf("Error binding URI: %v", err)
		c.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	log.Printf("Looking up user: %s", req.Username)

	user, err := s.store.GetUser(c, req.Username)
	if err != nil {
		log.Printf("Error getting user %s: %v", req.Username, err)
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	log.Printf("Found user: %+v", user)
	c.JSON(http.StatusOK, OtherUserResponse(user))
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
	fmt.Println("OTP sent successfully")
}

func (s *Server) checkOtp(c *gin.Context) {
	var req checkOtpRequest
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

	user, err := s.store.GetUser(c, payload.Username)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	otpKey := util.OtpKeyPrefix + user.Email
	/*hashedOTP, err := util.HashPassword(req.OTP)
	if err != nil {
		log.Printf("Error hashing OTP: %v", err)
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}*/

	// Check if the OTP exists in Redis
	otpStored, err := s.redis.Get(c, otpKey).Result()
	if err != nil {
		if s.redis == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "OTP not found"})
			return
		}
		log.Printf("Error getting OTP from Redis: %v", err)
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	ok = util.CheckPasswordHash(req.OTP, otpStored)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// If OTP is valid, delete it from Redis
	err = s.redis.Del(c, otpKey).Err()
	if err != nil {
		log.Printf("Error deleting OTP from Redis: %v", err)
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	res := userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) resetPassword(c *gin.Context) {
	var req resetPasswordRequest
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

	user, err := s.store.GetUser(c, payload.Username)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	hashedNewPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	if user.HashedPassword == hashedNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password cannot be the same as the old password"})
		return
	}

	updatedUser, err := s.store.UpdateUserPassword(c, db.UpdateUserPasswordParams{
		Username:       user.Username,
		HashedPassword: hashedNewPassword,
	})
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}

	res := userResponse{
		Username:          updatedUser.Username,
		FullName:          updatedUser.FullName,
		Email:             updatedUser.Email,
		PasswordChangedAt: updatedUser.PasswordChangedAt,
		CreatedAt:         updatedUser.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}
