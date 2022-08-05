package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/repository"
)

type AddLogUseCase interface {
	AddLog(id, senderAccountNumber, receiverAccountNumber, time string, isMerchant bool) *response.ErrorResp
}

type addLogUseCase struct {
	customerRepo repository.CustomerRepo
}

func (a *addLogUseCase) AddLog(id, senderAccountNumber, receiverAccountNumber, time string, isMerchant bool) *response.ErrorResp {
	return a.customerRepo.AddLogToHistory(senderAccountNumber, receiverAccountNumber, time, id, isMerchant)
}

func NewAddLogUseCase(custRepo repository.CustomerRepo) AddLogUseCase {
	return &addLogUseCase{
		customerRepo: custRepo,
	}
}
