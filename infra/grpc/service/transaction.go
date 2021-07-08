package service

import (
	"context"

	"github.com/thienry/code-bank/dto"
	"github.com/thienry/code-bank/infra/grpc/pb"
	"github.com/thienry/code-bank/usecase"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) Payment(ctx context.Context, in *pb.PaymentRequest) (*empty.Empty, error) {
	transactionDto := dto.Transaction{
		Name: in.GetCreditCard().GetName(),
		Number: in.GetCreditCard().GetNumber(),
		ExpirationMonth: in.GetCreditCard().GetExpirationMonth(),
		ExpirationYear: in.GetCreditCard().GetExpirationYear(),
		Amount: in.GetAmount(),
		CVV: in.GetCreditCard().GetCVV(),
		Store: in.GetStore(),
		Description: in.GetDescription(),
	}

	transaction, err := t.ProcessTransactionUseCase.ProcessTransaction(transactionDto)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}

	if transaction.Status != "approved" {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, "Transaction rejected by the bank!")
	}

	return &empty.Empty{}, nil
}
