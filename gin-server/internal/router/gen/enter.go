package gen

import api "shack/internal/api/v1"

type RouterGroup struct {
	ApikeyRouter
	GenerateRouter
	HistoryRouter
}

var (
	apikeyApi = api.ApiGroupApp.GenApiGroup.ApikeyApi
	generateApi = api.ApiGroupApp.GenApiGroup.GenerateApi
	historyApi = api.ApiGroupApp.GenApiGroup.HistoryApi
)
