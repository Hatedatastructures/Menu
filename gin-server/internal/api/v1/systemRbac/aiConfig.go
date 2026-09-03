package systemRbac

import (
	"strconv"

	"github.com/gin-gonic/gin"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"
)

type AiConfigApi struct{}

// CreateAiConfigHandler
// @Tags systemRbacaiConfigApi
// @Summary CreateAiConfigHandler 创建AI配置-后台使用
// @Description CreateAiConfigHandler 创建AI配置-后台使用
// @Param data body req.CreateAiConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateAiConfigRes}
// @Router /api/ai/config [POST]
func (s *AiConfigApi) CreateAiConfigHandler(c *gin.Context) {
	var req req.CreateAiConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := aiConfigService.CreateAiConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetAiConfigListHandler
// @Tags systemRbacaiConfigApi
// @Summary GetAiConfigListHandler 获取AI配置列表-后台使用
// @Description GetAiConfigListHandler 获取AI配置列表-后台使用
// @Param data body req.GetAiConfigListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetAiConfigListRes}
// @Router /api/ai/config [GET]
func (s *AiConfigApi) GetAiConfigListHandler(c *gin.Context) {
	var req req.GetAiConfigListReq
	// query 参数
	{
		val := c.Query("page")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Page = parsed

		}
	}
	{
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Size = parsed

		}
	}
	{
		val := c.Query("name")
		if val != "" {
			req.Name = val

		}
	}
	{
		val := c.Query("enabled")
		if val != "" {
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Enabled = parsed

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := aiConfigService.GetAiConfigList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetAiConfigHandler
// @Tags systemRbacaiConfigApi
// @Summary GetAiConfigHandler 获取AI配置详情-后台使用
// @Description GetAiConfigHandler 获取AI配置详情-后台使用
// @Param data body req.GetAiConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetAiConfigRes}
// @Router /api/ai/config/:id [GET]
func (s *AiConfigApi) GetAiConfigHandler(c *gin.Context) {
	var req req.GetAiConfigReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := aiConfigService.GetAiConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateAiConfigHandler
// @Tags systemRbacaiConfigApi
// @Summary UpdateAiConfigHandler 更新AI配置-后台使用
// @Description UpdateAiConfigHandler 更新AI配置-后台使用
// @Param data body req.UpdateAiConfigReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/ai/config/:id [PUT]
func (s *AiConfigApi) UpdateAiConfigHandler(c *gin.Context) {
	var req req.UpdateAiConfigReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

	}
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := aiConfigService.UpdateAiConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteAiConfigHandler
// @Tags systemRbacaiConfigApi
// @Summary DeleteAiConfigHandler 删除AI配置-后台使用
// @Description DeleteAiConfigHandler 删除AI配置-后台使用
// @Param data body req.DeleteAiConfigReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/ai/config/:id [DELETE]
func (s *AiConfigApi) DeleteAiConfigHandler(c *gin.Context) {
	var req req.DeleteAiConfigReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := aiConfigService.DeleteAiConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// BatchDeleteAiConfigHandler
// @Tags systemRbacaiConfigApi
// @Summary BatchDeleteAiConfigHandler 批量删除AI配置-后台使用
// @Description BatchDeleteAiConfigHandler 批量删除AI配置-后台使用
// @Param data body req.BatchDeleteAiConfigReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/ai/config [DELETE]
func (s *AiConfigApi) BatchDeleteAiConfigHandler(c *gin.Context) {
	var req req.BatchDeleteAiConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := aiConfigService.BatchDeleteAiConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// TestAiConfigHandler
// @Tags systemRbacaiConfigApi
// @Summary TestAiConfigHandler 测试AI接口-后台使用
// @Description TestAiConfigHandler 测试AI接口-后台使用
// @Param data body req.TestAiConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.TestAiConfigRes}
// @Router /api/ai/test [POST]
func (s *AiConfigApi) TestAiConfigHandler(c *gin.Context) {
	var req req.TestAiConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := aiConfigService.TestAiConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// AiChatHandler
// @Tags systemRbacaiConfigApi
// @Summary AiChatHandler AI对话调用-后台使用
// @Description AiChatHandler AI对话调用-后台使用
// @Param data body req.AiChatReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.AiChatRes}
// @Router /api/ai/chat [POST]
func (s *AiConfigApi) AiChatHandler(c *gin.Context) {
	var req req.AiChatReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := aiConfigService.AiChat(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ToggleAiConfigStatusHandler
// @Tags systemRbacaiConfigApi
// @Summary ToggleAiConfigStatusHandler 切换AI配置启用状态-后台使用
// @Description ToggleAiConfigStatusHandler 切换AI配置启用状态-后台使用
// @Param data body req.ToggleAiConfigStatusReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/ai/config/:id/status [PUT]
func (s *AiConfigApi) ToggleAiConfigStatusHandler(c *gin.Context) {
	var req req.ToggleAiConfigStatusReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

	}
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := aiConfigService.ToggleAiConfigStatus(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetEnabledAiConfigsHandler
// @Tags systemRbacaiConfigApi
// @Summary GetEnabledAiConfigsHandler 获取启用的AI配置列表-前台使用
// @Description GetEnabledAiConfigsHandler 获取启用的AI配置列表-前台使用
// @Success 200 {object} vo.Result{data=_.GetEnabledAiConfigsRes}
// @Router /api/ai/config/enabled [GET]
func (s *AiConfigApi) GetEnabledAiConfigsHandler(c *gin.Context) {
	data, err := aiConfigService.GetEnabledAiConfigs(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}