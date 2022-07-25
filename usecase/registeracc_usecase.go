package usecase

import (
	"mohamadelabror.me/simplebankappv2/model"
	"mohamadelabror.me/simplebankappv2/repository"
)

type RegisterAccountUseCase interface {
	RegisterAccount(id, accountNumber, userName, userPassword string, balance int) error
}

type registerAccountUseCase struct {
	customerRepo repository.CustomerRepo
}

func (r *registerAccountUseCase) RegisterAccount(id, accountNumber, userName, userPassword string, balance int) error {
	newAccount := model.NewCustomer(id, accountNumber, userName, userPassword, balance)
	err := r.customerRepo.RegisterNewAccount(newAccount)
	if err != nil {
		return err
	}
	return err
}

func NewRegisterAccountUseCase(custRepo repository.CustomerRepo) RegisterAccountUseCase {
	return &registerAccountUseCase{
		customerRepo: custRepo,
	}
}
