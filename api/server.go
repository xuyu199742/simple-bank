package api

import (
	"github.com/gin-gonic/gin"
	db "simple-bank/db/sqlc"
)

const (
	serverAdd   = "0.0.0.0:8080"
	successCode = iota
	errorCode
)

type (
	Server struct {
		db     *db.Store
		router *gin.Engine
	}

	Response struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
		Msg  string      `json:"msg"`
	}
)

func NewServer(db *db.Store) *Server {
	server := &Server{
		db:     db,
		router: gin.Default(),
	}
	server.setUpRouter()
	server.router.SetTrustedProxies(nil)
	return server
}

func (s *Server) setUpRouter() {
	s.router.POST("/account", s.CreateAccount)
	s.router.POST("/account/:id", s.GetAccount)
	s.router.GET("/accounts", s.ListAccount)
}

func (s *Server) resOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code: successCode,
		Data: data,
		Msg:  "success",
	})
}

func (s *Server) resFail(ctx *gin.Context, err error, code int) {
	ctx.JSON(code, Response{
		Code: errorCode,
		Msg:  err.Error(),
	})
}

func (s *Server) Star() error {
	return s.router.Run(serverAdd)
}
