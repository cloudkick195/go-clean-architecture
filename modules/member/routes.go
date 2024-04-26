package member

import (
	"go_clean_architecture/commons/middlewares"
	"go_clean_architecture/utils"

	"github.com/gin-gonic/gin"
)

type API struct{ Handler IHandler }

// Map auth routes
func (a *API) RegisterRoutes(router *gin.RouterGroup, mw middlewares.IMiddlewareManager) {
	h := a.Handler
	router.Group("/member").
		GET("/", mw.PermissionMiddleware(utils.CONSTANT_ADMIN), h.List()).
		GET("/:id", mw.PermissionMiddleware(utils.CONSTANT_ADMIN), h.Detail()).
		PUT("/:id", mw.PermissionMiddleware(utils.CONSTANT_ADMIN), h.UpdateMember()).
		GET("/profile", h.Get()).
		GET("/transaction-history", h.TransactionHistory()).
		POST("/deposit", h.Deposit()).
		PUT("/password", h.UpdatePassword())

}
