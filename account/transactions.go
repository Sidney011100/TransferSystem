package account

import (
	"context"
	"fmt"
	db "transferSystem/database"
	"transferSystem/model"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ProcessTransaction(sourceAccountId, destinationAccountId int64, amount decimal.Decimal) error {
	ctx := context.Background()
	sourceAcc, err := GetAccount(sourceAccountId)
	if err != nil {
		return fmt.Errorf("source " + err.Error())
	}

	destinationAcc, err := GetAccount(destinationAccountId)
	if err != nil {
		return fmt.Errorf("destination " + err.Error())
	}

	//create transaction in db
	transId, err := db.CreateTransaction(ctx, sourceAcc, destinationAcc, amount)
	if err != nil {
		return fmt.Errorf(ErrFailedToCreateTransaction, err)
	}

	// lock to ensure no one takes the funds if its sufficient.
	if !hasSufficientFunds(sourceAcc, amount) {
		return fmt.Errorf(ErrAccountHasInsufficientFunds, sourceAccountId)
	}

	err = updateTransactionSrcAcc(ctx, transId, sourceAcc, amount)

	// unlock

	err = updateTransactionDestAcc(ctx, transId, destinationAcc, amount)

	return err
}

func updateTransactionSrcAcc(ctx context.Context, id uuid.UUID, account *model.Account, amount decimal.Decimal) error {
	err := updateAccount(ctx, account, amount.Neg())
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
	err := updateAccount(ctx, account, amount)
	if err != nil {
		return err
	}

	err = db.UpdateTransactionDest(ctx, id, account)
	if err != nil {
		return err
	}
	return nil
}
