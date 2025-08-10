package handler

import (
	"github.com/gin-gonic/gin"
)

const (
	ErrInvalidAccount = "invalid account id %s"
	ErrInvalidJson    = "invalid input %s"
	ErrAccountTaken   = "account ID %d already taken, please choose another"

	ErrTransactionFailed = "transaction failed %s"
)

type Response struct {
	Msg string `json:"err"`
}

func doResp(c *gin.Context, data interface{}, err error) {
	httpRespCode := 200
	if err != nil {
		respObj := &Response{
			Msg: err.Error(),
		}
		c.JSON(httpRespCode, respObj)
		return
	}
	c.JSON(httpRespCode, data)
}
