package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "simple-bank/db/sqlc"
)

type (
	createAccountReq struct {
		Owner    string `json:"owner" binding:"required" `
		Currency string `json:"currency" binding:"required,oneof=USD RMB"`
	}
	getAccountReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	listAccountReq struct {
		PageSize int32 `form:"page_size"  binding:"required,min=5,max=10"`
		Page     int32 `form:"page"  binding:"required,min=1"`
	}
)

func (s *Server) CreateAccount(ctx *gin.Context) {

	var arg = createAccountReq{}
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		s.resFail(ctx, err, 400)
		return
	}

	account, err := s.db.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    arg.Owner,
		Currency: arg.Currency,
	})
	if err != nil {
		s.resFail(ctx, err, 500)
		return
	}
	s.resOk(ctx, account)
}

func (s *Server) GetAccount(ctx *gin.Context) {

	var arg = getAccountReq{}
	if err := ctx.ShouldBindUri(&arg); err != nil {
		s.resFail(ctx, err, 400)
		return
	}

	account, err := s.db.GetAccount(ctx, arg.ID)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			s.resFail(ctx, err, 404)
			return
		}
		s.resFail(ctx, err, 500)
		return
	}
	s.resOk(ctx, account)
}

func (s *Server) ListAccount(ctx *gin.Context) {

	var arg = listAccountReq{}
	if err := ctx.ShouldBindQuery(&arg); err != nil {
		s.resFail(ctx, err, 400)
		return
	}

	accounts, err := s.db.ListAccounts(ctx, db.ListAccountsParams{
		Limit: arg.PageSize, Offset: (arg.Page - 1) * arg.PageSize,
	})
	if err != nil {
		s.resFail(ctx, err, 500)
		return
	}
	s.resOk(ctx, accounts)
}
