package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ProfileApi struct{}

// GetProfileHandler
// @Tags systemRbacprofileApi
// @Summary GetProfileHandler 获取个人信息-前台使用
// @Description GetProfileHandler 获取个人信息-前台使用
// @Success 200 {object} vo.Result{data=_.GetProfileRes}
// @Router /api/profile [GET]
func (s *ProfileApi) GetProfileHandler(c *gin.Context) {
	data, err := profileService.GetProfile(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateProfileHandler
// @Tags systemRbacprofileApi
// @Summary UpdateProfileHandler 更新个人信息-前台使用
// @Description UpdateProfileHandler 更新个人信息-前台使用
// @Param data body req.UpdateProfileReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateProfileRes}
// @Router /api/profile [PUT]
func (s *ProfileApi) UpdateProfileHandler(c *gin.Context) {
	var req req.UpdateProfileReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := profileService.UpdateProfile(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdatePasswordHandler
// @Tags systemRbacprofileApi
// @Summary UpdatePasswordHandler 修改密码-前台使用
// @Description UpdatePasswordHandler 修改密码-前台使用
// @Param data body req.UpdatePasswordReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/profile/password [PUT]
func (s *ProfileApi) UpdatePasswordHandler(c *gin.Context) {
	var req req.UpdatePassword1Req
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := profileService.UpdatePassword(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UpdateAvatarHandler
// @Tags systemRbacprofileApi
// @Summary UpdateAvatarHandler 修改头像-前台使用
// @Description UpdateAvatarHandler 修改头像-前台使用
// @Param data body req.UpdateAvatarReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateAvatarRes}
// @Router /api/profile/avatar [PUT]
func (s *ProfileApi) UpdateAvatarHandler(c *gin.Context) {
	var req req.UpdateAvatarReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := profileService.UpdateAvatar(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UploadAvatarHandler
// @Tags systemRbacprofileApi
// @Summary UploadAvatarHandler 上传头像-前台使用
// @Description UploadAvatarHandler 上传头像-前台使用
// @Param data body req.UploadAvatarReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UploadAvatarRes}
// @Router /api/profile/avatar/upload [POST]
func (s *ProfileApi) UploadAvatarHandler(c *gin.Context) {
	var req req.UploadAvatarReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := profileService.UploadAvatar(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
