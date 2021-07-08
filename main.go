package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pg"
	"github.com/thienry/code-bank/domain"
	"github.com/thienry/code-bank/infra/repository"
	"github.com/thienry/code-bank/usecase"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Name = "Thiago"
	cc.Number = "00215153"
	cc.ExpirationMonth = 10
	cc.ExpirationYear = 2021
	cc.CVV = 857
	cc.Limit = 2000
	cc.Balance = 0

	repo := repository.NewTransactionRepository(db)
	repo.CreateCreditCard(*cc)
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepository(db)
	usecase := usecase.NewUseCaseTransaction(transactionRepository)
	return usecase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s 					sslmode=disable",
		"db", 5432, "postgres", "root", "codebank",
	)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Error on connect database...")
	}

	return db
}
