package database

import (
	"context"
	"transferSystem/model"
)

func CreateAccount(ctx context.Context, acc *model.NewAccount) error {
	sql := `INSERT INTO t_account (account_id, balance) VALUES ($1, $2)`
	_, err := conn.Exec(ctx, sql, acc.AccountId, acc.InitialBalance)
	return err
}

func GetAccount(ctx context.Context, accid int64) (*model.Account, error) {
	account := &model.Account{}
	sql := `SELECT * FROM t_account WHERE account_id = $1`
	row := conn.QueryRow(ctx, sql, accid)
	err := row.Scan(&account.AccountId, &account.Balance)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func UpdateAccount(ctx context.Context, acc *model.Account) error {
	sql := `UPDATE t_account SET balance = $1 WHERE account_id = $2`
	_, err := conn.Exec(ctx, sql, acc.Balance, acc.AccountId)
	return err
}
