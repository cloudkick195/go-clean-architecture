package server

import (
	"go_clean_architecture/commons/middlewares"
	"go_clean_architecture/commons/uows"
	"go_clean_architecture/modules/auth"
	"go_clean_architecture/modules/comission"
	configModule "go_clean_architecture/modules/config"
	"go_clean_architecture/modules/deposit"
	depositProvider "go_clean_architecture/modules/deposit/provider"
	"go_clean_architecture/modules/depositLog"
	"go_clean_architecture/modules/log"
	"go_clean_architecture/modules/member"

	"go_clean_architecture/modules/transfer"
	transferProvider "go_clean_architecture/modules/transfer/provider"
	"go_clean_architecture/modules/transferLog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type API interface {
	RegisterRoutes(router *gin.RouterGroup, mw middlewares.IMiddlewareManager)
}

type IPublicApi interface {
	PublicRoutes(router *gin.RouterGroup, mw middlewares.IMiddlewareManager)
}

func (s *Server) InitRouter() {
	routerDefault := s.router

	routerDefault.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "API endpoint not found",
		})
	})
	router := routerDefault.Group("/api/v1")
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to my API!",
		})
	})

	// repository
	repoUOW := uows.NewRepoUnitOfWork(s.db)
	//Usecase
	logUsecase := log.NewUsecase(repoUOW)
	middlewareManager := middlewares.NewMiddlewareManager(logUsecase, s.db)

	//Factory struct
	providerFactory := transferProvider.NewProviderFactory(repoUOW, middlewareManager)
	providerDepositFactory := &depositProvider.ProviderFactory{}
	depositStrategyFactory := deposit.DepositStrategyFactory{}

	comissionUsecase := comission.NewUsecase(repoUOW)

	configUsecase := configModule.NewUsecase(repoUOW)
	transferLogUsecase := transferLog.NewUsecase(repoUOW)
	depositLogUsecase := depositLog.NewUsecase(repoUOW)
	transferUsecase := transfer.NewUsecase(repoUOW, transferLogUsecase, configUsecase, providerFactory)
	depositUsecase := deposit.NewUsecase(repoUOW, depositLogUsecase, configUsecase, providerDepositFactory, depositStrategyFactory)
	memberUsecase := member.NewUsecase(repoUOW, depositUsecase, transferUsecase, comissionUsecase)

	authUsecase := auth.NewUsecase(memberUsecase)

	//handler
	transferHandler := transfer.NewHandler(transferUsecase)
	depositHandler := deposit.NewHandler(depositUsecase)
	authHandler := auth.NewHandler(authUsecase)
	memberHandler := member.NewHandler(memberUsecase, middlewareManager)

	router.Use(middlewareManager.PerClientRateLimiterForAll())
	router.Use(middlewareManager.ErrorMiddleware())
	//admin group

	authApis := []IPublicApi{
		&transfer.API{Handler: transferHandler},
		&auth.API{Handler: authHandler},
	}
	for _, a := range authApis {
		a.PublicRoutes(router, middlewareManager)
	}

	adminGroup := router.Group("/auth")
	adminGroup.Use(middlewareManager.AuthMiddleware())
	apis := []API{
		&member.API{Handler: memberHandler},
		&deposit.API{Handler: depositHandler},
	}

	for _, a := range apis {
		a.RegisterRoutes(adminGroup, middlewareManager)
	}
}
