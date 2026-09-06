package router

import (
	"shack/internal/router/gen"
	"shack/internal/router/menu"
	"shack/internal/router/system"
	"shack/internal/router/systemRbac"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	Gen gen.RouterGroup
	Menu menu.RouterGroup
	System     system.RouterGroup
	SystemRbac systemRbac.RouterGroup
}
