# Golang REST and gRPC API - work in progress
Golang REST API and 2 (micro)services - work in progress

## Installation
1. After cloning this repository, cd into it and run `docker-compose up -d`.
This will bring up the postgress db in a docker installation. Remember, if you destroy the docker container the data will be lost. We don't use volumes here.
2. Copy the contents of `etc/createTablesWithData.sql` into your db console and run it in order to create the tables and insert seed data. Caution, it will install 1M records.
3. generate code from proto with make proto
4. run the server `go run ./salesservice/cmd/sales/main.go`
5. run the client `go run ./productservice/cmd/products/main.go`
