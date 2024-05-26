package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount , valid := server.validAccount(ctx,req.ToAccountID,req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayload_key).(*token.Payload)
	if fromAccount.Owner != authPayload.Username { 
		err := errors.New("from account doens't belong to the authenticated user")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ , valid = server.validAccount(ctx,req.ToAccountID,req.Currency)
	if !valid {
		return
	}
	// if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
	// 	return
	// }

	// if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
	// 	return
	// }

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// Pass the db.Store instance and context along with the arguments
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account , bool) {
	account, err := server.store.GetAccount(ctx, int32(accountId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account,false
	}

	if account.Currency != currency {
		err := fmt.Errorf("Account [%d] currency mismatch: %s vs %s ", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account,true
}
