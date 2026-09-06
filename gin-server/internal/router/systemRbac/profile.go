//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type ProfileRouter struct{}

func (s *ProfileRouter) InitProfileRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	profileRouter := Router.Group("")
	utils.RegisterApi(profileRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/profile", "获取个人信息-前台使用", profileApi.GetProfileHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/profile", "更新个人信息-前台使用", profileApi.UpdateProfileHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/profile/password", "修改密码-前台使用", profileApi.UpdatePasswordHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/profile/avatar", "修改头像-前台使用", profileApi.UpdateAvatarHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/profile/avatar/upload", "上传头像-前台使用", profileApi.UploadAvatarHandler),
	)
	return profileRouter
}
