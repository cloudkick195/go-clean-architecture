package auth

import (
	"net/http"

	"go_clean_architecture/commons"
	"go_clean_architecture/modules/member"
	"go_clean_architecture/utils"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
}

type handler struct {
	Usecase IUsecase
}

func NewHandler(Usecase IUsecase) IHandler {
	return &handler{
		Usecase: Usecase,
	}
}

func (h *handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := &member.CreateInput{Ip: c.ClientIP()}
		if err := c.ShouldBindJSON(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		if err := h.Usecase.Register(c, params); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(nil))
	}
}

func (h *handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := &LoginInput{}
		if err := c.ShouldBindJSON(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		if err := utils.Validate.Struct(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		ouput, err := h.Usecase.Login(c, params)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(ouput))
	}
}
