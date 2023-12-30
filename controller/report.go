package controller

import (
	"GinProject/model"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
)

type ReportController struct {
}

func Report() ReportController {
	return ReportController{}
}

func (r ReportController) AddReport(c *gin.Context) {
	report := &model.Report{}
	err := c.BindJSON(report)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.AddReport(report)
	response.RspSuccess(c, res)
}

func (r ReportController) JudgeReport(c *gin.Context) {
	report := &model.Report{}
	err := c.BindJSON(report)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.JudgeReport(report)
	response.RspSuccess(c, res)
}
