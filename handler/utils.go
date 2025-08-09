package handler

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

const (
	ErrInvalidAmount         = "invalid amount/balance %s"
	ErrInvalidAccount        = "invalid account id %s"
	ErrInvalidJson           = "invalid input %s"
	ErrInvalidTransferAmount = "invalid transfer amount %s, err %s"
	ErrAccountTaken          = "account ID %d already taken, please choose another"

	ErrTransactionFailed = "transaction failed %s"
)

type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func doResp(c *gin.Context, data interface{}, err error) {
	httpRespCode := 200
	errCode := 0
	errMsg := ""
	if err != nil {
		errCode = -1
		errMsg = err.Error()
	}
	respObj := &Response{
		Status: errCode,
		Msg:    errMsg,
		Data:   data,
	}
	c.JSON(httpRespCode, respObj)
}

func isStringValidNumber(s string) bool {
	re := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return re.MatchString(s)
}
