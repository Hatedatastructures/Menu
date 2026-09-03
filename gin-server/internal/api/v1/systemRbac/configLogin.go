package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigLoginApi struct{}

// GetLoginConfigHandler
// @Tags systemRbacconfigLoginApi
// @Summary GetLoginConfigHandler 获取登录配置-后台使用
// @Description GetLoginConfigHandler 获取登录配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetLoginConfigRes}
// @Router /api/config/login [GET]
func (s *ConfigLoginApi) GetLoginConfigHandler(c *gin.Context) {
	data, err := configLoginService.GetLoginConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveLoginConfigHandler
// @Tags systemRbacconfigLoginApi
// @Summary SaveLoginConfigHandler 保存登录配置-后台使用
// @Description SaveLoginConfigHandler 保存登录配置-后台使用
// @Param data body req.SaveLoginConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveLoginConfigRes}
// @Router /api/config/login [POST]
func (s *ConfigLoginApi) SaveLoginConfigHandler(c *gin.Context) {
	var req req.SaveLoginConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configLoginService.SaveLoginConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
