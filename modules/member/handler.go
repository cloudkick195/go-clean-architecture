package member

import (
	"errors"
	"net/http"
	"strconv"

	"go_clean_architecture/commons"
	"go_clean_architecture/commons/middlewares"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/utils"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	Get() gin.HandlerFunc
	Detail() gin.HandlerFunc
	Deposit() gin.HandlerFunc
	UpdatePassword() gin.HandlerFunc
	UpdateMember() gin.HandlerFunc
	TransactionHistory() gin.HandlerFunc
	List() gin.HandlerFunc
}

type handler struct {
	Usecase IUsecase
	mw      middlewares.IMiddlewareManager
}

func NewHandler(Usecase IUsecase, mw middlewares.IMiddlewareManager) IHandler {
	return &handler{
		Usecase: Usecase,
		mw:      mw,
	}
}

func (h *handler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		resultMember, err := h.Usecase.Profile(c, member.ID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(resultMember))
	}
}

func (h *handler) Detail() gin.HandlerFunc {
	return func(c *gin.Context) {
		IdStr := c.Param("id")
		Id, err := strconv.ParseUint(IdStr, 10, 64)
		if err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		resultMember, err := h.Usecase.Profile(c, uint(Id))
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(resultMember))
	}
}

func (h *handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := &Query{}
		if err := c.ShouldBindQuery(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		if err := utils.Validate.Struct(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		output, err := h.Usecase.List(c, &params.Pagination, &params.FilterMember)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.NewSuccessResponse(output, &params.Pagination, &params.FilterMember))
	}
}

func (h *handler) TransactionHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		params := &TransactionHistoryQuery{}
		if err := c.ShouldBindQuery(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		if err := utils.Validate.Struct(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		params.Id = member.ID
		output, err := h.Usecase.TransactionHistory(c, params)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(output))
	}
}

func (h *handler) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		params := &UpdatePasswordInput{}
		if err := c.ShouldBindJSON(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		if err := utils.Validate.Struct(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		params.Id = member.ID
		h.mw.SingleRequestAddUser(member.ID)
		err := h.Usecase.UpdatePassword(c, params)
		h.mw.SingleRequestReleaseUser(member.ID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(nil))
	}
}

func (h *handler) UpdateMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		IdStr := c.Param("id")
		IdInt64, err := strconv.ParseUint(IdStr, 10, 64)
		Id := uint(IdInt64)
		if err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		params := &UpdateMember{}
		if err := c.ShouldBindJSON(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		if err := utils.Validate.Struct(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		params.Id = Id
		h.mw.SingleRequestAddUser(Id)
		err = h.Usecase.UpdateMember(c, params)
		h.mw.SingleRequestReleaseUser(Id)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(nil))
	}
}

func (h *handler) Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		params := &DepositInput{}

		if err := c.ShouldBindJSON(params); err != nil {
			panic(commons.ErrInvalidRequest(err))
		}
		h.mw.SingleRequestAddUser(member.ID)
		err := h.Usecase.Deposit(c, member.ID, params, h.mw)
		h.mw.SingleRequestReleaseUser(member.ID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(nil))
	}
}
