package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type DictApi struct{}

// GetDictionaryListHandler
// @Tags systemRbacdictApi
// @Summary GetDictionaryListHandler 获取字典列表-后台使用
// @Description GetDictionaryListHandler 获取字典列表-后台使用
// @Success 200 {object} vo.Result{data=_.GetDictionaryListRes}
// @Router /api/dictionaries [GET]
func (s *DictApi) GetDictionaryListHandler(c *gin.Context) {
	data, err := dictService.GetDictionaryList(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CreateDictionaryHandler
// @Tags systemRbacdictApi
// @Summary CreateDictionaryHandler 创建字典-后台使用
// @Description CreateDictionaryHandler 创建字典-后台使用
// @Param data body req.CreateDictionaryReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateDictionaryRes}
// @Router /api/dictionaries [POST]
func (s *DictApi) CreateDictionaryHandler(c *gin.Context) {
	var req req.CreateDictionaryReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictService.CreateDictionary(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDictionaryHandler
// @Tags systemRbacdictApi
// @Summary GetDictionaryHandler 字典详情-后台使用
// @Description GetDictionaryHandler 字典详情-后台使用
// @Param data body req.GetDictionaryReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryRes}
// @Router /api/dictionaries/:id [GET]
func (s *DictApi) GetDictionaryHandler(c *gin.Context) {
	var req req.GetDictionaryReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictService.GetDictionary(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateDictionaryHandler
// @Tags systemRbacdictApi
// @Summary UpdateDictionaryHandler 更新字典-后台使用
// @Description UpdateDictionaryHandler 更新字典-后台使用
// @Param data body req.UpdateDictionaryReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/dictionaries/:id [PUT]
func (s *DictApi) UpdateDictionaryHandler(c *gin.Context) {
	var req req.UpdateDictionaryReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = parsed
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
	err := dictService.UpdateDictionary(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteDictionaryHandler
// @Tags systemRbacdictApi
// @Summary DeleteDictionaryHandler 删除字典-后台使用
// @Description DeleteDictionaryHandler 删除字典-后台使用
// @Param data body req.DeleteDictionaryReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/dictionaries/:id [DELETE]
func (s *DictApi) DeleteDictionaryHandler(c *gin.Context) {
	var req req.DeleteDictionaryReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := dictService.DeleteDictionary(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetDictsByTypesHandler
// @Tags systemRbacdictApi
// @Summary GetDictsByTypesHandler 批量字典查询-前台/后台使用
// @Description GetDictsByTypesHandler 批量字典查询-前台/后台使用
// @Param data body req.GetDictsByTypesReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictsByTypesRes}
// @Router /api/dicts [GET]
func (s *DictApi) GetDictsByTypesHandler(c *gin.Context) {
	var req req.GetDictsByTypesReq
	// query 参数
	{
		val := c.Query("types")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Types = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictService.GetDictsByTypes(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetAllDictsHandler
// @Tags systemRbacdictApi
// @Summary GetAllDictsHandler 获取全部字典-前台/后台使用
// @Description GetAllDictsHandler 获取全部字典-前台/后台使用
// @Success 200 {object} vo.Result{data=_.GetAllDictsRes}
// @Router /api/dicts/all [GET]
func (s *DictApi) GetAllDictsHandler(c *gin.Context) {
	data, err := dictService.GetAllDicts(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// RefreshDictsHandler
// @Tags systemRbacdictApi
// @Summary RefreshDictsHandler 字典缓存刷新-后台使用
// @Description RefreshDictsHandler 字典缓存刷新-后台使用
// @Success 200 {object} vo.Result{data=_.RefreshDictsRes}
// @Router /api/dicts/refresh [POST]
func (s *DictApi) RefreshDictsHandler(c *gin.Context) {
	data, err := dictService.RefreshDicts(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
