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

type CasbinApi struct{}

// UpdateCasbinHandler
// @Tags systemRbaccasbinApi
// @Summary UpdateCasbinHandler 更新权限-后台使用
// @Description UpdateCasbinHandler 更新权限-后台使用
// @Param data body req.UpdateCasbinReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/casbin [POST]
func (s *CasbinApi) UpdateCasbinHandler(c *gin.Context) {
	var req req.UpdateCasbinReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := casbinService.UpdateCasbin(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetPolicyPathByAuthorityIdHandler
// @Tags systemRbaccasbinApi
// @Summary GetPolicyPathByAuthorityIdHandler 获取权限列表-后台使用
// @Description GetPolicyPathByAuthorityIdHandler 获取权限列表-后台使用
// @Param data body req.GetPolicyPathByAuthorityIdReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetPolicyPathByAuthorityIdRes}
// @Router /api/casbin/:authorityId [GET]
func (s *CasbinApi) GetPolicyPathByAuthorityIdHandler(c *gin.Context) {
	var req req.GetPolicyPathByAuthorityIdReq
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
	data, err := casbinService.GetPolicyPathByAuthorityId(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ClearCasbinHandler
// @Tags systemRbaccasbinApi
// @Summary ClearCasbinHandler 清除权限-后台使用
// @Description ClearCasbinHandler 清除权限-后台使用
// @Param data body req.ClearCasbinReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/casbin [DELETE]
func (s *CasbinApi) ClearCasbinHandler(c *gin.Context) {
	var req req.ClearCasbinReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := casbinService.ClearCasbin(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// FreshCasbinHandler
// @Tags systemRbaccasbinApi
// @Summary FreshCasbinHandler 刷新Casbin-后台使用
// @Description FreshCasbinHandler 刷新Casbin-后台使用
// @Success 200 {object} vo.Result{}
// @Router /api/casbin/refresh [POST]
func (s *CasbinApi) FreshCasbinHandler(c *gin.Context) {
	err := casbinService.FreshCasbin(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(2000, "操作失败")))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}