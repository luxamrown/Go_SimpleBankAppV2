package usecase

import "mohamadelabror.me/simplebankappv2/repository"

type AddTransactionDetailUseCase interface {
	AddTransactionDetail(id, historyId, customerId, message string, amount int) error
}

type addTransactionDetailUseCase struct {
	customerRepo repository.CustomerRepo
}

func (a *addTransactionDetailUseCase) AddTransactionDetail(id, historyId, customerId, message string, amount int) error {
	return a.customerRepo.AddTransactionDetail(id, historyId, customerId, message, amount)
}

func NewAddTransactionDetailUseCase(custRepo repository.CustomerRepo) AddTransactionDetailUseCase {
	return &addTransactionDetailUseCase{
		customerRepo: custRepo,
	}
}
