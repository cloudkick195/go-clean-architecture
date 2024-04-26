package deposit

import (
	"go_clean_architecture/commons"
	"go_clean_architecture/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	List() gin.HandlerFunc
	UpdateDeposit() gin.HandlerFunc
	Detail() gin.HandlerFunc
}

type handler struct {
	useCase IUsecase
}

func NewHandler(useCase IUsecase) IHandler {
	return &handler{
		useCase: useCase,
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

		output, err := h.useCase.List(c, &params.Pagination, &params.FilterDeposit)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.NewSuccessResponse(output, &params.Pagination, &params.FilterDeposit))
	}
}

func (h *handler) Detail() gin.HandlerFunc {
	return func(c *gin.Context) {
		IdStr := c.Param("id")
		Id, err := strconv.ParseUint(IdStr, 10, 64)
		if err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		result, err := h.useCase.Detail(c, uint(Id))
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(result))
	}
}

func (h *handler) UpdateDeposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		IdStr := c.Param("id")
		IdInt64, err := strconv.ParseUint(IdStr, 10, 64)
		Id := uint(IdInt64)
		if err != nil {
			panic(commons.ErrInvalidRequest(err))
		}

		if err := h.useCase.UpdateDeposit(c, Id); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, commons.SimpleSuccessResponse(nil))
	}
}
