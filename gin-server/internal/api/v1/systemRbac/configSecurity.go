package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigSecurityApi struct{}

// GetSecurityConfigHandler
// @Tags systemRbacconfigSecurityApi
// @Summary GetSecurityConfigHandler 获取安全配置-后台使用
// @Description GetSecurityConfigHandler 获取安全配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetSecurityConfigRes}
// @Router /api/config/security [GET]
func (s *ConfigSecurityApi) GetSecurityConfigHandler(c *gin.Context) {
	data, err := configSecurityService.GetSecurityConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveSecurityConfigHandler
// @Tags systemRbacconfigSecurityApi
// @Summary SaveSecurityConfigHandler 保存安全配置-后台使用
// @Description SaveSecurityConfigHandler 保存安全配置-后台使用
// @Param data body req.SaveSecurityConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveSecurityConfigRes}
// @Router /api/config/security [POST]
func (s *ConfigSecurityApi) SaveSecurityConfigHandler(c *gin.Context) {
	var req req.SaveSecurityConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configSecurityService.SaveSecurityConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GenerateRSAKeysHandler
// @Tags systemRbacconfigSecurityApi
// @Summary GenerateRSAKeysHandler 生成RSA密钥对-后台使用
// @Description GenerateRSAKeysHandler 生成RSA密钥对-后台使用
// @Success 200 {object} vo.Result{data=_.GenerateRSAKeysRes}
// @Router /api/config/security/generate-keys [POST]
func (s *ConfigSecurityApi) GenerateRSAKeysHandler(c *gin.Context) {
	data, err := configSecurityService.GenerateRSAKeys(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}