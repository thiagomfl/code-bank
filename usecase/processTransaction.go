package usecase

import (
	"encoding/json"
	"os"
	"time"

	"github.com/thienry/code-bank/domain"
	"github.com/thienry/code-bank/dto"
	"github.com/thienry/code-bank/infra/kafka"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
	KafkaProducer kafka.KafkaProducer
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	creditCard := u.hydrateCreditCard(transactionDto)

	ccBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	creditCard.ID = ccBalanceAndLimit.ID
	creditCard.Limit = ccBalanceAndLimit.Limit
	creditCard.Balance = ccBalanceAndLimit.Balance

	t := u.newtransaction(transactionDto, ccBalanceAndLimit)
	t.ProcessAndValidate(creditCard)

	err = u.TransactionRepository.SaveTransaction(*t, *creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	transactionDto.ID = t.ID
	transactionDto.CreatedAt = t.CreatedAt

	transactionJson, err := json.Marshal(transactionDto)
	if err != nil {
		return domain.Transaction{}, err
	}
	
	err = u.KafkaProducer.Publish(string(transactionJson), os.Getenv("KAFKA_TRANSACTIONS_TOPIC"))
	if err != nil {
		return domain.Transaction{}, err
	}

	return *t, nil
}

func (u UseCaseTransaction) hydrateCreditCard(transactionDto dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	
	creditCard.CVV = transactionDto.CVV
	creditCard.Name = transactionDto.Name
	creditCard.Number = transactionDto.Number
	creditCard.ExpirationMonth = transactionDto.ExpirationMonth
	creditCard.ExpirationYear = transactionDto.ExpirationYear

	return creditCard
}

func (u UseCaseTransaction) newtransaction(transactionDto dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()
	
	t.CreditCardId = cc.ID
	t.Amount = transactionDto.Amount
	t.Store = transactionDto.Store
	t.Description = transactionDto.Description
	t.CreatedAt = time.Now()

	return t
}
