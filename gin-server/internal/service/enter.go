package service

import (
	"shack/internal/service/gen"
	"shack/internal/service/menu"
	"shack/internal/service/system"
	"shack/internal/service/systemRbac"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	GenServiceGroup gen.ServiceGroup
	MenuServiceGroup menu.ServiceGroup
	SystemRbacServiceGroup systemRbac.ServiceGroup
	SystemServiceGroup     system.ServiceGroup
}
