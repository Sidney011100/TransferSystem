package account

import (
	"log"
	"testing"
	db "transferSystem/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db.InitDatabase("postgres://user:password@localhost:5432/test_transfer_db?sslmode=disable")
	defer db.CloseDatabase()
}

func TestCreateAccount(t *testing.T) {
	db.InitDatabase("postgres://postgres:password@localhost:5432/test_transfer_db?sslmode=disable")
	err := CreateAccount(123652, "100.00")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	//err = CreateAccount("xx", "100.00")
	//if err == nil {
	//	t.Fatalf("expected error, got no error")
	//}

	err = CreateAccount(246, "1.9.2")
	if err == nil {
		t.Fatalf("expected error, got no error")
	}

	err = CreateAccount(369, "-100.89")
	if err == nil {
		t.Fatalf("expected error, got no error")
	}
}
