package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigPasswordApi struct{}

// GetPasswordConfigHandler
// @Tags systemRbacconfigPasswordApi
// @Summary GetPasswordConfigHandler 获取密码配置-后台使用
// @Description GetPasswordConfigHandler 获取密码配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetPasswordConfigRes}
// @Router /api/config/password [GET]
func (s *ConfigPasswordApi) GetPasswordConfigHandler(c *gin.Context) {
	data, err := configPasswordService.GetPasswordConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SavePasswordConfigHandler
// @Tags systemRbacconfigPasswordApi
// @Summary SavePasswordConfigHandler 保存密码配置-后台使用
// @Description SavePasswordConfigHandler 保存密码配置-后台使用
// @Param data body req.SavePasswordConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SavePasswordConfigRes}
// @Router /api/config/password [POST]
func (s *ConfigPasswordApi) SavePasswordConfigHandler(c *gin.Context) {
	var req req.SavePasswordConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configPasswordService.SavePasswordConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
