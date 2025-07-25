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
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

type Server struct {
	store      *db.Store
	tokenMaker util.Maker
	router     *gin.Engine
	redis      *redis.Client
	geminiKey  string
}

func NewServer(store *db.Store, redisClient *redis.Client) (*Server, error) {
	/*err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}*/

	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			break
		}
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

	geminiAI := viper.GetString("GEMINI_API_KEY")
	if geminiAI == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set in the environment variables")
	}

	ginRouter := gin.Default()

	// CORS middleware for development - allow all origins
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	ginRouter.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
	})

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		router:     ginRouter,
		redis:      redisClient,
		geminiKey:  geminiAI,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("accountType", validType)
	}

	return server, nil
}

func (s *Server) MountHandlers() {
	public := s.router.Group("/")
	auth := s.router.Group("/").Use(AuthMiddleware(s.tokenMaker))

	// Public routes (no auth)
	public.POST("/register", s.createUser)
	public.GET("/tests", s.TestRoute)
	public.POST("/forgotpassword", s.forgotPassword)
	public.POST("/login", s.loginUser)
	public.GET("/users/:username", s.getOtherUser)
	// Authenticated routes
	auth.POST("/accounts", s.createAccount)
	auth.GET("/accounts", s.listAccounts)
	auth.GET("/accounts/:id", s.getAccount)
	auth.GET("/user/:username", s.getUser)
	auth.POST("/transfers", s.createTransfer)
	auth.POST("/checkotp", s.checkOtp)
	auth.POST("/resetpassword", s.resetPassword)
	auth.GET("/accountsbytype", s.getAccountListByOwnerAndType)
	auth.GET("/transfers/:username", s.getTransferBetweenUser)
	auth.GET("/transfers/user", s.listTransactedUsersWithTotals)
	//Groups Routes
	auth.POST("/groups", s.createGroup)
	auth.GET("/groups", s.listGroups)
	auth.GET("/groups/:id", s.getGroup)
	auth.POST("/groups/:id/add", s.addMemberToGroup)
	auth.GET("/groups/:id/accounts", s.getGroupMembers)
	auth.POST("/groups/:id/updatename", s.updateGroupName)
	auth.POST("/groups/:id/leave", s.leaveGroup)
	auth.POST("/groups/:id/delete", s.deleteGroup)
	//Location
	auth.POST("/location", s.createLocation)
	//auth.POST("/location/:id", s.updateLocation)
	auth.GET("/location", s.getLocation)

	//prompt
	auth.POST("/prompt", s.PromptAPI)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
