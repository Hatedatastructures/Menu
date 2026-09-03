package v1

import (
	"shack/internal/api/v1/gen"
	"shack/internal/api/v1/menu"
	"shack/internal/api/v1/system"
	"shack/internal/api/v1/systemRbac"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	GenApiGroup gen.ApiGroup
	MenuApiGroup menu.ApiGroup
	SystemApiGroup     system.ApiGroup
	SystemRbacApiGroup systemRbac.ApiGroup
}
