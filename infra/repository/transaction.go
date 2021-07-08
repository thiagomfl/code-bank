package repository

import (
	"database/sql"
	"errors"

	"github.com/thienry/code-bank/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db }
}

func (t *TransactionRepository) SaveTransaction(transaction domain.Transaction, creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`
		insert into transactions(id, credit_card, amount, status, description, store, created_at)
		values($1, $2, $3, $4, $5, $6, $7)
	`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)

	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		err = t.updateBalance(creditCard)

		if err != nil {
			return err
		}
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepository) updateBalance(creditCard domain.CreditCard) error {
	_, err := t.db.Exec("update credit_cards set balance = $1 where id = $2", creditCard.Balance, creditCard.ID)
	
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepository) CreateCreditCard(creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`
		insert into credit_cards(id, name, number, expiration_month, expiration_year, cvv, balance, limit)
		values($1, $2, $3, $4, $5, $6, $7, $8)
	`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)

	if err != nil {
		return err
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepository) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare("select id, balance, balance_limit from credit_cards where number=$1")

	if err != nil {
		return c, err
	}

	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card doesn`t exists")
	}

	return c, nil
}
 