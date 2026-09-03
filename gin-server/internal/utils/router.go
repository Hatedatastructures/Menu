package utils

import (
	"fmt"
	"shack/internal/global"
	"shack/internal/model/common"
	"shack/internal/model/systemRbac"

	"github.com/gin-gonic/gin"
)

// GetRouterTree
// -------- 返回组装后的路由 --------
func GetRouterTree(routers *[]systemRbac.Router) []systemRbac.Router {
	routerMap := map[int][]systemRbac.Router{}
	for _, router := range *routers {
		routerMap[router.ParentID] = append(routerMap[router.ParentID], router)
	}
	var routerTree []systemRbac.Router
	for _, router := range *routers {
		router.Children = routerMap[router.ID]
		if router.ParentID == -1 {
			routerTree = append(routerTree, router)
		}
	}
	return routerTree
}

func NewRegisterApiParam(apiMethod common.HttpType, apiUrl, apiComment string, handle func(ctx *gin.Context)) systemRbac.RegisterApiParam {
	return systemRbac.RegisterApiParam{
		ApiMethod:  apiMethod,
		ApiUrl:     apiUrl,
		ApiComment: apiComment,
		Handle:     handle,
	}
}

// RegisterApi 根据Api接口简介是否存在,进行判断是否已注册
func RegisterApi(router gin.IRoutes, apiSet map[string]struct{}, prefix string, params ...systemRbac.RegisterApiParam) {

	for _, param := range params {
		// 使用 path-method 作为唯一键
		fullPath := fmt.Sprintf("%s%s", prefix, param.ApiUrl)
		uniqueKey := fmt.Sprintf("%s-%s", fullPath, param.ApiMethod.MethodString())

		_, exist := apiSet[uniqueKey]
		if !exist {
			api := systemRbac.Api{
				Path:   fullPath,
				Name:   param.ApiComment,
				Method: param.ApiMethod.MethodString(),
			}
			// 使用 FirstOrCreate 避免重复插入错误
			if err := global.GVA_DB.Where("path = ? AND method = ?", fullPath, param.ApiMethod.MethodString()).
				FirstOrCreate(&api).Error; err != nil {
				panic(fmt.Sprintf("Api注册失败,%s", err.Error()))
			}
			apiSet[uniqueKey] = struct{}{}
		}
		switch param.ApiMethod {
		case common.HttpGet:
			router.GET(param.ApiUrl, param.Handle)
		case common.HttpPost:
			router.POST(param.ApiUrl, param.Handle)
		case common.HttpPatch:
			router.PATCH(param.ApiUrl, param.Handle)
		case common.HttpDelete:
			router.DELETE(param.ApiUrl, param.Handle)
		default:
			router.PUT(param.ApiUrl, param.Handle)
		}
	}
}
