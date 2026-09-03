package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigEmailApi struct{}

// GetEmailConfigHandler
// @Tags systemRbacconfigEmailApi
// @Summary GetEmailConfigHandler 获取邮件配置-后台使用
// @Description GetEmailConfigHandler 获取邮件配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetEmailConfigRes}
// @Router /api/config/email [GET]
func (s *ConfigEmailApi) GetEmailConfigHandler(c *gin.Context) {
	data, err := configEmailService.GetEmailConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveEmailConfigHandler
// @Tags systemRbacconfigEmailApi
// @Summary SaveEmailConfigHandler 保存邮件配置-后台使用
// @Description SaveEmailConfigHandler 保存邮件配置-后台使用
// @Param data body req.SaveEmailConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveEmailConfigRes}
// @Router /api/config/email [POST]
func (s *ConfigEmailApi) SaveEmailConfigHandler(c *gin.Context) {
	var req req.SaveEmailConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configEmailService.SaveEmailConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// TestEmailHandler
// @Tags systemRbacconfigEmailApi
// @Summary TestEmailHandler 测试邮件发送-后台使用
// @Description TestEmailHandler 测试邮件发送-后台使用
// @Param data body req.TestEmailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.TestEmailRes}
// @Router /api/config/email/test [POST]
func (s *ConfigEmailApi) TestEmailHandler(c *gin.Context) {
	var req req.TestEmailReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configEmailService.TestEmail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
