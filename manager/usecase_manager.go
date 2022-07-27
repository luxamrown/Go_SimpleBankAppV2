package manager

import (
	"mohamadelabror.me/simplebankappv2/usecase"
)

type UseCaseManager interface {
	RegisterAccountUseCase() usecase.RegisterAccountUseCase
	LoginUseCase() usecase.LoginUseCase
	LogoutUseCase() usecase.LogoutUseCase
	TransferUseCase() usecase.TransferUseCase
	AddLogUseCase() usecase.AddLogUseCase
}

type useCaseManager struct {
	repo RepoManager
}

func (u *useCaseManager) RegisterAccountUseCase() usecase.RegisterAccountUseCase {
	return usecase.NewRegisterAccountUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) LoginUseCase() usecase.LoginUseCase {
	return usecase.NewLoginUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) LogoutUseCase() usecase.LogoutUseCase {
	return usecase.NewLogoutUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) TransferUseCase() usecase.TransferUseCase {
	return usecase.NewTransferUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) AddLogUseCase() usecase.AddLogUseCase {
	return usecase.NewAddLogUseCase(u.repo.CustomerRepo())
}

func NewUseCaseManager(repo RepoManager) UseCaseManager {
	return &useCaseManager{
		repo: repo,
	}
}
