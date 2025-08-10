package database

import (
	"context"
	"log"
	"time"
	"transferSystem/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

const (
	//t_account = "t_account"
	//t_transaction = "t_transaction"
	ErrDupKey = "ERROR: duplicate key value violates unique constraint"
)

var conn *pgx.Conn

func InitDatabase(dsn string) {
	var err error
	conn, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDatabase() {
	if conn != nil {
		conn.Close(context.Background())
	}
}

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

func CreateTransaction(ctx context.Context, srcAcc, destAcc *model.Account, fund decimal.Decimal) (uuid.UUID, error) {
	id := uuid.New()
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	sql := `INSERT INTO t_transaction(uuid, source_account_id, source_init, source_after, destination_account_id, destination_init, destination_after, amount, create_time, modify_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := conn.Exec(ctx, sql, id, srcAcc.AccountId, srcAcc.Balance, srcAcc.Balance, destAcc.AccountId, destAcc.Balance, destAcc.Balance, fund.String(), timeNow, timeNow)
	return id, err
}

func UpdateTransactionSource(ctx context.Context, transId uuid.UUID, srcAcc *model.Account) error {
	sql := `UPDATE t_transaction 
				SET source_after = $1,  
					modify_time = $2
				WHERE uuid = $3`
	_, err := conn.Exec(ctx, sql, srcAcc.Balance, time.Now().Format("2006-01-02 15:04:05"), transId)
	return err
}

func UpdateTransactionDest(ctx context.Context, transId uuid.UUID, destAcc *model.Account) error {
	sql := `UPDATE t_transaction 
			SET destination_after = $1,
			    success = TRUE,
				modify_time = $2
            WHERE uuid = $3`
	_, err := conn.Exec(ctx, sql, destAcc.Balance, time.Now().Format("2006-01-02 15:04:05"), transId)
	return err
}
