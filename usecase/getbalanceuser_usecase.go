package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/repository"
)

type GetBalanceUserUseCase interface {
	GetBalance(idCustomer, pin string) (*int, *response.ErrorResp)
}

type getBalanceUserUseCase struct {
	customerRepo repository.CustomerRepo
}

func (g *getBalanceUserUseCase) GetBalance(idCustomer, pin string) (*int, *response.ErrorResp) {
	return g.customerRepo.GetBalanceUser(idCustomer, pin)
}

func NewGetBalanceUserUseCase(custRepo repository.CustomerRepo) GetBalanceUserUseCase {
	return &getBalanceUserUseCase{
		customerRepo: custRepo,
	}
}
