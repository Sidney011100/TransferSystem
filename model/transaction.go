package model

type NewTransaction struct {
	SourceAccountId      int64  `json:"source_account_id" binding:"required"`
	DestinationAccountId int64  `json:"destination_account_id" binding:"required"`
	Amount               string `json:"amount" binding:"required"`
}
