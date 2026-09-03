package systemRbac

import (
	"strconv"

	"github.com/gin-gonic/gin"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"
	biz_err "shack/internal/error"
)

type ConfigGroupApi struct{}

// CreateConfigGroupHandler
// @Tags systemRbacconfigGroupApi
// @Summary CreateConfigGroupHandler 创建配置分组-后台使用
// @Description CreateConfigGroupHandler 创建配置分组-后台使用
// @Param data body req.CreateConfigGroupReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateConfigGroupRes}
// @Router /api/config/groups [POST]
func (s *ConfigGroupApi) CreateConfigGroupHandler(c *gin.Context) {
	var req req.CreateConfigGroupReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configGroupService.CreateConfigGroup(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetConfigGroupListHandler
// @Tags systemRbacconfigGroupApi
// @Summary GetConfigGroupListHandler 配置分组列表-后台使用
// @Description GetConfigGroupListHandler 配置分组列表-后台使用
// @Param data body req.GetConfigGroupListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetConfigGroupListRes}
// @Router /api/config/groups [GET]
func (s *ConfigGroupApi) GetConfigGroupListHandler(c *gin.Context) {
	var req req.GetConfigGroupListReq
	// query 参数
	{
		val := c.Query("page")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
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
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Size = parsed
		}
	}
	{
		val := c.Query("keyword")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Keyword = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configGroupService.GetConfigGroupList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetConfigGroupDetailHandler
// @Tags systemRbacconfigGroupApi
// @Summary GetConfigGroupDetailHandler 获取配置分组详情-后台使用
// @Description GetConfigGroupDetailHandler 获取配置分组详情-后台使用
// @Param data body req.GetConfigGroupDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetConfigGroupDetailRes}
// @Router /api/config/groups/:id [GET]
func (s *ConfigGroupApi) GetConfigGroupDetailHandler(c *gin.Context) {
	var req req.GetConfigGroupDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configGroupService.GetConfigGroupDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateConfigGroupHandler
// @Tags systemRbacconfigGroupApi
// @Summary UpdateConfigGroupHandler 更新配置分组-后台使用
// @Description UpdateConfigGroupHandler 更新配置分组-后台使用
// @Param data body req.UpdateConfigGroupReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/config/groups/:id [PUT]
func (s *ConfigGroupApi) UpdateConfigGroupHandler(c *gin.Context) {
	var req req.UpdateConfigGroupReq
	// path 参数
	{
		val := c.Param("id")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed
	}
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := configGroupService.UpdateConfigGroup(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteConfigGroupHandler
// @Tags systemRbacconfigGroupApi
// @Summary DeleteConfigGroupHandler 删除配置分组-后台使用
// @Description DeleteConfigGroupHandler 删除配置分组-后台使用
// @Param data body req.DeleteConfigGroupReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/config/groups/:id [DELETE]
func (s *ConfigGroupApi) DeleteConfigGroupHandler(c *gin.Context) {
	var req req.DeleteConfigGroupReq
	// path 参数
	{
		val := c.Param("id")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := configGroupService.DeleteConfigGroup(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// BatchDeleteConfigGroupsHandler
// @Tags systemRbacconfigGroupApi
// @Summary BatchDeleteConfigGroupsHandler 批量删除配置分组-后台使用
// @Description BatchDeleteConfigGroupsHandler 批量删除配置分组-后台使用
// @Param data body req.BatchDeleteConfigGroupsReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/config/groups [DELETE]
func (s *ConfigGroupApi) BatchDeleteConfigGroupsHandler(c *gin.Context) {
	var req req.BatchDeleteConfigGroupsReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := configGroupService.BatchDeleteConfigGroups(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetConfigByCodeHandler
// @Tags systemRbacconfigGroupApi
// @Summary GetConfigByCodeHandler 根据编码获取配置-后台使用
// @Description GetConfigByCodeHandler 根据编码获取配置-后台使用
// @Param data body req.GetConfigByCodeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetConfigByCodeRes}
// @Router /api/config/:code [GET]
func (s *ConfigGroupApi) GetConfigByCodeHandler(c *gin.Context) {
	var req req.GetConfigByCodeReq
	// path 参数
	{
		val := c.Param("code")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Code = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configGroupService.GetConfigByCode(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveConfigHandler
// @Tags systemRbacconfigGroupApi
// @Summary SaveConfigHandler 保存配置-后台使用
// @Description SaveConfigHandler 保存配置-后台使用
// @Param data body req.SaveConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveConfigRes}
// @Router /api/config/:code [POST]
func (s *ConfigGroupApi) SaveConfigHandler(c *gin.Context) {
	var req req.SaveConfigReq
	// path 参数
	{
		val := c.Param("code")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Code = parsed
	}
	{
		val := c.Param("code")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Code = parsed
	}
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configGroupService.SaveConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// RefreshConfigCacheHandler
// @Tags systemRbacconfigGroupApi
// @Summary RefreshConfigCacheHandler 刷新配置缓存-后台使用
// @Description RefreshConfigCacheHandler 刷新配置缓存-后台使用
// @Success 200 {object} vo.Result{data=_.RefreshConfigCacheRes}
// @Router /api/config/refresh [POST]
func (s *ConfigGroupApi) RefreshConfigCacheHandler(c *gin.Context) {
	data, err := configGroupService.RefreshConfigCache(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}