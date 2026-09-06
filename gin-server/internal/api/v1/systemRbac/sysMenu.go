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

type SysMenuApi struct{}

// GetMenuTreeHandler
// @Tags systemRbacsysMenuApi
// @Summary GetMenuTreeHandler 获取动态菜单树-后台使用
// @Description GetMenuTreeHandler 获取动态菜单树-后台使用
// @Param data body req.GetMenuTreeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetMenuTreeRes}
// @Router /menu/tree/:authorityId [GET]
func (s *SysMenuApi) GetMenuTreeHandler(c *gin.Context) {
	var req req.GetMenuTreeReq
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

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysMenuService.GetMenuTree(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetInfoListHandler
// @Tags systemRbacsysMenuApi
// @Summary GetInfoListHandler 获取基础菜单列表(分页)-后台使用
// @Description GetInfoListHandler 获取基础菜单列表(分页)-后台使用
// @Param data body req.GetInfoListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetInfoListRes}
// @Router /menu/list/:authorityId [GET]
func (s *SysMenuApi) GetInfoListHandler(c *gin.Context) {
	var req req.GetInfoListReq
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

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysMenuService.GetInfoList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// AddBaseMenuHandler
// @Tags systemRbacsysMenuApi
// @Summary AddBaseMenuHandler 添加基础菜单-后台使用
// @Description AddBaseMenuHandler 添加基础菜单-后台使用
// @Param data body req.AddBaseMenuReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.AddBaseMenuRes}
// @Router /menu/create [POST]
func (s *SysMenuApi) AddBaseMenuHandler(c *gin.Context) {
	var req req.AddBaseMenuReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysMenuService.AddBaseMenu(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetBaseMenuTreeHandler
// @Tags systemRbacsysMenuApi
// @Summary GetBaseMenuTreeHandler 获取基础菜单树-后台使用
// @Description GetBaseMenuTreeHandler 获取基础菜单树-后台使用
// @Param data body req.GetBaseMenuTreeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetBaseMenuTreeRes}
// @Router /menu/base/tree/:authorityId [GET]
func (s *SysMenuApi) GetBaseMenuTreeHandler(c *gin.Context) {
	var req req.GetBaseMenuTreeReq
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

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysMenuService.GetBaseMenuTree(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// AddMenuAuthorityHandler
// @Tags systemRbacsysMenuApi
// @Summary AddMenuAuthorityHandler 为角色分配菜单-后台使用
// @Description AddMenuAuthorityHandler 为角色分配菜单-后台使用
// @Param data body req.AddMenuAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /menu/authority [POST]
func (s *SysMenuApi) AddMenuAuthorityHandler(c *gin.Context) {
	var req req.AddMenuAuthorityReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysMenuService.AddMenuAuthority(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetMenuAuthorityHandler
// @Tags systemRbacsysMenuApi
// @Summary GetMenuAuthorityHandler 获取角色菜单权限-后台使用
// @Description GetMenuAuthorityHandler 获取角色菜单权限-后台使用
// @Param data body req.GetMenuAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetMenuAuthorityRes}
// @Router /menu/authority/:authorityId [GET]
func (s *SysMenuApi) GetMenuAuthorityHandler(c *gin.Context) {
	var req req.GetMenuAuthorityReq
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

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysMenuService.GetMenuAuthority(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}