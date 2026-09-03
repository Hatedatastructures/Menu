package systemRbac

import api "shack/internal/api/v1"

type RouterGroup struct {
	AiConfigRouter
	ApiRouter
	ApiTestRouter
	AuthRouter
	CaptchaClickRouter
	CaptchaRotateRouter
	CaptchaSliderRouter
	CasbinRouter
	ChatRouter
	ConfigEmailRouter
	ConfigEmailTemplateRouter
	ConfigGroupRouter
	ConfigLoginRouter
	ConfigPasswordRouter
	ConfigRegisterRouter
	ConfigSecurityRouter
	ConfigSystemRouter
	DictDetailRouter
	DictRouter
	EmailLimitRouter
	EmailLogRouter
	EmailSendRouter
	EmailStatisticsRouter
	EmailTaskRouter
	EmailTemplateRouter
	FileRouter
	LoginLogRouter
	NoticeRouter
	OprationRouter
	ProfileRouter
	SysAuthorityBtnRouter
	SysAuthorityRouter
	SysBaseMenuRouter
	SysMenuRouter
	SysMonitorRouter
	SysUserRouter
	TemplateRouter
	VerifyCodeRouter
}

var (
	aiConfigApi = api.ApiGroupApp.SystemRbacApiGroup.AiConfigApi
	apiApi = api.ApiGroupApp.SystemRbacApiGroup.ApiApi
	apiTestApi = api.ApiGroupApp.SystemRbacApiGroup.ApiTestApi
	authApi = api.ApiGroupApp.SystemRbacApiGroup.AuthApi
	captchaClickApi = api.ApiGroupApp.SystemRbacApiGroup.CaptchaClickApi
	captchaRotateApi = api.ApiGroupApp.SystemRbacApiGroup.CaptchaRotateApi
	captchaSliderApi = api.ApiGroupApp.SystemRbacApiGroup.CaptchaSliderApi
	casbinApi = api.ApiGroupApp.SystemRbacApiGroup.CasbinApi
	chatApi = api.ApiGroupApp.SystemRbacApiGroup.ChatApi
	configEmailApi         = api.ApiGroupApp.SystemRbacApiGroup.ConfigEmailApi
	configEmailTemplateApi = api.ApiGroupApp.SystemRbacApiGroup.ConfigEmailTemplateApi
	configGroupApi         = api.ApiGroupApp.SystemRbacApiGroup.ConfigGroupApi
	configLoginApi         = api.ApiGroupApp.SystemRbacApiGroup.ConfigLoginApi
	configPasswordApi      = api.ApiGroupApp.SystemRbacApiGroup.ConfigPasswordApi
	configRegisterApi      = api.ApiGroupApp.SystemRbacApiGroup.ConfigRegisterApi
	configSecurityApi      = api.ApiGroupApp.SystemRbacApiGroup.ConfigSecurityApi
	configSystemApi        = api.ApiGroupApp.SystemRbacApiGroup.ConfigSystemApi
	dictApi                = api.ApiGroupApp.SystemRbacApiGroup.DictApi
	dictDetailApi          = api.ApiGroupApp.SystemRbacApiGroup.DictDetailApi
	emailLimitApi = api.ApiGroupApp.SystemRbacApiGroup.EmailLimitApi
	emailLogApi = api.ApiGroupApp.SystemRbacApiGroup.EmailLogApi
	emailSendApi = api.ApiGroupApp.SystemRbacApiGroup.EmailSendApi
	emailStatisticsApi = api.ApiGroupApp.SystemRbacApiGroup.EmailStatisticsApi
	emailTaskApi = api.ApiGroupApp.SystemRbacApiGroup.EmailTaskApi
	emailTemplateApi = api.ApiGroupApp.SystemRbacApiGroup.EmailTemplateApi
	fileApi = api.ApiGroupApp.SystemRbacApiGroup.FileApi
	loginLogApi            = api.ApiGroupApp.SystemRbacApiGroup.LoginLogApi
	noticeApi = api.ApiGroupApp.SystemRbacApiGroup.NoticeApi
	oprationApi = api.ApiGroupApp.SystemRbacApiGroup.OprationApi
	profileApi = api.ApiGroupApp.SystemRbacApiGroup.ProfileApi
	sysAuthorityApi    = api.ApiGroupApp.SystemRbacApiGroup.SysAuthorityApi
	sysAuthorityBtnApi = api.ApiGroupApp.SystemRbacApiGroup.SysAuthorityBtnApi
	sysBaseMenuApi     = api.ApiGroupApp.SystemRbacApiGroup.SysBaseMenuApi
	sysMenuApi         = api.ApiGroupApp.SystemRbacApiGroup.SysMenuApi
	sysMonitorApi      = api.ApiGroupApp.SystemRbacApiGroup.SysMonitorApi
	sysUserApi         = api.ApiGroupApp.SystemRbacApiGroup.SysUserApi
	templateApi = api.ApiGroupApp.SystemRbacApiGroup.TemplateApi
	verifyCodeApi = api.ApiGroupApp.SystemRbacApiGroup.VerifyCodeApi
)
