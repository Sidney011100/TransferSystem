package account

import (
	"log"
	"math/rand"
	"os"
	"testing"
	db "transferSystem/database"

	"github.com/joho/godotenv"
)

var randomId []int64

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found")
	}

	db.InitDatabase(os.Getenv("DATABASE_URL_TEST"))
	defer db.CloseDatabase()

	randomId = make([]int64, 3)
	for i := 0; i < 3; i++ {
		randInt := rand.Int63n(1000)
		randomId[i] = randInt
	}
	m.Run()

}

func TestCreateAccount(t *testing.T) {
	inputBalanceToExpectedResult := map[string]bool{
		"1.9.2":   false,
		"-100.89": false,
		"100.00":  true,
	}

	for input, expectedResult := range inputBalanceToExpectedResult {
		err := CreateAccount(randomId[0], input)
		if (err != nil) == expectedResult {
			if expectedResult {
				t.Fatalf("expected no error, got %v", err)
			} else {
				t.Fatalf("expected error, got no error for balance %s", input)
			}
		}
	}
}

func TestGetAccount(t *testing.T) {
	inputBalanceToExpectedResult := map[int64]bool{
		randomId[0]: true,
		randomId[1]: false,
	}
	for input, expectedResult := range inputBalanceToExpectedResult {
		_, err := GetAccount(input)
		if (err != nil) == expectedResult {
			if expectedResult {
				t.Fatalf("expected no error, got %v", err)
			} else {
				t.Fatalf("expected error, got no error for balance %d", input)
			}
		}
	}
}
