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

type DictDetailApi struct{}

// GetDictionaryDetailListHandler
// @Tags systemRbacdictDetailApi
// @Summary GetDictionaryDetailListHandler 获取字典详情列表-后台使用
// @Description GetDictionaryDetailListHandler 获取字典详情列表-后台使用
// @Param data body req.GetDictionaryDetailListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryDetailListRes}
// @Router /api/dictionary-details [GET]
func (s *DictDetailApi) GetDictionaryDetailListHandler(c *gin.Context) {
	var req req.GetDictionaryDetailListReq
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
		val := c.Query("label")
		if val != "" {
			req.Label = val
		}
	}
	{
		val := c.Query("value")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Value = parsed
		}
	}
	{
		val := c.Query("status")
		if val != "" {
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Status = parsed
		}
	}
	{
		val := c.Query("sysDictionaryID")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.SysDictionaryID = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictDetailService.GetDictionaryDetailList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CreateDictionaryDetailHandler
// @Tags systemRbacdictDetailApi
// @Summary CreateDictionaryDetailHandler 创建字典详情-后台使用
// @Description CreateDictionaryDetailHandler 创建字典详情-后台使用
// @Param data body req.CreateDictionaryDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateDictionaryDetailRes}
// @Router /api/dictionary-details [POST]
func (s *DictDetailApi) CreateDictionaryDetailHandler(c *gin.Context) {
	var req req.CreateDictionaryDetailReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictDetailService.CreateDictionaryDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDictionaryDetailHandler
// @Tags systemRbacdictDetailApi
// @Summary GetDictionaryDetailHandler 字典详情详情-后台使用
// @Description GetDictionaryDetailHandler 字典详情详情-后台使用
// @Param data body req.GetDictionaryDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryDetailRes}
// @Router /api/dictionary-details/:id [GET]
func (s *DictDetailApi) GetDictionaryDetailHandler(c *gin.Context) {
	var req req.GetDictionaryDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
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
	data, err := dictDetailService.GetDictionaryDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateDictionaryDetailHandler
// @Tags systemRbacdictDetailApi
// @Summary UpdateDictionaryDetailHandler 更新字典详情-后台使用
// @Description UpdateDictionaryDetailHandler 更新字典详情-后台使用
// @Param data body req.UpdateDictionaryDetailReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/dictionary-details/:id [PUT]
func (s *DictDetailApi) UpdateDictionaryDetailHandler(c *gin.Context) {
	var req req.UpdateDictionaryDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
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
	err := dictDetailService.UpdateDictionaryDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteDictionaryDetailHandler
// @Tags systemRbacdictDetailApi
// @Summary DeleteDictionaryDetailHandler 删除字典详情-后台使用
// @Description DeleteDictionaryDetailHandler 删除字典详情-后台使用
// @Param data body req.DeleteDictionaryDetailReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/dictionary-details/:id [DELETE]
func (s *DictDetailApi) DeleteDictionaryDetailHandler(c *gin.Context) {
	var req req.DeleteDictionaryDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
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
	err := dictDetailService.DeleteDictionaryDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetDictionaryListByIdHandler
// @Tags systemRbacdictDetailApi
// @Summary GetDictionaryListByIdHandler 按字典ID获取字典全部内容-后台使用
// @Description GetDictionaryListByIdHandler 按字典ID获取字典全部内容-后台使用
// @Param data body req.GetDictionaryListByIdReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryListByIdRes}
// @Router /api/dictionary-details/dictionary/:dictionaryId [GET]
func (s *DictDetailApi) GetDictionaryListByIdHandler(c *gin.Context) {
	var req req.GetDictionaryListByIdReq
	// path 参数
	{
		val := c.Param("dictionaryId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.DictionaryId = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictDetailService.GetDictionaryListById(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDictionaryListByTypeHandler
// @Tags systemRbacdictDetailApi
// @Summary GetDictionaryListByTypeHandler 按字典type获取字典全部内容-后台使用
// @Description GetDictionaryListByTypeHandler 按字典type获取字典全部内容-后台使用
// @Param data body req.GetDictionaryListByTypeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryListByTypeRes}
// @Router /api/dictionary-details/type/:type [GET]
func (s *DictDetailApi) GetDictionaryListByTypeHandler(c *gin.Context) {
	var req req.GetDictionaryListByTypeReq
	// path 参数
	{
		val := c.Param("type")
		req.Type = val
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictDetailService.GetDictionaryListByType(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDictionaryInfoByValueHandler
// @Tags systemRbacdictDetailApi
// @Summary GetDictionaryInfoByValueHandler 按字典ID和value获取单条字典内容-后台使用
// @Description GetDictionaryInfoByValueHandler 按字典ID和value获取单条字典内容-后台使用
// @Param data body req.GetDictionaryInfoByValueReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryInfoByValueRes}
// @Router /api/dictionary-details/:dictionaryId/value/:value [GET]
func (s *DictDetailApi) GetDictionaryInfoByValueHandler(c *gin.Context) {
	var req req.GetDictionaryInfoByValueReq
	// path 参数
	{
		val := c.Param("dictionaryId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.DictionaryId = parsed
	}
	{
		val := c.Param("value")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Value = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictDetailService.GetDictionaryInfoByValue(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDictionaryInfoByTypeValueHandler
// @Tags systemRbacdictDetailApi
// @Summary GetDictionaryInfoByTypeValueHandler 按字典type和value获取单条字典内容-后台使用
// @Description GetDictionaryInfoByTypeValueHandler 按字典type和value获取单条字典内容-后台使用
// @Param data body req.GetDictionaryInfoByTypeValueReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetDictionaryInfoByTypeValueRes}
// @Router /api/dictionary-details/type/:type/value/:value [GET]
func (s *DictDetailApi) GetDictionaryInfoByTypeValueHandler(c *gin.Context) {
	var req req.GetDictionaryInfoByTypeValueReq
	// path 参数
	{
		val := c.Param("type")
		req.Type = val
	}
	{
		val := c.Param("value")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Value = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := dictDetailService.GetDictionaryInfoByTypeValue(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}