package api

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
    db "github.com/jingtian168/admin-api-go/db/sqlc"
	"github.com/jingtian168/admin-api-go/token"
	"github.com/jingtian168/admin-api-go/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	if err := InitTrans("zh"); err != nil {
		log.Fatalf("初始化验证翻译器错误: %s", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(Cors())
	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/users/renew_access", server.renewAccessToken)
	

	authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouters.GET("/getUserInfo", server.getUserInfo)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
