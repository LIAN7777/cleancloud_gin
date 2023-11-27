package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RspData struct {
	Code    RspCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty当data为空时,不展示这个字段
}

func RspError(ctx *gin.Context, c RspCode) {
	rd := &RspData{
		Code:    c,
		Message: c.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func RspErrorWithMsg(ctx *gin.Context, code RspCode, data interface{}) {
	rd := &RspData{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func RspSuccess(ctx *gin.Context, data interface{}) {
	rd := &RspData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}
