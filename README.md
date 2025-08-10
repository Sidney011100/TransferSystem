// goland version
# Internal Transfer System

# Prerequisties
1. Docker Desktop
2. Go (version `go1.23.5 darwin/amd64`)

## Big Idea
1. Run postgres using Docker, create the database and insert the tables with `go-lang/migrate`
2. (Optional) Configuring the port to call the app on local
3. Running the app!

### Creating the database
1. make reset-db
2. make start-db
2. make create-db
3. make migrate

### To run the application:



## Others
### Initial Manual Set Up 
#### Installing Postgres
1. download docker desktop
2. run postgres container
3. `psql --version`
4. If it does not exist, `brew install postgresql`
5. On docker, find the postgres image and run it

#### Database set up
1. `psql -h localhost -U postgres`
2.  To view the current databases `\l`
3. Create the database for this app `CREATE DATABASE transfer_db`


#### Inserting the tables
We are using `golang-migrate/migrate` tool
1. Install with  `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
2. The tables are written in `migrations/001_init.up.sql`. Run
```bash
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/transfer_db?sslmode=disable" up
```



