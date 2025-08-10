package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	db "transferSystem/database"
	"transferSystem/model"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

func GetAccount(id int64) (*model.Account, error) {
	ctx := context.Background()
	res, err := db.GetAccount(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return res, fmt.Errorf(ErrAccountNotFound, id)
	}
	return res, err
}

func CreateAccount(req *model.NewAccount) error {
	initBalance := req.InitialBalance
	if !isStringValidNumber(req.InitialBalance) {
		return fmt.Errorf(ErrInvalidAmount, initBalance)
	}

	id := req.AccountId
	if id <= 0 {
		return errors.New("id must be a positive integer")
	}

	balance, err := decimal.NewFromString(initBalance)
	if err != nil {
		return fmt.Errorf(ConversionFailed, err)
	}

	if balance.IsNegative() {
		return fmt.Errorf(ErrInitialBalanceNotPositive)
	}

	newAccount := model.NewAccount{
		AccountId:      id,
		InitialBalance: initBalance,
	}

	ctx := context.Background()
	err = db.CreateAccount(ctx, &newAccount)
	return err
}

func hasSufficientFunds(account *model.Account, fund decimal.Decimal) bool {
	balance, err := decimal.NewFromString(account.Balance)
	if err != nil {
		log.Fatalf("get account balance err: %v", err)
		return false
	}
	if balance.Sub(fund).IsNegative() {
		return false
	}
	return true
}

func UpdateAccount(ctx context.Context, account *model.Account, fund decimal.Decimal) error {
	balance, err := decimal.NewFromString(account.Balance)
	if err != nil {
		return fmt.Errorf(ConversionFailed, err)
	}
	newBalance := balance.Add(fund)
	if newBalance.IsNegative() {
		return fmt.Errorf(ErrAccountHasInsufficientFunds, account.AccountId)
	}
	account.Balance = newBalance.String()
	return db.UpdateAccount(ctx, account)
}
