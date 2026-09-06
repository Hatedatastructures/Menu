// router 解析
package systemRbac

import (
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type DictDetailRouter struct{}

func (s *DictDetailRouter) InitDictDetailRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	dictDetailRouter := Router.Group("")
	utils.RegisterApi(dictDetailRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionary-details", "获取字典详情列表-后台使用", dictDetailApi.GetDictionaryDetailListHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/dictionary-details", "创建字典详情-后台使用", dictDetailApi.CreateDictionaryDetailHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionary-details/id/:id", "字典详情详情-后台使用", dictDetailApi.GetDictionaryDetailHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/dictionary-details/id/:id", "更新字典详情-后台使用", dictDetailApi.UpdateDictionaryDetailHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/dictionary-details/id/:id", "删除字典详情-后台使用", dictDetailApi.DeleteDictionaryDetailHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionary-details/dictionary/:dictionaryId", "按字典ID获取字典全部内容-后台使用", dictDetailApi.GetDictionaryListByIdHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionary-details/type/:type", "按字典type获取字典全部内容-后台使用", dictDetailApi.GetDictionaryListByTypeHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionary-details/:dictionaryId/value/:value", "按字典ID和value获取单条字典内容-后台使用", dictDetailApi.GetDictionaryInfoByValueHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/dictionary-details/type/:type/value/:value", "按字典type和value获取单条字典内容-后台使用", dictDetailApi.GetDictionaryInfoByTypeValueHandler),
	)
	return dictDetailRouter
}
