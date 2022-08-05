package usecase

import (
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/repository"
)

type LoginUseCase interface {
	Login(username, password, token string) (string, *response.ErrorResp)
}

type loginUseCase struct {
	customerRepo repository.CustomerRepo
}

func (l *loginUseCase) Login(username, password, token string) (string, *response.ErrorResp) {
	idLoginAccount, err := l.customerRepo.Login(username, password, token)
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
