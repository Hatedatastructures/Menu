package systemRbac

import (
	"shack/internal/service"
	"shack/internal/service/systemRbac"
)

type ApiGroup struct {
	AiConfigApi
	ApiApi
	ApiTestApi
	AuthApi
	CaptchaClickApi
	CaptchaRotateApi
	CaptchaSliderApi
	CasbinApi
	ChatApi
	ConfigEmailApi
	ConfigEmailTemplateApi
	ConfigEncryptApi
	ConfigGroupApi
	ConfigLoginApi
	ConfigPasswordApi
	ConfigRegisterApi
	ConfigSecurityApi
	ConfigSystemApi
	DictApi
	DictDetailApi
	EmailLimitApi
	EmailLogApi
	EmailSendApi
	EmailStatisticsApi
	EmailTaskApi
	EmailTemplateApi
	FileApi
	LoginLogApi
	NoticeApi
	OprationApi
	ProfileApi
	SysAuthorityApi
	SysAuthorityBtnApi
	SysBaseMenuApi
	SysMenuApi
	SysMonitorApi
	SysUserApi
	TemplateApi
	VerifyCodeApi
}

var (
	aiConfigService            = service.ServiceGroupApp.SystemRbacServiceGroup.AiConfigService
	apiService                 = service.ServiceGroupApp.SystemRbacServiceGroup.ApiService
	apiTestService             = service.ServiceGroupApp.SystemRbacServiceGroup.ApiTestService
	authService                = service.ServiceGroupApp.SystemRbacServiceGroup.AuthService
	captchaClickService        = service.ServiceGroupApp.SystemRbacServiceGroup.CaptchaClickService
	captchaRotateService       = service.ServiceGroupApp.SystemRbacServiceGroup.CaptchaRotateService
	captchaSliderService       = service.ServiceGroupApp.SystemRbacServiceGroup.CaptchaSliderService
	casbinService              = service.ServiceGroupApp.SystemRbacServiceGroup.CasbinService
	chatService                = service.ServiceGroupApp.SystemRbacServiceGroup.ChatService
	configEmailService         = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigEmailService
	configEmailTemplateService = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigEmailTemplateService
	configEncryptService       = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigEncryptService
	configGroupService         = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigGroupService
	configLoginService         = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigLoginService
	configPasswordService      = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigPasswordService
	configRegisterService      = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigRegisterService
	configSecurityService      = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigSecurityService
	configSystemService        = service.ServiceGroupApp.SystemRbacServiceGroup.ConfigSystemService
	dictDetailService          = service.ServiceGroupApp.SystemRbacServiceGroup.DictDetailService
	dictService                = service.ServiceGroupApp.SystemRbacServiceGroup.DictService
	emailLimitService          = service.ServiceGroupApp.SystemRbacServiceGroup.EmailLimitService
	emailLogService            = service.ServiceGroupApp.SystemRbacServiceGroup.EmailLogService
	emailSendService           = service.ServiceGroupApp.SystemRbacServiceGroup.EmailSendService
	emailStatisticsService     = service.ServiceGroupApp.SystemRbacServiceGroup.EmailStatisticsService
	emailTaskService           = service.ServiceGroupApp.SystemRbacServiceGroup.EmailTaskService
	emailTemplateService       = service.ServiceGroupApp.SystemRbacServiceGroup.EmailTemplateService
	fileService                = service.ServiceGroupApp.SystemRbacServiceGroup.FileService
	loginLogService            = service.ServiceGroupApp.SystemRbacServiceGroup.LoginLogService
	noticeService              = service.ServiceGroupApp.SystemRbacServiceGroup.NoticeService
	oprationService            = service.ServiceGroupApp.SystemRbacServiceGroup.OprationService
	profileService             = service.ServiceGroupApp.SystemRbacServiceGroup.ProfileService
	sysAuthorityBtnService     = service.ServiceGroupApp.SystemRbacServiceGroup.SysAuthorityBtnService
	sysAuthorityService        = service.ServiceGroupApp.SystemRbacServiceGroup.SysAuthorityService
	sysBaseMenuService         = service.ServiceGroupApp.SystemRbacServiceGroup.SysBaseMenuService
	sysMenuService             = service.ServiceGroupApp.SystemRbacServiceGroup.SysMenuService
	sysMonitorService          = systemRbac.GetSysMonitorService()
	sysUserService             = service.ServiceGroupApp.SystemRbacServiceGroup.SysUserService
	templateService            = service.ServiceGroupApp.SystemRbacServiceGroup.TemplateService
	verifyCodeService          = service.ServiceGroupApp.SystemRbacServiceGroup.VerifyCodeService
)
