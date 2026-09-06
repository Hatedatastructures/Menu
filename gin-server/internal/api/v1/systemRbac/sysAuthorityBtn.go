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

type SysAuthorityBtnApi struct{}

// GetAuthorityBtnHandler
// @Tags systemRbacsysAuthorityBtnApi
// @Summary GetAuthorityBtnHandler 获取角色按钮权限-后台使用
// @Description GetAuthorityBtnHandler 获取角色按钮权限-后台使用
// @Param data body req.GetAuthorityBtnReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetAuthorityBtnRes}
// @Router /authoritybtn/:authorityId/:menuID [GET]
func (s *SysAuthorityBtnApi) GetAuthorityBtnHandler(c *gin.Context) {
	var req req.GetAuthorityBtnReq
	// path 参数
	{
		val := c.Param("authorityId")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.AuthorityId = uint(parsed)

	}
	{
		val := c.Param("menuID")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.MenuID = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityBtnService.GetAuthorityBtn(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SetAuthorityBtnHandler
// @Tags systemRbacsysAuthorityBtnApi
// @Summary SetAuthorityBtnHandler 设置角色按钮权限-后台使用
// @Description SetAuthorityBtnHandler 设置角色按钮权限-后台使用
// @Param data body req.SetAuthorityBtnReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /authoritybtn [PUT]
func (s *SysAuthorityBtnApi) SetAuthorityBtnHandler(c *gin.Context) {
	var req req.SetAuthorityBtnReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysAuthorityBtnService.SetAuthorityBtn(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// CanRemoveAuthorityBtnHandler
// @Tags systemRbacsysAuthorityBtnApi
// @Summary CanRemoveAuthorityBtnHandler 检查角色按钮是否可以删除-后台使用
// @Description CanRemoveAuthorityBtnHandler 检查角色按钮是否可以删除-后台使用
// @Param data body req.CanRemoveAuthorityBtnReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /authoritybtn/check/:id [DELETE]
func (s *SysAuthorityBtnApi) CanRemoveAuthorityBtnHandler(c *gin.Context) {
	var req req.CanRemoveAuthorityBtnReq
	// path 参数
	{
		val := c.Param("id")
		req.Id = val

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysAuthorityBtnService.CanRemoveAuthorityBtn(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}