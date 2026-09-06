package gen

import (

	"github.com/gin-gonic/gin"
	req "shack/internal/model/gen/request"
	_ "shack/internal/model/gen/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"

)

type ApikeyApi struct{}

// GetMyApiKeyHandler
// @Tags genapikeyApi
// @Summary GetMyApiKeyHandler 获取我的API Key配置
// @Description GetMyApiKeyHandler 获取我的API Key配置
// @Success 200 {object} vo.Result{data=_.GetMyApiKeyRes}
// @Router /gen/apikey [GET]
func (s *ApikeyApi) GetMyApiKeyHandler(c *gin.Context) {
	data, err := apikeyService.GetMyApiKey(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveApiKeyHandler
// @Tags genapikeyApi
// @Summary SaveApiKeyHandler 保存/更新API Key
// @Description SaveApiKeyHandler 保存/更新API Key
// @Param data body req.SaveApiKeyReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveApiKeyRes}
// @Router /gen/apikey [POST]
func (s *ApikeyApi) SaveApiKeyHandler(c *gin.Context) {
	var req req.SaveApiKeyReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := apikeyService.SaveApiKey(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteApiKeyHandler
// @Tags genapikeyApi
// @Summary DeleteApiKeyHandler 删除API Key
// @Description DeleteApiKeyHandler 删除API Key
// @Success 200 {object} vo.Result{}
// @Router /gen/apikey [DELETE]
func (s *ApikeyApi) DeleteApiKeyHandler(c *gin.Context) {
	err := apikeyService.DeleteApiKey(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// TestApiKeyHandler
// @Tags genapikeyApi
// @Summary TestApiKeyHandler 测试API Key是否可用
// @Description TestApiKeyHandler 测试API Key是否可用
// @Param data body req.TestApiKeyReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.TestApiKeyRes}
// @Router /gen/apikey/test [POST]
func (s *ApikeyApi) TestApiKeyHandler(c *gin.Context) {
	var req req.TestApiKeyReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := apikeyService.TestApiKey(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}