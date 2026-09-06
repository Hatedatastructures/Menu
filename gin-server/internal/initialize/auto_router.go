package initialize

import (
	"shack/internal/router"

	"github.com/gin-gonic/gin"
)

func initAutoRouter(PublicGroup, PrivateGroup *gin.RouterGroup, apiSet map[string]struct{}) {

	systemRbacRouter := router.RouterGroupApp.SystemRbac
	systemRbacRouter.InitApiRouter(PublicGroup, apiSet) // 公开

	systemRbacRouter.InitOprationRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigEmailRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigEmailTemplateRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigGroupRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigLoginRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigPasswordRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigRegisterRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigSecurityRouter(PublicGroup, apiSet)
	systemRbacRouter.InitConfigSystemRouter(PublicGroup, apiSet)
	systemRbacRouter.InitSysMonitorRouter(PublicGroup, apiSet)
	systemRbacRouter.InitLoginLogRouter(PublicGroup, apiSet)

	systemRbacRouter.InitDictRouter(PublicGroup, apiSet)
	systemRbacRouter.InitDictDetailRouter(PublicGroup, apiSet)
	systemRbacRouter.InitSysAuthorityRouter(PublicGroup, apiSet)
	systemRbacRouter.InitSysAuthorityBtnRouter(PublicGroup, apiSet)
	systemRbacRouter.InitSysBaseMenuRouter(PublicGroup, apiSet)
	systemRbacRouter.InitSysMenuRouter(PublicGroup, apiSet)
	systemRbacRouter.InitSysUserRouter(PublicGroup, apiSet)
	systemRbacRouter.InitCasbinRouter(PublicGroup, apiSet)
	systemRbacRouter.InitAuthRouter(PublicGroup, apiSet)
	systemRbacRouter.InitFileRouter(PublicGroup, apiSet)
	systemRbacRouter.InitCaptchaClickRouter(PublicGroup, apiSet)
	systemRbacRouter.InitChatRouter(PublicGroup, apiSet)
	systemRbacRouter.InitNoticeRouter(PublicGroup, apiSet)
	systemRbacRouter.InitCaptchaSliderRouter(PublicGroup, apiSet)
	systemRbacRouter.InitCaptchaRotateRouter(PublicGroup, apiSet)
	systemRbacRouter.InitAiConfigRouter(PublicGroup, apiSet)
	systemRbacRouter.InitEmailLimitRouter(PublicGroup, apiSet)
	systemRbacRouter.InitEmailLogRouter(PublicGroup, apiSet)
	systemRbacRouter.InitEmailSendRouter(PublicGroup, apiSet)
	systemRbacRouter.InitEmailStatisticsRouter(PublicGroup, apiSet)
	systemRbacRouter.InitEmailTaskRouter(PublicGroup, apiSet)
	systemRbacRouter.InitEmailTemplateRouter(PublicGroup, apiSet)
	systemRbacRouter.InitVerifyCodeRouter(PublicGroup, apiSet)
	systemRbacRouter.InitProfileRouter(PublicGroup, apiSet)

	systemRbacRouter.InitTemplateRouter(PublicGroup, apiSet)
	systemRbacRouter.InitApiTestRouter(PublicGroup, apiSet)
	genRouter := router.RouterGroupApp.Gen
	genRouter.InitApikeyRouter(PrivateGroup, apiSet)
	genRouter.InitGenerateRouter(PrivateGroup, apiSet)
	genRouter.InitHistoryRouter(PrivateGroup, apiSet)


}
