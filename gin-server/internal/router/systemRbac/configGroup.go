// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigGroupRouter struct{}

func (s *ConfigGroupRouter) InitConfigGroupRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configGroupRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configGroupRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/groups", "创建配置分组-后台使用", configGroupApi.CreateConfigGroupHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/groups", "配置分组列表-后台使用", configGroupApi.GetConfigGroupListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/groups/:id", "获取配置分组详情-后台使用", configGroupApi.GetConfigGroupDetailHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/config/groups/:id", "更新配置分组-后台使用", configGroupApi.UpdateConfigGroupHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/config/groups/:id", "删除配置分组-后台使用", configGroupApi.DeleteConfigGroupHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/config/groups", "批量删除配置分组-后台使用", configGroupApi.BatchDeleteConfigGroupsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/:code", "根据编码获取配置-后台使用", configGroupApi.GetConfigByCodeHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/:code", "保存配置-后台使用", configGroupApi.SaveConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/refresh", "刷新配置缓存-后台使用", configGroupApi.RefreshConfigCacheHandler),
	)
	return configGroupRouter
}
