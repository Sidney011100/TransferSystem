package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"transferSystem/database"
	"transferSystem/internal"
	"transferSystem/model"

	"github.com/gin-gonic/gin"
)

func UserCreateAccount(c *gin.Context) {
	req := &model.NewAccount{}
	if err := c.BindJSON(req); err != nil {
		doResp(c, nil, fmt.Errorf(ErrInvalidJson, err))
		return
	}

	err := internal.CreateAccount(req)
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
	res, err := internal.GetAccount(i)
	doResp(c, res, err)
}

func UserTransaction(c *gin.Context) {
	req := &model.NewTransaction{}
	if err := c.BindJSON(req); err != nil {
		doResp(c, nil, fmt.Errorf(ErrInvalidJson, err))
		return
	}

	sourceAcc, err := internal.ProcessTransaction(context.Background(), req)
	if err != nil {
		doResp(c, sourceAcc, fmt.Errorf(ErrTransactionFailed, err))
		return
	}

	doResp(c, sourceAcc, nil)
	return
}
