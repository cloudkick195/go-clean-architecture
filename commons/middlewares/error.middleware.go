package middlewares

import (
	"bytes"
	"encoding/json"
	"go_clean_architecture/commons"
	"go_clean_architecture/modules/log"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middlewareManager) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStream, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqStream))

		defer func() {
			if err := recover(); err != nil {
				var appErr *commons.AppError
				if errConvert, ok := err.(*commons.AppError); ok {
					appErr = errConvert
				} else {
					appErr = commons.ErrInternal(err.(error))
				}

				appErrStr, _ := json.Marshal(appErr)
				//reqStr, _ := json.Marshal(c.Request)

				m.logUsecase.Create(c, &log.CreateInput{
					StatusCode: appErr.StatusCode,
					RootErr:    string(appErrStr),
					Message:    appErr.Message,
					Log:        appErr.Log,
					Key:        appErr.Key,
					Api:        c.FullPath(),
					Request:    strings.ReplaceAll(string(reqStream), " ", ""),
					Ip:         c.ClientIP() + "-" + c.Request.Header.Get("Secure-Token"),
				})
				c.AbortWithStatusJSON(appErr.StatusCode, log.Output{
					StatusCode: appErr.StatusCode,
					Message:    appErr.Message,
					Key:        appErr.Key,
				})
				return
			}
		}()
		c.Next()
	}
}
