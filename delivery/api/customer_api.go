package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"mohamadelabror.me/simplebankappv2/delivery/appreq"
	"mohamadelabror.me/simplebankappv2/delivery/jwt"
	"mohamadelabror.me/simplebankappv2/delivery/response"
	"mohamadelabror.me/simplebankappv2/usecase"
	"mohamadelabror.me/simplebankappv2/util"
)

type CustomerApi struct {
	registerUseCase                usecase.RegisterAccountUseCase
	loginUseCase                   usecase.LoginUseCase
	logoutUseCase                  usecase.LogoutUseCase
	transferUseCase                usecase.TransferUseCase
	addLogUseCase                  usecase.AddLogUseCase
	addTransactionDetailUsecase    usecase.AddTransactionDetailUseCase
	getTransactionDetailUseCase    usecase.GetTransactionDetailUseCase
	getAllTransactionDetailUseCase usecase.GetAllTransactionUseCase
	getBalanceUserUseCase          usecase.GetBalanceUserUseCase
}

func (cu *CustomerApi) UserRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newAcc appreq.RegisterReq
		response := response.NewResponse(c)

		err := c.ShouldBindJSON(&newAcc)
		if err != nil {
			response.NewErrorResponse(400, "Cant bind json", nil)
			return
		}
		newAccountNumber := util.GenerateAccountNumber()
		err = cu.registerUseCase.RegisterAccount(util.GenerateUuid(), newAccountNumber, newAcc.UserName, newAcc.UserPin, newAcc.UserPassword, newAcc.Balance)

		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		response.NewSuccesMessage(200, "succes", gin.H{"account_number": newAccountNumber})
	}
}

func (cu *CustomerApi) UserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var credential jwt.Credential
		response := response.NewResponse(c)

		err := c.ShouldBindJSON(&credential)
		if err != nil {
			response.NewErrorResponse(400, "Cant bind json", nil)
			return
		}
		jwtToken, err := jwt.GenerateToken(credential.Username, "luxamrown@corp.id")
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		selectedId, err := cu.loginUseCase.Login(credential.Username, credential.Password, jwtToken)
		if err != nil {
			response.NewErrorResponse(401, err.Error(), nil)
			return
		}
		response.NewSuccesMessage(200, "succes", gin.H{"token": jwtToken, "id": selectedId})
	}
}

func (cu *CustomerApi) UserLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var credential appreq.LogoutReq
		response := response.NewResponse(c)

		err := c.ShouldBindJSON(&credential)
		if err != nil {
			response.NewErrorResponse(400, "Cant bind json", nil)
			return
		}
		err = cu.logoutUseCase.Logout(credential.Id)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		response.NewSuccesMessage(200, "succes", nil)
	}
}

func (cu *CustomerApi) UserTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transferReq appreq.TransactionReq
		timeNow := time.Now().Format("2006-01-02 15:04")
		idHistory := util.GenerateUuid()
		idTransactionDetails := util.GenerateUuid()
		err := c.BindJSON(&transferReq)
		response := response.NewResponse(c)
		if err != nil {
			response.NewErrorResponse(400, "Cant bind json", nil)
			return
		}
		err = cu.transferUseCase.Transfer(transferReq.SenderId, transferReq.SenderPin, transferReq.SenderAccNumber, transferReq.ReceiverAccountNumber, transferReq.Amount, transferReq.IsMerchant)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		err = cu.addLogUseCase.AddLog(idHistory, transferReq.SenderAccNumber, transferReq.ReceiverAccountNumber, timeNow, transferReq.IsMerchant)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		err = cu.addTransactionDetailUsecase.AddTransactionDetail(idTransactionDetails, idHistory, transferReq.SenderId, transferReq.Message, transferReq.Amount)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		newTransactionDetail, err := cu.getTransactionDetailUseCase.GetTransactionDetail(idTransactionDetails, idHistory, transferReq.IsMerchant)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		response.NewSuccesMessage(200, "succes", newTransactionDetail)
	}
}

func (cu *CustomerApi) GetAllTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		response := response.NewResponse(c)

		Transactions, err := cu.getAllTransactionDetailUseCase.GetAllTransactionDetail(userId)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		response.NewSuccesMessage(200, "succes", Transactions)
	}
}

func (cu *CustomerApi) GetBalanceUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionReq appreq.TransactionReq
		response := response.NewResponse(c)

		err := c.BindJSON(&transactionReq)
		if err != nil {

			response.NewErrorResponse(400, "Cant bind json", nil)
			return
		}
		balance, err := cu.getBalanceUserUseCase.GetBalance(transactionReq.SenderId, transactionReq.SenderPin)
		if err != nil {
			response.NewErrorResponse(400, err.Error(), nil)
			return
		}
		response.NewSuccesMessage(200, "succes", gin.H{"balance": balance})

	}
}

func NewCustomerApi(customerRoute *gin.RouterGroup, registerUseCase usecase.RegisterAccountUseCase, loginUseCase usecase.LoginUseCase, logoutUsecase usecase.LogoutUseCase, transferUsecase usecase.TransferUseCase, addLogUseCase usecase.AddLogUseCase, addTransactionDetail usecase.AddTransactionDetailUseCase, getTransactionDetaiUseCase usecase.GetTransactionDetailUseCase, getAllTransactionDetailUseCase usecase.GetAllTransactionUseCase, getBalanceUserUseCase usecase.GetBalanceUserUseCase) {
	api := CustomerApi{
		registerUseCase:                registerUseCase,
		loginUseCase:                   loginUseCase,
		logoutUseCase:                  logoutUsecase,
		transferUseCase:                transferUsecase,
		addLogUseCase:                  addLogUseCase,
		addTransactionDetailUsecase:    addTransactionDetail,
		getTransactionDetailUseCase:    getTransactionDetaiUseCase,
		getAllTransactionDetailUseCase: getAllTransactionDetailUseCase,
		getBalanceUserUseCase:          getBalanceUserUseCase,
	}
	customerRoute.POST("/register", api.UserRegister())
	customerRoute.POST("/login", api.UserLogin())
	customerRoute.POST("/logout", api.UserLogout())
	customerRoute.POST("/transfer", api.UserTransfer())
	customerRoute.GET("/transfer/:userId", api.GetAllTransaction())
	customerRoute.POST("/transfer/getbalance", api.GetBalanceUser())
}
