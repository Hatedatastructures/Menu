package systemRbac

import (
	"strconv"

	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type SysAuthorityApi struct{}

// CreateAuthorityHandler
// @Tags systemRbacsysAuthorityApi
// @Summary CreateAuthorityHandler 创建角色-后台使用
// @Description CreateAuthorityHandler 创建角色-后台使用
// @Param data body req.CreateAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateAuthorityRes}
// @Router /authority/create [POST]
func (s *SysAuthorityApi) CreateAuthorityHandler(c *gin.Context) {
	var req req.CreateAuthorityReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.CreateAuthority(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CopyAuthorityHandler
// @Tags systemRbacsysAuthorityApi
// @Summary CopyAuthorityHandler 复制角色-后台使用
// @Description CopyAuthorityHandler 复制角色-后台使用
// @Param data body req.CopyAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CopyAuthorityRes}
// @Router /authority/copy [POST]
func (s *SysAuthorityApi) CopyAuthorityHandler(c *gin.Context) {
	var req req.CopyAuthorityReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.CopyAuthority(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateAuthorityHandler
// @Tags systemRbacsysAuthorityApi
// @Summary UpdateAuthorityHandler 更新角色-后台使用
// @Description UpdateAuthorityHandler 更新角色-后台使用
// @Param data body req.UpdateAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateAuthorityRes}
// @Router /authority/update [PUT]
func (s *SysAuthorityApi) UpdateAuthorityHandler(c *gin.Context) {
	var req req.UpdateAuthorityReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.UpdateAuthority(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteAuthorityHandler
// @Tags systemRbacsysAuthorityApi
// @Summary DeleteAuthorityHandler 删除角色-后台使用
// @Description DeleteAuthorityHandler 删除角色-后台使用
// @Param data body req.DeleteAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /authority/:authorityId [DELETE]
func (s *SysAuthorityApi) DeleteAuthorityHandler(c *gin.Context) {
	var req req.DeleteAuthorityReq
	// path 参数
	{
		val := c.Param("authorityId")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
			return
		}
		req.AuthorityId = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	err := sysAuthorityService.DeleteAuthority(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetAuthorityInfoListHandler
// @Tags systemRbacsysAuthorityApi
// @Summary GetAuthorityInfoListHandler 获取角色列表-后台使用
// @Description GetAuthorityInfoListHandler 获取角色列表-后台使用
// @Param data body req.GetAuthorityInfoListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetAuthorityInfoListRes}
// @Router /authority/list [GET]
func (s *SysAuthorityApi) GetAuthorityInfoListHandler(c *gin.Context) {
	var req req.GetAuthorityInfoListReq
	// query 参数
	{
		val := c.Query("authorityId")
		if val != "" {
			parsed, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
				return
			}
			req.AuthorityId = uint(parsed)

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.GetAuthorityInfoList(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetStructAuthorityListHandler
// @Tags systemRbacsysAuthorityApi
// @Summary GetStructAuthorityListHandler 获取角色结构列表-后台使用
// @Description GetStructAuthorityListHandler 获取角色结构列表-后台使用
// @Param data body req.GetStructAuthorityListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetStructAuthorityListRes}
// @Router /authority/structlist/:authorityId [GET]
func (s *SysAuthorityApi) GetStructAuthorityListHandler(c *gin.Context) {
	var req req.GetStructAuthorityListReq
	// path 参数
	{
		val := c.Param("authorityId")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
			return
		}
		req.AuthorityId = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.GetStructAuthorityList(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetAuthorityInfoHandler
// @Tags systemRbacsysAuthorityApi
// @Summary GetAuthorityInfoHandler 获取角色信息-后台使用
// @Description GetAuthorityInfoHandler 获取角色信息-后台使用
// @Param data body req.GetAuthorityInfoReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetAuthorityInfoRes}
// @Router /authority/info/:authorityId [GET]
func (s *SysAuthorityApi) GetAuthorityInfoHandler(c *gin.Context) {
	var req req.GetAuthorityInfoReq
	// path 参数
	{
		val := c.Param("authorityId")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
			return
		}
		req.AuthorityId = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.GetAuthorityInfo(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SetDataAuthorityHandler
// @Tags systemRbacsysAuthorityApi
// @Summary SetDataAuthorityHandler 设置角色数据权限-后台使用
// @Description SetDataAuthorityHandler 设置角色数据权限-后台使用
// @Param data body req.SetDataAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /authority/data [PUT]
func (s *SysAuthorityApi) SetDataAuthorityHandler(c *gin.Context) {
	var req req.SetDataAuthorityReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	err := sysAuthorityService.SetDataAuthority(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// SetMenuAuthorityHandler
// @Tags systemRbacsysAuthorityApi
// @Summary SetMenuAuthorityHandler 设置角色菜单权限-后台使用
// @Description SetMenuAuthorityHandler 设置角色菜单权限-后台使用
// @Param data body req.SetMenuAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /authority/menu [PUT]
func (s *SysAuthorityApi) SetMenuAuthorityHandler(c *gin.Context) {
	var req req.SetMenuAuthorityReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	err := sysAuthorityService.SetMenuAuthority(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetParentAuthorityIDHandler
// @Tags systemRbacsysAuthorityApi
// @Summary GetParentAuthorityIDHandler 获取父角色ID-后台使用
// @Description GetParentAuthorityIDHandler 获取父角色ID-后台使用
// @Param data body req.GetParentAuthorityIDReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetParentAuthorityIDRes}
// @Router /authority/parent/:authorityId [GET]
func (s *SysAuthorityApi) GetParentAuthorityIDHandler(c *gin.Context) {
	var req req.GetParentAuthorityIDReq
	// path 参数
	{
		val := c.Param("authorityId")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(400, vo.Fail(c, "", biz_err.New(biz_err.PARAM_ERROR, "参数错误")))
			return
		}
		req.AuthorityId = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	data, err := sysAuthorityService.GetParentAuthorityID(c, req)
	if err != nil {
		c.JSON(400, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
