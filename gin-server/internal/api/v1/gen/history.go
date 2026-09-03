package gen

import (
	req "shack/internal/model/gen/request"
	_ "shack/internal/model/gen/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HistoryApi struct{}

// GetHistoryListHandler
// @Tags genhistoryApi
// @Summary GetHistoryListHandler 获取生成历史列表
// @Description GetHistoryListHandler 获取生成历史列表
// @Param data body req.GetHistoryListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetHistoryListRes}
// @Router /gen/history [GET]
func (s *HistoryApi) GetHistoryListHandler(c *gin.Context) {
	var req req.GetHistoryListReq
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
		val := c.Query("pageSize")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.PageSize = parsed

		}
	}
	{
		val := c.Query("status")
		if val != "" {
			req.Status = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := historyService.GetHistoryList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetHistoryDetailHandler
// @Tags genhistoryApi
// @Summary GetHistoryDetailHandler 获取历史详情(含答案解析)
// @Description GetHistoryDetailHandler 获取历史详情(含答案解析)
// @Param data body req.GetHistoryDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetHistoryDetailRes}
// @Router /gen/history/:id [GET]
func (s *HistoryApi) GetHistoryDetailHandler(c *gin.Context) {
	var req req.GetHistoryDetailReq
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
	data, err := historyService.GetHistoryDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteHistoryHandler
// @Tags genhistoryApi
// @Summary DeleteHistoryHandler 删除历史记录
// @Description DeleteHistoryHandler 删除历史记录
// @Param data body req.DeleteHistoryReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /gen/history/:id [DELETE]
func (s *HistoryApi) DeleteHistoryHandler(c *gin.Context) {
	var req req.DeleteHistoryReq
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
	err := historyService.DeleteHistory(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}
