package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/model"
	"mohamadelabror.me/simplebankappv2/repository"
)

type GetAllTransactionUseCase interface {
	GetAllTransactionDetail(idCustomer string) ([]model.TransactionDetail, *response.ErrorResp)
}

type getAllTransactionUseCase struct {
	customerRepo repository.CustomerRepo
}

func (g *getAllTransactionUseCase) GetAllTransactionDetail(idCustomer string) ([]model.TransactionDetail, *response.ErrorResp) {
	return g.customerRepo.GetAllTransactionDetail(idCustomer)
}

func NewGetAllTransactionUseCase(custRepo repository.CustomerRepo) GetAllTransactionUseCase {
	return &getAllTransactionUseCase{
		customerRepo: custRepo,
	}
}
