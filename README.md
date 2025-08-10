# Internal Transfer System
This is an internal transfers application, it aims to facilitate transactions of the same currency between accounts. 
It was built with Golang and runs on postgres. 

# Prerequisites
1. Install and run Docker Desktop
2. Go (version `go1.23.5 or newer`)
3. [`golang-migrate/tool`](https://github.com/golang-migrate/migrate)

# Setting up
### 1. Clone Repo
```bash
git clone https://github.com/Sidney011100/TransferSystem.git
```
2. cd into repo `cd TransferSystem`
3. Get dependencies 
```bash
go mod tidy
```


### 2. Create Database
1. `make reset-db`
Expected Outcome 
```bash
    docker stop postgres
    postgres
    postgres
    Waiting for Postgres to start...
    Postgres is ready!
    PGPASSWORD=postgres dropdb -h localhost -p 5432 -U postgres transfer_system_db || true
    Creating database...
    PGPASSWORD=postgres createdb -h localhost -p 5432 -U postgres transfer_system_db || true
    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/transfer_system_db?sslmode=disable" up
    1/u init (15.012048ms)
    2/u create_account (26.047674ms)
```

### 3. Run Application
run application 
```bash
go run main.go
```


## Features
### 1. Create Accounts `[POST] api/v1/accounts`
Description: Creates an account based on the account id and initial balance.

Assumptions: 
1. account_id is more than 0, and cannot be more than 9223372036854775807.
2. An account cannot be opened with a negative initial balance.
3. A user cannot open an account with an account_id previously used.
4. If it is successful, returning null is acceptable. Otherwise, an error will be encapsulated and sent.

Example of request
```json
{
    "account_id": 456,
    "initial_balance": "100.12345"
}
```

Example of error response
```json
{
  "err": "account ID 2456 already taken, please choose another"
}
```

### 2. Get Accounts `[GET] api/v1/accounts/{account_id}`
Description: Returns an account previously set up, with an account_id. 

Example of response
```json
{
    "account_id": 124,
    "balance": "100.23344"
}
```


### 3. Make Transactions`[POST]  api/v1/transactions`
Description: Transfers an `amount` from a `source_account_id` to a `destination_account_id`.

Assumptions
1. Both source account and destination account has to exist.
2. An `amount` has to be more than 0 for a transaction to take place. 
3. A transaction cannot be made if it leaves the source account with a negative balance. 
4. A transaction is made from the source account, by the source user. A balance of the source account will be displayed for the user

Example of request
```json
{
    "source_account_id": 456,
    "destination_account_id": 123,
    "amount": "100.12345"
}
```

Example of successful result, source account balance
```json
{
  "account_id": 2456,
  "balance": "99.32345"
}
```

Example of unsuccessful result
```json
{
  "err": "transaction failed account 2456 has insufficient funds, current balance 99.32345"
}
```


## Testing
### Running test
Test cases were written in `/test/internal_test.go`

To run the testcase `-go test ./test/*`

Expected Result
```bash
Creating database...
PGPASSWORD=postgres createdb -h localhost -p 5432 -U postgres test_transfer_db || true
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/test_transfer_db?sslmode=disable" up
1/u init (17.482237ms)
2/u create_account (30.10501ms)
=== RUN   TestGetAccount
--- PASS: TestGetAccount (0.01s)
=== RUN   TestCreateAccount
--- PASS: TestCreateAccount (0.00s)
=== RUN   TestUpdateAccount
--- PASS: TestUpdateAccount (0.01s)
=== RUN   TestTransaction
--- PASS: TestTransaction (0.04s)
PASS
PGPASSWORD=postgres dropdb -h localhost -p 5432 -U postgres test_transfer_db || true
```



## Others
### Alternative Set Up Process
#### Installing Postgres
1. download docker desktop
2. run postgres container
3. `psql --version`
4. If it does not exist, `brew install postgresql`
5. On docker, find the postgres image and run it

#### Database set up
1. `psql -h localhost -U postgres`
2.  To view the current databases `\l`
3. Create the database for this app `CREATE DATABASE transfer_system_db`


#### Inserting the tables
We are using `golang-migrate/migrate` tool
1. Install with  `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
2. The tables are written in `migrations` folder. Run the command below to create the tables and populate some accounts.
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/transfer_system_db?sslmode=disable" up
```



