package system

import api "shack/internal/api/v1"

type RouterGroup struct {
	ApiRouter
	JwtRouter
	SysRouter
	BaseRouter
	InitRouter
	MenuRouter
	UserRouter
	CasbinRouter
	AuthorityRouter
	DictionaryRouter
	OperationRecordRouter
	DictionaryDetailRouter
	AuthorityBtnRouter
	SysExportTemplateRouter
	SysParamsRouter
	SysVersionRouter
}

var (
	dbApi               = api.ApiGroupApp.SystemApiGroup.DBApi
	jwtApi              = api.ApiGroupApp.SystemApiGroup.JwtApi
	baseApi             = api.ApiGroupApp.SystemApiGroup.BaseApi
	casbinApi           = api.ApiGroupApp.SystemApiGroup.CasbinApi
	systemApi           = api.ApiGroupApp.SystemApiGroup.SystemApi
	sysParamsApi        = api.ApiGroupApp.SystemApiGroup.SysParamsApi
	authorityApi        = api.ApiGroupApp.SystemApiGroup.AuthorityApi
	apiRouterApi        = api.ApiGroupApp.SystemApiGroup.SystemApiApi
	dictionaryApi       = api.ApiGroupApp.SystemApiGroup.DictionaryApi
	authorityBtnApi     = api.ApiGroupApp.SystemApiGroup.AuthorityBtnApi
	authorityMenuApi    = api.ApiGroupApp.SystemApiGroup.AuthorityMenuApi

	operationRecordApi  = api.ApiGroupApp.SystemApiGroup.OperationRecordApi

	dictionaryDetailApi = api.ApiGroupApp.SystemApiGroup.DictionaryDetailApi

	exportTemplateApi   = api.ApiGroupApp.SystemApiGroup.SysExportTemplateApi
	sysVersionApi       = api.ApiGroupApp.SystemApiGroup.SysVersionApi
)
