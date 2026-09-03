//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type DictRouter struct{}

func (s *DictRouter) InitDictRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	dictRouter := Router.Group("")
	utils.RegisterApi(dictRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionaries", "获取字典列表-后台使用", dictApi.GetDictionaryListHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/dictionaries", "创建字典-后台使用", dictApi.CreateDictionaryHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionaries/:id", "字典详情-后台使用", dictApi.GetDictionaryHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/dictionaries/:id", "更新字典-后台使用", dictApi.UpdateDictionaryHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/dictionaries/:id", "删除字典-后台使用", dictApi.DeleteDictionaryHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dicts", "批量字典查询-前台/后台使用", dictApi.GetDictsByTypesHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dicts/all", "获取全部字典-前台/后台使用", dictApi.GetAllDictsHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/dicts/refresh", "字典缓存刷新-后台使用", dictApi.RefreshDictsHandler),
	)
	return dictRouter
}
