package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"mohamadelabror.me/simplebankappv2/delivery/appreq"
	"mohamadelabror.me/simplebankappv2/delivery/jwt"
	"mohamadelabror.me/simplebankappv2/usecase"
	"mohamadelabror.me/simplebankappv2/util"
)

type CustomerApi struct {
	registerUseCase usecase.RegisterAccountUseCase
	loginUseCase    usecase.LoginUseCase
	logoutUseCase   usecase.LogoutUseCase
	transferUseCase usecase.TransferUseCase
	addLogUseCase   usecase.AddLogUseCase
}

func (cu *CustomerApi) UserRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newAcc appreq.RegisterReq
		err := c.ShouldBindJSON(&newAcc)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "cant bind json",
			})
			return
		}
		newAccountNumber := util.GenerateAccountNumber()
		err = cu.registerUseCase.RegisterAccount(util.GenerateUuid(), newAccountNumber, newAcc.UserName, newAcc.UserPin, newAcc.UserPassword, newAcc.Balance)
		if err != nil {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message":        "success",
			"account_number": newAccountNumber,
		})
	}
}

func (cu *CustomerApi) UserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var credential jwt.Credential
		err := c.ShouldBindJSON(&credential)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "cant bind json",
			})
			return
		}
		jwtToken, err := jwt.GenerateToken(credential.Username, "luxamrown@corp.id")
		if err != nil {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		selectedId, err := cu.loginUseCase.Login(credential.Username, credential.Password, jwtToken)
		if err != nil {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"token":   jwtToken,
			"id":      selectedId,
		})
	}
}

func (cu *CustomerApi) UserLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var credential appreq.LogoutReq
		err := c.ShouldBindJSON(&credential)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "cant bind json",
			})
			return
		}
		err = cu.logoutUseCase.Logout(credential.Id)
		if err != nil {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	}
}

func (cu *CustomerApi) UserTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transferReq appreq.TransferReq
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		err := c.BindJSON(&transferReq)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "cant bind json",
			})
			return
		}
		err = cu.transferUseCase.Transfer(transferReq.SenderId, transferReq.SenderPin, transferReq.SenderAccNumber, transferReq.ReceiverAccountNumber, transferReq.Amount, transferReq.IsMerchant)
		if err != nil {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		err = cu.addLogUseCase.AddLog(util.GenerateUuid(), transferReq.SenderAccNumber, transferReq.ReceiverAccountNumber, timeNow, transferReq.IsMerchant)
		if err != nil {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
		})
	}
}

// func (cu *CustomerApi) UserTest() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var req appreq.RequestTest
// 		err := c.BindJSON(&req)
// 		if err != nil {
// 			c.JSON(401, gin.H{
// 				"message": "cant bind json",
// 			})
// 			return
// 		}
// 		// if req.Lain == nil {
// 		// 	c.JSON(200, "lain kosong")
// 		// 	return
// 		// }
// 		c.JSON(200, req)

// 	}
// }

func NewCustomerApi(customerRoute *gin.RouterGroup, registerUseCase usecase.RegisterAccountUseCase, loginUseCase usecase.LoginUseCase, logoutUsecase usecase.LogoutUseCase, transferUsecase usecase.TransferUseCase, addLogUseCase usecase.AddLogUseCase) {
	api := CustomerApi{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		logoutUseCase:   logoutUsecase,
		transferUseCase: transferUsecase,
		addLogUseCase:   addLogUseCase,
	}
	customerRoute.POST("/register", api.UserRegister())
	customerRoute.POST("/login", api.UserLogin())
	customerRoute.POST("/logout", api.UserLogout())
	customerRoute.POST("/transfer", api.UserTransfer())
}
