package model

type Account struct {
	AccountId int64  `json:"account_id" binding:"required"`
	Balance   string `json:"balance"`
}

type NewAccount struct {
	AccountId      int64  `json:"account_id" binding:"required"`
	InitialBalance string `json:"initial_balance" binding:"required"`
}
