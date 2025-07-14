package api

import (
	db "examples/SimpleBankProject/db/sqlc"
	"examples/SimpleBankProject/util"
	"fmt"

	//"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Server struct {
	store      *db.Store
	tokenMaker util.Maker
	router     *gin.Engine
	redis      *redis.Client
}

func NewServer(store *db.Store, redisClient *redis.Client) (*Server, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	viper.AutomaticEnv()
	secret := viper.GetString("TOKEN_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("TOKEN_SECRET is not set in the environment variables")
	}

	tokenMaker, err := util.NewPasetoMaker(secret)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
		redis:      redisClient,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

func (s *Server) MountHandlers() {
	public := s.router.Group("/")
	auth := s.router.Group("/").Use(AuthMiddleware(s.tokenMaker))

	// ✅ Public routes (no auth)
	public.POST("/register", s.createUser)
	public.GET("/tests", s.TestRoute)
	public.POST("/forgotpassword", s.forgotPassword)
	public.POST("/login", s.loginUser)

	// ✅ Authenticated routes
	auth.POST("/accounts", s.createAccount)
	auth.GET("/accounts", s.listAccounts)
	auth.GET("/accounts/:id", s.getAccount)
	auth.GET("/user/:username", s.getUser)
	auth.POST("/transfers", s.createTransfer)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
