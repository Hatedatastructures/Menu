package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigRegisterApi struct{}

// GetRegisterConfigHandler
// @Tags systemRbacconfigRegisterApi
// @Summary GetRegisterConfigHandler 获取注册配置-后台使用
// @Description GetRegisterConfigHandler 获取注册配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetRegisterConfigRes}
// @Router /api/config/register [GET]
func (s *ConfigRegisterApi) GetRegisterConfigHandler(c *gin.Context) {
	data, err := configRegisterService.GetRegisterConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveRegisterConfigHandler
// @Tags systemRbacconfigRegisterApi
// @Summary SaveRegisterConfigHandler 保存注册配置-后台使用
// @Description SaveRegisterConfigHandler 保存注册配置-后台使用
// @Param data body req.SaveRegisterConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveRegisterConfigRes}
// @Router /api/config/register [POST]
func (s *ConfigRegisterApi) SaveRegisterConfigHandler(c *gin.Context) {
	var req req.SaveRegisterConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configRegisterService.SaveRegisterConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}