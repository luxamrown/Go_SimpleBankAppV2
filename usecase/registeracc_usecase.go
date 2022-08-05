package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/model"
	"mohamadelabror.me/simplebankappv2/repository"
)

type RegisterAccountUseCase interface {
	RegisterAccount(id, accountNumber, userName, userPin, userPassword string, balance int) *response.ErrorResp
}

type registerAccountUseCase struct {
	customerRepo repository.CustomerRepo
}

func (r *registerAccountUseCase) RegisterAccount(id, accountNumber, userName, userPin, userPassword string, balance int) *response.ErrorResp {
	newAccount := model.NewCustomer(id, accountNumber, userName, userPin, userPassword, balance)
	return r.customerRepo.RegisterNewAccount(newAccount)

}

func NewRegisterAccountUseCase(custRepo repository.CustomerRepo) RegisterAccountUseCase {
	return &registerAccountUseCase{
		customerRepo: custRepo,
	}
}
