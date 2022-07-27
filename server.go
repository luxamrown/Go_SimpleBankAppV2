package main

import (
	"github.com/gin-gonic/gin"
	"mohamadelabror.me/simplebankappv2/config"
	"mohamadelabror.me/simplebankappv2/delivery/api"
	"mohamadelabror.me/simplebankappv2/delivery/jwt"
)

type AppServer interface {
	Run()
}
type appServer struct {
	routerEngine *gin.Engine
	cfg          config.Config
}

func (a *appServer) initHandlers() {
	a.v1()
}

func (a *appServer) v1() {
	a.routerEngine.Use(jwt.AuthTokenMiddleware())
	customerApiGroup := a.routerEngine.Group("/bank")
	api.NewCustomerApi(customerApiGroup, a.cfg.UseCaseManager.RegisterAccountUseCase(), a.cfg.UseCaseManager.LoginUseCase(), a.cfg.UseCaseManager.LogoutUseCase(), a.cfg.UseCaseManager.TransferUseCase(), a.cfg.UseCaseManager.AddLogUseCase())
}

func (a *appServer) Run() {
	a.initHandlers()
	err := a.routerEngine.Run(a.cfg.ApiConfig.Url)
	if err != nil {
		panic(err)
	}
}

func Server() AppServer {
	r := gin.Default()
	c := config.NewConfig(".", "config")
	return &appServer{
		routerEngine: r,
		cfg:          c,
	}
}
