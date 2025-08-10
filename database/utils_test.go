package database

import (
	"context"
	"log"
	"os"
	"os/exec"
	"testing"
	"transferSystem/model"

	"github.com/joho/godotenv"
)

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
	InitDatabase(dsn)
	code := m.Run()

	CloseDatabase()
	cmd = exec.Command("make", "drop-test-db")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to run 'make drop-test-db': %v", err)
	}
	os.Exit(code)
}

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()
	inputBalanceToExpectedResult := []struct {
		acc     *model.NewAccount
		wantErr bool
	}{
		{
			acc:     &model.NewAccount{AccountId: 1, InitialBalance: "100.00"},
			wantErr: false,
		},
		{
			acc:     &model.NewAccount{AccountId: 2, InitialBalance: "132"},
			wantErr: true,
		},
		{
			acc:     &model.NewAccount{AccountId: 3, InitialBalance: "-100.89"},
			wantErr: true,
		},
	}

	for _, testcase := range inputBalanceToExpectedResult {
		err := CreateAccount(ctx, testcase.acc)
		if (err != nil) != testcase.wantErr {
			t.Errorf("CreateAccount(%v) error = %v, wantErr %v", testcase.acc, err, testcase.wantErr)
		}
	}
}
