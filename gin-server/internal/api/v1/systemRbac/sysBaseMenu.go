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

type SysBaseMenuApi struct{}

// DeleteBaseMenuHandler
// @Tags systemRbacsysBaseMenuApi
// @Summary DeleteBaseMenuHandler 删除基础菜单-后台使用
// @Description DeleteBaseMenuHandler 删除基础菜单-后台使用
// @Param data body req.DeleteBaseMenuReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /basemenu/:id [DELETE]
func (s *SysBaseMenuApi) DeleteBaseMenuHandler(c *gin.Context) {
	var req req.DeleteBaseMenuReq
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
	err := sysBaseMenuService.DeleteBaseMenu(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UpdateBaseMenuHandler
// @Tags systemRbacsysBaseMenuApi
// @Summary UpdateBaseMenuHandler 更新基础菜单-后台使用
// @Description UpdateBaseMenuHandler 更新基础菜单-后台使用
// @Param data body req.UpdateBaseMenuReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateBaseMenuRes}
// @Router /basemenu/update [PUT]
func (s *SysBaseMenuApi) UpdateBaseMenuHandler(c *gin.Context) {
	var req req.UpdateBaseMenuReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysBaseMenuService.UpdateBaseMenu(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetBaseMenuByIdHandler
// @Tags systemRbacsysBaseMenuApi
// @Summary GetBaseMenuByIdHandler 获取基础菜单详情-后台使用
// @Description GetBaseMenuByIdHandler 获取基础菜单详情-后台使用
// @Param data body req.GetBaseMenuByIdReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetBaseMenuByIdRes}
// @Router /basemenu/:id [GET]
func (s *SysBaseMenuApi) GetBaseMenuByIdHandler(c *gin.Context) {
	var req req.GetBaseMenuByIdReq
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
	data, err := sysBaseMenuService.GetBaseMenuById(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}