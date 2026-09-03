//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type TemplateRouter struct{}

func (s *TemplateRouter) InitTemplateRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	templateRouter := Router.Group("")
	utils.RegisterApi(templateRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/system/template", "创建模板-后台使用", templateApi.CreateTemplateHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/system/template/:id", "删除模板-后台使用", templateApi.DeleteTemplateHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/system/template/:id", "更新模板-后台使用", templateApi.UpdateTemplateHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/system/template/:id", "获取模板详情-后台使用", templateApi.GetTemplateDetailHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/system/template/list", "获取模板分页列表-后台使用", templateApi.GetTemplateListHandler),
	)
	return templateRouter
}
