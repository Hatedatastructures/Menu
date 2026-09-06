package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigSystemApi struct{}

// GetSystemConfigHandler
// @Tags systemRbacconfigSystemApi
// @Summary GetSystemConfigHandler 获取系统配置-后台使用
// @Description GetSystemConfigHandler 获取系统配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetSystemConfigRes}
// @Router /api/config/system [GET]
func (s *ConfigSystemApi) GetSystemConfigHandler(c *gin.Context) {
	data, err := configSystemService.GetSystemConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveSystemConfigHandler
// @Tags systemRbacconfigSystemApi
// @Summary SaveSystemConfigHandler 保存系统配置-后台使用
// @Description SaveSystemConfigHandler 保存系统配置-后台使用
// @Param data body req.SaveSystemConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveSystemConfigRes}
// @Router /api/config/system [POST]
func (s *ConfigSystemApi) SaveSystemConfigHandler(c *gin.Context) {
	var req req.SaveSystemConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configSystemService.SaveSystemConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
