package middlewares

import (
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middlewareManager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header của request.
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			panic(commons.ErrorUnAuthorized())
		}
		tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
		if tokenStr == "" {
			panic(commons.ErrorUnAuthorized())
		}
		jwtProvider := commons.NewJwtProvider()
		payload, err := jwtProvider.Validate(tokenStr)

		if err != nil {
			panic(commons.ErrorUnAuthorized())
		}
		var input *models.Member

		if err := m.db.WithContext(c).Where(payload.Id).First(&input).Error; err != nil {
			panic(err)
		}
		if input == nil {
			panic(commons.ErrorUnAuthorized())
		}
		if input.Token != tokenStr {
			panic(commons.ErrorUnAuthorized())
		}
		if input.Status != models.MemberStatusActive {
			panic(commons.ErrorUnAuthorized())
		}

		c.Set("member", input)

		c.Next()
	}
}

func (m *middlewareManager) PermissionMiddleware(roleStr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		if member.Role != roleStr {
			panic(commons.ErrorUnAuthorized())
		}
		c.Next()
	}
}
