package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/algorithm9/flash-deal/pkg/errorx"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errorx.CodeOK.Int(),
		Message: "success",
		Data:    data,
	})
}

// Fail is a unified handler to convert error to errorx.Error and respond
func Fail(c *gin.Context, err error) {
	if e, ok := err.(errorx.Error); ok {
		Respond(c, e)
	} else {
		Respond(c, errorx.Internal(err.Error()))
	}
}

func JSON(c *gin.Context, e errorx.Error) {
	c.JSON(e.HTTPStatus(), Response{
		Code:    e.Code(),
		Message: e.Message(),
		Data:    nil,
	})
}

func Respond(c *gin.Context, e errorx.Error) {
	JSON(c, e.Log())
}
