package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type EmailSendApi struct{}

// SendEmailHandler
// @Tags systemRbacemailSendApi
// @Summary SendEmailHandler 发送普通邮件-后台使用
// @Description SendEmailHandler 发送普通邮件-后台使用
// @Param data body req.SendEmailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SendEmailRes}
// @Router /email/send [POST]
func (s *EmailSendApi) SendEmailHandler(c *gin.Context) {
	var req req.SendEmailReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailSendService.SendEmail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SendTemplateEmailHandler
// @Tags systemRbacemailSendApi
// @Summary SendTemplateEmailHandler 发送模板邮件-后台使用
// @Description SendTemplateEmailHandler 发送模板邮件-后台使用
// @Param data body req.SendTemplateEmailReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/send-template [POST]
func (s *EmailSendApi) SendTemplateEmailHandler(c *gin.Context) {
	var req req.SendTemplateEmailReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := emailSendService.SendTemplateEmail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// SendBatchEmailHandler
// @Tags systemRbacemailSendApi
// @Summary SendBatchEmailHandler 批量发送邮件-后台使用
// @Description SendBatchEmailHandler 批量发送邮件-后台使用
// @Param data body req.SendBatchEmailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SendBatchEmailRes}
// @Router /email/send-batch [POST]
func (s *EmailSendApi) SendBatchEmailHandler(c *gin.Context) {
	var req req.SendBatchEmailReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailSendService.SendBatchEmail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
