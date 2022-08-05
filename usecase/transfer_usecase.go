package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/repository"
)

type TransferUseCase interface {
	Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber string, amount int, isMerchant bool) *response.ErrorResp
}

type transferUseCase struct {
	customerRepo repository.CustomerRepo
}

func (t *transferUseCase) Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber string, amount int, isMerchant bool) *response.ErrorResp {
	return t.customerRepo.Transfer(sender_id, sender_pin, senderAccountNumber, receiverAccountNumber, amount, isMerchant)
}

func NewTransferUseCase(custRepo repository.CustomerRepo) TransferUseCase {
	return &transferUseCase{
		customerRepo: custRepo,
	}
}
