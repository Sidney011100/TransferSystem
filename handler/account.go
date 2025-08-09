package handler

import (
	"fmt"
	"strconv"
	"strings"
	"transferSystem/account"
	"transferSystem/database"
	"transferSystem/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func UserCreateAccount(c *gin.Context) {
	req := &model.NewAccount{}
	if err := c.BindJSON(req); err != nil {
		doResp(c, nil, fmt.Errorf(ErrInvalidJson, err))
		return
	}

	initBalance := req.InitialBalance
	if !isStringValidNumber(req.InitialBalance) {
		doResp(c, nil, fmt.Errorf(ErrInvalidAmount, initBalance))
	}

	err := account.CreateAccount(req.AccountId, req.InitialBalance)
	if err != nil && strings.Contains(err.Error(), database.ErrDupKey) {
		doResp(c, nil, fmt.Errorf(ErrAccountTaken, req.AccountId))
		return
	}

	doResp(c, nil, err)
	return
}

func UserGetAccount(c *gin.Context) {
	accId := c.Param("account_id")
	i, err := strconv.ParseInt(accId, 10, 64)
	if err != nil {
		doResp(c, nil, fmt.Errorf(ErrInvalidAccount, err))
		return
	}
	res, err := account.GetAccount(i)
	doResp(c, res, err)
}

func UserTransaction(c *gin.Context) {
	req := &model.NewTransaction{}
	if err := c.BindJSON(req); err != nil {
		doResp(c, nil, fmt.Errorf(ErrInvalidJson, err))
		return
	}

	inputAmount := req.Amount
	if !isStringValidNumber(inputAmount) {
		doResp(c, nil, fmt.Errorf(ErrInvalidAmount, inputAmount))
		return
	}

	transferAmount, err := decimal.NewFromString(inputAmount)
	if err != nil {
		doResp(c, nil, fmt.Errorf(ErrInvalidTransferAmount, inputAmount, err))
		return
	}

	err = account.ProcessTransaction(req.SourceAccountId, req.DestinationAccountId, transferAmount)
	if err != nil {
		doResp(c, nil, fmt.Errorf(ErrTransactionFailed, err))
		return
	}

	doResp(c, nil, nil)
	return
}
