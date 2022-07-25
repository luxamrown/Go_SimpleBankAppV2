package api

import (
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
		err = cu.registerUseCase.RegisterAccount(util.GenerateUuid(), newAcc.AccountNumber, newAcc.UserName, newAcc.UserPassword, newAcc.Balance)
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

func NewCustomerApi(customerRoute *gin.RouterGroup, registerUseCase usecase.RegisterAccountUseCase, loginUseCase usecase.LoginUseCase, logoutUsecase usecase.LogoutUseCase) {
	api := CustomerApi{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		logoutUseCase:   logoutUsecase,
	}
	customerRoute.POST("/register", api.UserRegister())
	customerRoute.POST("/login", api.UserLogin())
	customerRoute.POST("/logout", api.UserLogout())
}
