package menu

import (
	_ "shack/internal/model/menu/response"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type HealthApi struct{}

// HealthzHandler
// @Tags menuhealthApi
// @Summary HealthzHandler 系统健康检查模块 对应 C++ server /healthz 和 /readyz 健康检查
// @Description HealthzHandler 系统健康检查模块 对应 C++ server /healthz 和 /readyz 健康检查
// @Success 200 {object} vo.Result{data=_.HealthzRes}
// @Router /healthz [GET]
func (s *HealthApi) HealthzHandler(c *gin.Context) {
	data, err := healthService.Healthz(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ReadyzHandler
// @Tags menuhealthApi
// @Summary ReadyzHandler 就绪检查
// @Description ReadyzHandler 就绪检查
// @Success 200 {object} vo.Result{data=_.ReadyzRes}
// @Router /readyz [GET]
func (s *HealthApi) ReadyzHandler(c *gin.Context) {
	data, err := healthService.Readyz(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
