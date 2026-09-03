package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type EmailTaskApi struct{}

// CreateEmailTaskHandler
// @Tags systemRbacemailTaskApi
// @Summary CreateEmailTaskHandler 创建邮件任务-后台使用
// @Description CreateEmailTaskHandler 创建邮件任务-后台使用
// @Param data body req.CreateEmailTaskReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/task [POST]
func (s *EmailTaskApi) CreateEmailTaskHandler(c *gin.Context) {
	var req req.CreateEmailTaskReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := emailTaskService.CreateEmailTask(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UpdateEmailTaskHandler
// @Tags systemRbacemailTaskApi
// @Summary UpdateEmailTaskHandler 更新邮件任务-后台使用
// @Description UpdateEmailTaskHandler 更新邮件任务-后台使用
// @Param data body req.UpdateEmailTaskReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/task/:id [PUT]
func (s *EmailTaskApi) UpdateEmailTaskHandler(c *gin.Context) {
	var req req.UpdateEmailTaskReq
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
	err := emailTaskService.UpdateEmailTask(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteEmailTaskHandler
// @Tags systemRbacemailTaskApi
// @Summary DeleteEmailTaskHandler 删除邮件任务-后台使用
// @Description DeleteEmailTaskHandler 删除邮件任务-后台使用
// @Param data body req.DeleteEmailTaskReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/task/:id [DELETE]
func (s *EmailTaskApi) DeleteEmailTaskHandler(c *gin.Context) {
	var req req.DeleteEmailTaskReq
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
	err := emailTaskService.DeleteEmailTask(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetEmailTaskListHandler
// @Tags systemRbacemailTaskApi
// @Summary GetEmailTaskListHandler 获取邮件任务列表-后台使用
// @Description GetEmailTaskListHandler 获取邮件任务列表-后台使用
// @Param data body req.GetEmailTaskListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailTaskListRes}
// @Router /email/task/list [GET]
func (s *EmailTaskApi) GetEmailTaskListHandler(c *gin.Context) {
	var req req.GetEmailTaskListReq
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
	data, err := emailTaskService.GetEmailTaskList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// StartEmailTaskHandler
// @Tags systemRbacemailTaskApi
// @Summary StartEmailTaskHandler 启动邮件任务-后台使用
// @Description StartEmailTaskHandler 启动邮件任务-后台使用
// @Param data body req.StartEmailTaskReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/task/start/:id [POST]
func (s *EmailTaskApi) StartEmailTaskHandler(c *gin.Context) {
	var req req.StartEmailTaskReq
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
	err := emailTaskService.StartEmailTask(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// StopEmailTaskHandler
// @Tags systemRbacemailTaskApi
// @Summary StopEmailTaskHandler 停止邮件任务-后台使用
// @Description StopEmailTaskHandler 停止邮件任务-后台使用
// @Param data body req.StopEmailTaskReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/task/stop/:id [POST]
func (s *EmailTaskApi) StopEmailTaskHandler(c *gin.Context) {
	var req req.StopEmailTaskReq
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
	err := emailTaskService.StopEmailTask(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// RunEmailTaskHandler
// @Tags systemRbacemailTaskApi
// @Summary RunEmailTaskHandler 手动执行邮件任务-后台使用
// @Description RunEmailTaskHandler 手动执行邮件任务-后台使用
// @Param data body req.RunEmailTaskReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.RunEmailTaskRes}
// @Router /email/task/run/:id [POST]
func (s *EmailTaskApi) RunEmailTaskHandler(c *gin.Context) {
	var req req.RunEmailTaskReq
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
	data, err := emailTaskService.RunEmailTask(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}