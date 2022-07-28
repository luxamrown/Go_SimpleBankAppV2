package usecase

import (
	"mohamadelabror.me/simplebankappv2/model"
	"mohamadelabror.me/simplebankappv2/repository"
)

type GetTransactionDetailUseCase interface {
	GetTransactionDetail(idDetail, idHistory string, isMerchant bool) (model.TransactionDetail, error)
}

type getTransactionDetailUseCase struct {
	customerRepo repository.CustomerRepo
}

func (g *getTransactionDetailUseCase) GetTransactionDetail(idDetail, idHistory string, isMerchant bool) (model.TransactionDetail, error) {
	return g.customerRepo.GetTransactionDetail(idDetail, idHistory, isMerchant)
}

func NewGetTransactionDetailUseCase(custRepo repository.CustomerRepo) GetTransactionDetailUseCase {
	return &getTransactionDetailUseCase{
		customerRepo: custRepo,
	}
}
