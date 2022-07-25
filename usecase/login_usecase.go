package usecase

import (
	"mohamadelabror.me/simplebankappv2/repository"
)

type LoginUseCase interface {
	Login(username, password, token string) (string, error)
}

type loginUseCase struct {
	customerRepo repository.CustomerRepo
}

func (l *loginUseCase) Login(username, password, token string) (string, error) {
	idLoginAccount, err := l.customerRepo.Login(username, password)
	if err != nil {
		return "", err
	}
	err = l.customerRepo.SaveToken(idLoginAccount, token)
	if err != nil {
		return "", err
	}
	return idLoginAccount, nil
}

func NewLoginUseCase(custRepo repository.CustomerRepo) LoginUseCase {
	return &loginUseCase{
		customerRepo: custRepo,
	}
}
