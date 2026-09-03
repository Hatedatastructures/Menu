package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type SysUserApi struct{}

// RegisterHandler
// @Tags systemRbacsysUserApi
// @Summary RegisterHandler 用户注册-后台使用
// @Description RegisterHandler 用户注册-后台使用
// @Param data body req.RegisterReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.RegisterRes}
// @Router /user/register [POST]
func (s *SysUserApi) RegisterHandler(c *gin.Context) {
	var req req.RegisterReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysUserService.Register(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// LoginHandler
// @Tags systemRbacsysUserApi
// @Summary LoginHandler 用户登录-后台使用
// @Description LoginHandler 用户登录-后台使用
// @Param data body req.LoginReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.LoginRes}
// @Router /user/login [POST]
func (s *SysUserApi) LoginHandler(c *gin.Context) {
	var req req.LoginReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysUserService.Login(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ChangePasswordHandler
// @Tags systemRbacsysUserApi
// @Summary ChangePasswordHandler 修改用户密码-后台使用
// @Description ChangePasswordHandler 修改用户密码-后台使用
// @Param data body req.ChangePasswordReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/changepassword [PUT]
func (s *SysUserApi) ChangePasswordHandler(c *gin.Context) {
	var req req.ChangePasswordReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysUserService.ChangePassword(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetUserInfoListHandler
// @Tags systemRbacsysUserApi
// @Summary GetUserInfoListHandler 获取用户列表-后台使用
// @Description GetUserInfoListHandler 获取用户列表-后台使用
// @Param data body req.GetUserInfoListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetUserInfoListRes}
// @Router /user/list [GET]
func (s *SysUserApi) GetUserInfoListHandler(c *gin.Context) {
	var req req.GetUserInfoListReq
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
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Size = parsed
		}
	}
	{
		val := c.Query("username")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Username = parsed
		}
	}
	{
		val := c.Query("nickName")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.NickName = parsed
		}
	}
	{
		val := c.Query("phone")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Phone = parsed
		}
	}
	{
		val := c.Query("email")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Email = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysUserService.GetUserInfoList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SetUserAuthorityHandler
// @Tags systemRbacsysUserApi
// @Summary SetUserAuthorityHandler 设置用户角色(单一)-后台使用
// @Description SetUserAuthorityHandler 设置用户角色(单一)-后台使用
// @Param data body req.SetUserAuthorityReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/:id/authority [PUT]
func (s *SysUserApi) SetUserAuthorityHandler(c *gin.Context) {
	var req req.SetUserAuthorityReq
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
	err := sysUserService.SetUserAuthority(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// SetUserAuthoritiesHandler
// @Tags systemRbacsysUserApi
// @Summary SetUserAuthoritiesHandler 设置用户角色(多角色)-后台使用
// @Description SetUserAuthoritiesHandler 设置用户角色(多角色)-后台使用
// @Param data body req.SetUserAuthoritiesReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/authorities [PUT]
func (s *SysUserApi) SetUserAuthoritiesHandler(c *gin.Context) {
	var req req.SetUserAuthoritiesReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysUserService.SetUserAuthorities(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteUserHandler
// @Tags systemRbacsysUserApi
// @Summary DeleteUserHandler 删除用户-后台使用
// @Description DeleteUserHandler 删除用户-后台使用
// @Param data body req.DeleteUserReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/:id [DELETE]
func (s *SysUserApi) DeleteUserHandler(c *gin.Context) {
	var req req.DeleteUserReq
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
	err := sysUserService.DeleteUser(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// SetUserInfoHandler
// @Tags systemRbacsysUserApi
// @Summary SetUserInfoHandler 设置用户信息-后台使用
// @Description SetUserInfoHandler 设置用户信息-后台使用
// @Param data body req.SetUserInfoReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/info [PUT]
func (s *SysUserApi) SetUserInfoHandler(c *gin.Context) {
	var req req.SetUserInfoReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysUserService.SetUserInfo(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// SetSelfInfoHandler
// @Tags systemRbacsysUserApi
// @Summary SetSelfInfoHandler 设置自身信息-后台使用
// @Description SetSelfInfoHandler 设置自身信息-后台使用
// @Param data body req.SetSelfInfoReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/selfinfo [PUT]
func (s *SysUserApi) SetSelfInfoHandler(c *gin.Context) {
	var req req.SetSelfInfoReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysUserService.SetSelfInfo(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// SetSelfSettingHandler
// @Tags systemRbacsysUserApi
// @Summary SetSelfSettingHandler 设置用户配置-后台使用
// @Description SetSelfSettingHandler 设置用户配置-后台使用
// @Param data body req.SetSelfSettingReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /userselfsetting [PUT]
func (s *SysUserApi) SetSelfSettingHandler(c *gin.Context) {
	var req req.SetSelfSettingReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysUserService.SetSelfSetting(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetUserInfHandler
// @Tags systemRbacsysUserApi
// @Summary GetUserInfHandler 获取用户信息-后台使用
// @Description GetUserInfHandler 获取用户信息-后台使用
// @Param data body req.GetUserInfReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetUserInfRes}
// @Router /user/uuid/:uuid [GET]
func (s *SysUserApi) GetUserInfHandler(c *gin.Context) {
	var req req.GetUserInfReq
	// path 参数
	{
		val := c.Param("uuid")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Uuid = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysUserService.GetUserInf(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// FindUserByIdHandler
// @Tags systemRbacsysUserApi
// @Summary FindUserByIdHandler 通过ID获取用户信息-后台使用
// @Description FindUserByIdHandler 通过ID获取用户信息-后台使用
// @Param data body req.FindUserByIdReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.FindUserByIdRes}
// @Router /user/id/:id [GET]
func (s *SysUserApi) FindUserByIdHandler(c *gin.Context) {
	var req req.FindUserByIdReq
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
	data, err := sysUserService.FindUserById(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// FindUserByUuidHandler
// @Tags systemRbacsysUserApi
// @Summary FindUserByUuidHandler 通过UUID获取用户信息-后台使用
// @Description FindUserByUuidHandler 通过UUID获取用户信息-后台使用
// @Param data body req.FindUserByUuidReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.FindUserByUuidRes}
// @Router /user/byuuid/:uuid [GET]
func (s *SysUserApi) FindUserByUuidHandler(c *gin.Context) {
	var req req.FindUserByUuidReq
	// path 参数
	{
		val := c.Param("uuid")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Uuid = parsed
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysUserService.FindUserByUuid(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ResetPasswordHandler
// @Tags systemRbacsysUserApi
// @Summary ResetPasswordHandler 重置用户密码-后台使用
// @Description ResetPasswordHandler 重置用户密码-后台使用
// @Param data body req.ResetPasswordReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /user/resetpassword [PUT]
func (s *SysUserApi) ResetPasswordHandler(c *gin.Context) {
	var req req.ResetPasswordReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := sysUserService.ResetPassword(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}
