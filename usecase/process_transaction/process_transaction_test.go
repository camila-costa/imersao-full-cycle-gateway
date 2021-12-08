package process_transaction

import (
	"testing"
	"time"
	"github.com/camila-costa/imersao-gateway/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	mock_repository "github.com/camila-costa/imersao-gateway/domain/repository/mock"
	mock_broker "github.com/camila-costa/imersao-gateway/adapter/broker/mock"
)

func TestProcessTransaction_ExecuteInvalidCreditCard(t *testing.T) {
	input := TransactionDtoInput{
		ID: "1",
		AccountID: "1",
		CreditCardNumber: "400000000000000",
		CreditCardName: "Teste da Silva",
		CreditCardExpirationMonth: 12,
		CreditCardExpirationYear: time.Now().Year(),
		CreditCardExpirationCVV: 123,
		Amount: 200,
	}
	expectedOutput := TransactionDtoOutput{
		ID: "1",
		Status: entity.REJECTED,
		ErrorMessage: "invalid credit card number",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock_repository.NewMockTransactionRepository(ctrl)
	repositoryMock.EXPECT().
		Insert(input.ID, input.AccountID, input.Amount, expectedOutput.Status, expectedOutput.ErrorMessage).
		Return(nil)

	producerMock := mock_broker.NewMockProducerInterface(ctrl)
	producerMock.EXPECT().Publish(expectedOutput, []byte(input.ID), "transactions_result")

	usecase := NewProcessTransaction(repositoryMock, producerMock, "transactions_result")
	output, err := usecase.Execute(input)
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestProcessTransaction_ExecuteRejectedTransaction(t *testing.T) {
	input := TransactionDtoInput{
		ID: "1",
		AccountID: "1",
		CreditCardNumber: "4193523830170205",
		CreditCardName: "Teste da Silva",
		CreditCardExpirationMonth: 12,
		CreditCardExpirationYear: time.Now().Year(),
		CreditCardExpirationCVV: 123,
		Amount: 1200,
	}
	expectedOutput := TransactionDtoOutput{
		ID: "1",
		Status: entity.REJECTED,
		ErrorMessage: "you dont have limit for this transaction",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock_repository.NewMockTransactionRepository(ctrl)
	repositoryMock.EXPECT().
		Insert(input.ID, input.AccountID, input.Amount, expectedOutput.Status, expectedOutput.ErrorMessage).
		Return(nil)

	producerMock := mock_broker.NewMockProducerInterface(ctrl)
	producerMock.EXPECT().Publish(expectedOutput, []byte(input.ID), "transactions_result")

	usecase := NewProcessTransaction(repositoryMock, producerMock, "transactions_result")
	output, err := usecase.Execute(input)
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestProcessTransaction_ExecuteApprovedTransaction(t *testing.T) {
	input := TransactionDtoInput{
		ID: "1",
		AccountID: "1",
		CreditCardNumber: "4193523830170205",
		CreditCardName: "Teste da Silva",
		CreditCardExpirationMonth: 12,
		CreditCardExpirationYear: time.Now().Year(),
		CreditCardExpirationCVV: 123,
		Amount: 900,
	}
	expectedOutput := TransactionDtoOutput{
		ID: "1",
		Status: entity.APPROVED,
		ErrorMessage: "",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mock_repository.NewMockTransactionRepository(ctrl)
	repositoryMock.EXPECT().
		Insert(input.ID, input.AccountID, input.Amount, expectedOutput.Status, expectedOutput.ErrorMessage).
		Return(nil)

	producerMock := mock_broker.NewMockProducerInterface(ctrl)
	producerMock.EXPECT().Publish(expectedOutput, []byte(input.ID), "transactions_result")

	usecase := NewProcessTransaction(repositoryMock, producerMock, "transactions_result")
	output, err := usecase.Execute(input)
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, output)
}