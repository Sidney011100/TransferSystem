package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"testing"
	db "transferSystem/database"
	"transferSystem/model"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

var randomId []int64

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, relying on env vars for test")
	}

	cmd := exec.Command("make", "create-test-db")
	cmd.Dir = ".."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run 'make create-test-db': %v", err)
	}

	cmd = exec.Command("make", "migrate-test")
	cmd.Dir = ".."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run 'make migrate-test': %v", err)
	}

	dsn := os.Getenv("DATABASE_URL_TEST")
	if dsn == "" {
		log.Fatal("DATABASE_URL_TEST must be set in .env or environment")
	}
	db.InitDatabase(dsn)
	randomId = make([]int64, 3)
	for i := 0; i < 3; i++ {
		randInt := rand.Int63n(1000)
		randomId[i] = randInt
	}
	code := m.Run()

	db.CloseDatabase()
	cmd = exec.Command("make", "drop-test-db")
	cmd.Dir = ".."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to run 'make drop-test-db': %v", err)
	}
	os.Exit(code)
}

func TestGetAccount(t *testing.T) {
	inputBalanceToExpectedResult := []struct {
		acc         *model.Account
		wantErr     bool
		expectedErr error
	}{
		{
			acc:         &model.Account{AccountId: 1001, Balance: "100.123"},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			acc:         &model.Account{AccountId: 357, Balance: "100.123"},
			wantErr:     true,
			expectedErr: fmt.Errorf(ErrAccountNotFound, 357),
		},
		{
			acc:         &model.Account{AccountId: 1002, Balance: "1000000.12378857812"},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, testcase := range inputBalanceToExpectedResult {
		acc, err := GetAccount(testcase.acc.AccountId)
		if (err != nil) != testcase.wantErr {
			t.Fatal(fmt.Errorf("TestGetAccount(%v) error = %v", testcase.acc, dealWithErrUnexpected(testcase.wantErr, err)))
		}
		if testcase.wantErr {
			if errors.Is(err, testcase.expectedErr) {
				t.Fatal(fmt.Errorf("TestGetAccount(%v) error = %v", testcase.acc, dealWithErrUnexpected(testcase.wantErr, err)))
			}
			continue
		}
		if !reflect.DeepEqual(acc, testcase.acc) {
			t.Errorf("TestGetAccount(%d) result different. Expected %v got %v", testcase.acc.AccountId, testcase.acc, acc)
		}
	}
}

func TestCreateAccount(t *testing.T) {
	inputBalanceToExpectedResult := []struct {
		acc        *model.NewAccount
		updatedAcc *model.Account
		wantErr    bool
	}{
		{
			acc:        &model.NewAccount{AccountId: randomId[0], InitialBalance: "1.9.2"},
			updatedAcc: nil,
			wantErr:    true,
		},
		{
			acc:        &model.NewAccount{AccountId: randomId[1], InitialBalance: "-100.89"}, // insufficient balance
			updatedAcc: nil,
			wantErr:    true,
		},
		{
			acc:        &model.NewAccount{AccountId: 1001, InitialBalance: "100.89"}, //duplicate
			updatedAcc: nil,
			wantErr:    true,
		},
		{
			acc:        &model.NewAccount{AccountId: 2001, InitialBalance: "8.23ff"}, //weird
			updatedAcc: nil,
			wantErr:    true,
		},
		{
			acc:        &model.NewAccount{AccountId: 1001, InitialBalance: "8.23245236"},
			updatedAcc: &model.Account{AccountId: 1001, Balance: "8.23245236"},
			wantErr:    true,
		},
	}

	for _, testcase := range inputBalanceToExpectedResult {
		newAcc := testcase.acc
		err := CreateAccount(newAcc)
		if (err != nil) != testcase.wantErr {
			t.Fatal(fmt.Errorf("TestCreateAccount(%v) error = %v", newAcc, dealWithErrUnexpected(testcase.wantErr, err)))
		}
		if testcase.wantErr {
			continue
		}
		acc, err := GetAccount(testcase.acc.AccountId)
		if acc != testcase.updatedAcc {
			t.Fatal(fmt.Errorf("TestCreateAccount(%v) expected acc %v, got %v", testcase.acc, testcase.updatedAcc, acc))
		}
	}
}

func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()
	inputBalanceToExpectedResult := []struct {
		acc        *model.Account
		val        string
		updatedAcc *model.Account
		wantErr    bool
	}{
		{
			acc:        &model.Account{AccountId: 3001, Balance: "100.123"},
			val:        "3.689",
			updatedAcc: &model.Account{AccountId: 3001, Balance: "103.812"},
			wantErr:    false,
		},
		{
			acc:        &model.Account{AccountId: 3002, Balance: "100.456"},
			val:        "-4.281",
			updatedAcc: &model.Account{AccountId: 3002, Balance: "96.175"},
			wantErr:    false,
		},
		{
			acc:        &model.Account{AccountId: 3003, Balance: "200.36923"},
			val:        "-230.89",
			updatedAcc: &model.Account{AccountId: 3003, Balance: "200.36923"},
			wantErr:    true,
		},
	}

	for _, testcase := range inputBalanceToExpectedResult {
		fund, err := decimal.NewFromString(testcase.val)
		if err != nil {
			t.Fatal(err)
		}
		err = UpdateAccount(ctx, testcase.acc, fund)
		if (err != nil) != testcase.wantErr {
			t.Fatal(fmt.Errorf("TestUpdateAccount(%v) error = %v", testcase.acc, dealWithErrUnexpected(testcase.wantErr, err)))
		}
		acc, err := GetAccount(testcase.acc.AccountId)
		if err != nil {
			t.Fatal(err)
		}
		if acc.Balance != testcase.updatedAcc.Balance {
			t.Errorf("TestUpdateAccount(%v) error: account balance expected %s, got %s", testcase.acc, testcase.updatedAcc.Balance, acc.Balance)
		}
	}
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	inputBalanceToExpectedResult := []struct {
		transaction  *model.NewTransaction
		resultSrcAcc *model.Account
		resultDstAcc *model.Account
		expectedErr  error
		wantErr      bool
	}{
		{
			transaction:  &model.NewTransaction{SourceAccountId: 4001, DestinationAccountId: 4002, Amount: "200"},
			resultSrcAcc: &model.Account{AccountId: 4001, Balance: "100.123"},
			resultDstAcc: &model.Account{AccountId: 4002, Balance: "100.456"},
			expectedErr:  fmt.Errorf(ErrAccountHasInsufficientFunds, 234),
			wantErr:      true,
		},
		{
			transaction:  &model.NewTransaction{SourceAccountId: 4001, DestinationAccountId: 4002, Amount: "200.1.2"},
			resultSrcAcc: &model.Account{AccountId: 4001, Balance: "100.123"},
			resultDstAcc: &model.Account{AccountId: 4002, Balance: "100.456"},
			expectedErr:  fmt.Errorf(ErrInvalidAmount, "200.1.2"),
			wantErr:      true,
		},
		{
			transaction:  &model.NewTransaction{SourceAccountId: 4001, DestinationAccountId: 4002, Amount: "67"},
			resultSrcAcc: &model.Account{AccountId: 4001, Balance: "33.123"},
			resultDstAcc: &model.Account{AccountId: 4002, Balance: "167.456"},
			expectedErr:  nil,
			wantErr:      false,
		},
		{
			transaction:  &model.NewTransaction{SourceAccountId: 4003, DestinationAccountId: 4002, Amount: "200.369"},
			resultSrcAcc: &model.Account{AccountId: 4003, Balance: "0"},
			resultDstAcc: &model.Account{AccountId: 4002, Balance: "367.825"},
			expectedErr:  nil,
			wantErr:      false,
		},
	}

	for _, testcase := range inputBalanceToExpectedResult {
		err := ProcessTransaction(ctx, testcase.transaction)
		if (err != nil) != testcase.wantErr {
			t.Errorf("TestProcessTransaction(%v) error = %v", testcase.transaction, dealWithErrUnexpected(testcase.wantErr, err))
		}
		acc, err := GetAccount(testcase.resultSrcAcc.AccountId)
		if err != nil {
			t.Fatal(err)
		}
		if acc.Balance != testcase.resultSrcAcc.Balance {
			t.Errorf("source account balance not expected")
		}
		acc, err = GetAccount(testcase.resultDstAcc.AccountId)
		if err != nil {
			t.Fatal(err)
		}
		if acc.Balance != testcase.resultDstAcc.Balance {
			t.Errorf("destination account balance not expected")
		}
	}
}

func dealWithErrUnexpected(wantErr bool, err error) error {
	if wantErr {
		return fmt.Errorf("expected error did not get")
	}
	return fmt.Errorf("not expecting error but got %s", err.Error())
}
