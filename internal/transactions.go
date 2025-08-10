package internal

import (
	"context"
	"fmt"
	db "transferSystem/database"
	"transferSystem/model"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ProcessTransaction(ctx context.Context, req *model.NewTransaction) (*model.Account, error) {
	inputAmount := req.Amount
	if !isStringValidNumber(inputAmount) {
		return nil, fmt.Errorf(ErrInvalidAmount, inputAmount)
	}

	amount, err := decimal.NewFromString(inputAmount)
	if err != nil {
		return nil, fmt.Errorf(ErrInvalidTransferAmount, inputAmount, err)
	}

	sourceAccountId := req.SourceAccountId
	destinationAccountId := req.DestinationAccountId

	sourceAcc, err := GetAccount(sourceAccountId)
	if err != nil {
		return nil, fmt.Errorf("source " + err.Error())
	}

	destinationAcc, err := GetAccount(destinationAccountId)
	if err != nil {
		return nil, fmt.Errorf("destination " + err.Error())
	}

	//create transaction in db
	transId, err := db.CreateTransaction(ctx, sourceAcc, destinationAcc, amount)
	if err != nil {
		return sourceAcc, fmt.Errorf(ErrFailedToCreateTransaction, err)
	}

	// lock to ensure no one takes the funds if its sufficient.
	currentSourceBalance, isSuffice := hasSufficientFunds(sourceAcc, amount)
	if !isSuffice {
		return sourceAcc, fmt.Errorf(ErrAccountHasInsufficientFunds, sourceAccountId, currentSourceBalance)
	}

	err = updateTransactionSrcAcc(ctx, transId, sourceAcc, amount)

	// unlock

	err = updateTransactionDestAcc(ctx, transId, destinationAcc, amount)

	sourceAcc, err = GetAccount(sourceAccountId)
	if err != nil {
		return nil, fmt.Errorf("updated source " + err.Error())
	}

	return sourceAcc, err
}

func updateTransactionSrcAcc(ctx context.Context, id uuid.UUID, account *model.Account, amount decimal.Decimal) error {
	err := UpdateAccount(ctx, account, amount.Neg())
	if err != nil {
		return err
	}

	err = db.UpdateTransactionSource(ctx, id, account)
	if err != nil {
		return err
	}
	return nil
}

func updateTransactionDestAcc(ctx context.Context, id uuid.UUID, account *model.Account, amount decimal.Decimal) error {
	err := UpdateAccount(ctx, account, amount)
	if err != nil {
		return err
	}

	err = db.UpdateTransactionDest(ctx, id, account)
	if err != nil {
		return err
	}
	return nil
}
