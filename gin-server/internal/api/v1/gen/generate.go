package gen

import (
	req "shack/internal/model/gen/request"
	_ "shack/internal/model/gen/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GenerateApi struct{}

// GenerateHandler
// @Tags gengenerateApi
// @Summary GenerateHandler 上传PDF并生成答案解析(异步)
// @Description GenerateHandler 上传PDF并生成答案解析(异步)
// @Param data body req.GenerateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GenerateRes}
// @Router /gen/generate [POST]
func (s *GenerateApi) GenerateHandler(c *gin.Context) {
	var req req.GenerateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := generateService.Generate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetGenerateResultHandler
// @Tags gengenerateApi
// @Summary GetGenerateResultHandler 获取生成进度/结果
// @Description GetGenerateResultHandler 获取生成进度/结果
// @Param data body req.GetGenerateResultReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetGenerateResultRes}
// @Router /gen/generate/:id [GET]
func (s *GenerateApi) GetGenerateResultHandler(c *gin.Context) {
	var req req.GetGenerateResultReq
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
	data, err := generateService.GetGenerateResult(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}