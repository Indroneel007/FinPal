package api

import (
	db "examples/SimpleBankProject/db/sqlc"
	"fmt"

	//"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Server struct {
	store      *db.Store
	tokenMaker Maker
	router     *gin.Engine
}

func NewServer(store *db.Store) (*Server, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	viper.AutomaticEnv()
	secret := viper.GetString("TOKEN_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("TOKEN_SECRET is not set in the environment variables")
	}

	tokenMaker, err := NewPasetoMaker(secret)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

func (s *Server) MountHandlers() {
	api := s.router.Group("/api")

	api.POST("/accounts", s.createAccount)
	api.GET("/accounts/:id", s.getAccount)
	api.GET("/accounts", s.listAccounts)
	api.POST("/transfers", s.createTransfer)

	api.GET("/tests", s.TestRoute)

	api.POST("/register", s.createUser)
	api.GET("/user/:username", s.getUser)
	api.POST("/login", s.loginUser)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
