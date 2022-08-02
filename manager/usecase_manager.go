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
	AddTransactionDetailUseCase() usecase.AddTransactionDetailUseCase
	GetTransactionDetailUseCase() usecase.GetTransactionDetailUseCase
	GetAllTransactionDetailUseCase() usecase.GetAllTransactionUseCase
	GetBalanceUserUseCase() usecase.GetBalanceUserUseCase
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

func (u *useCaseManager) AddTransactionDetailUseCase() usecase.AddTransactionDetailUseCase {
	return usecase.NewAddTransactionDetailUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) GetTransactionDetailUseCase() usecase.GetTransactionDetailUseCase {
	return usecase.NewGetTransactionDetailUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) GetAllTransactionDetailUseCase() usecase.GetAllTransactionUseCase {
	return usecase.NewGetAllTransactionUseCase(u.repo.CustomerRepo())
}

func (u *useCaseManager) GetBalanceUserUseCase() usecase.GetBalanceUserUseCase {
	return usecase.NewGetBalanceUserUseCase(u.repo.CustomerRepo())
}

func NewUseCaseManager(repo RepoManager) UseCaseManager {
	return &useCaseManager{
		repo: repo,
	}
}
