package usecase

import "mohamadelabror.me/simplebankappv2/repository"

type AddLogUseCase interface {
	AddLog(id, senderAccountNumber, receiverAccountNumber, time string, isMerchant bool) error
}

type addLogUseCase struct {
	customerRepo repository.CustomerRepo
}

func (a *addLogUseCase) AddLog(id, senderAccountNumber, receiverAccountNumber, time string, isMerchant bool) error {
	return a.customerRepo.AddLogToHistory(senderAccountNumber, receiverAccountNumber, time, id, isMerchant)
}

func NewAddLogUseCase(custRepo repository.CustomerRepo) AddLogUseCase {
	return &addLogUseCase{
		customerRepo: custRepo,
	}
}
