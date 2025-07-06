package api

import (
	db "examples/SimpleBankProject/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	return server
}

func (s *Server) MountHandlers() {
	api := s.router.Group("/api")

	api.POST("/accounts", s.createAccount)
	api.GET("/tests", s.TestRoute)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
