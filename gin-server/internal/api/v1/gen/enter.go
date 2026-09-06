package gen

import "shack/internal/service"

type ApiGroup struct {
	ApikeyApi
	GenerateApi
	HistoryApi
}

var (
	apikeyService = service.ServiceGroupApp.GenServiceGroup.ApikeyService
	generateService = service.ServiceGroupApp.GenServiceGroup.GenerateService
	historyService = service.ServiceGroupApp.GenServiceGroup.HistoryService
)
