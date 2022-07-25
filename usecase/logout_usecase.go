package usecase

import "mohamadelabror.me/simplebankappv2/repository"

type LogoutUseCase interface {
	Logout(id string) error
}

type logoutUseCase struct {
	customerRepo repository.CustomerRepo
}

func (l *logoutUseCase) Logout(id string) error {
	return l.customerRepo.Logout(id)
}

func NewLogoutUseCase(custRepo repository.CustomerRepo) LogoutUseCase {
	return &logoutUseCase{
		customerRepo: custRepo,
	}
}
