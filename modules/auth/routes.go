package auth

import (
	"go_clean_architecture/commons/middlewares"

	"github.com/gin-gonic/gin"
)

type API struct{ Handler IHandler }

// Map auth routes
func (a *API) PublicRoutes(router *gin.RouterGroup, mw middlewares.IMiddlewareManager) {
	h := a.Handler
	router.POST("/register", h.Register())
	router.POST("/login", h.Login())
}
