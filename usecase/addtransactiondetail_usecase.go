package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/repository"
)

type AddTransactionDetailUseCase interface {
	AddTransactionDetail(id, historyId, customerId, message string, amount int) *response.ErrorResp
}

type addTransactionDetailUseCase struct {
	customerRepo repository.CustomerRepo
}

func (a *addTransactionDetailUseCase) AddTransactionDetail(id, historyId, customerId, message string, amount int) *response.ErrorResp {
	return a.customerRepo.AddTransactionDetail(id, historyId, customerId, message, amount)
}

func NewAddTransactionDetailUseCase(custRepo repository.CustomerRepo) AddTransactionDetailUseCase {
	return &addTransactionDetailUseCase{
		customerRepo: custRepo,
	}
}
