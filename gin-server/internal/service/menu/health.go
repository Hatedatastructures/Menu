package menu

import (
	"github.com/gin-gonic/gin"
	res "shack/internal/model/menu/response"
)

type HealthService struct{}

// 系统健康检查模块 对应 C++ server /healthz 和 /readyz 健康检查
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *HealthService) Healthz(
	ctx *gin.Context,
) (rs res.HealthzRes, err error) {
	return rs, nil
}

// 就绪检查
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *HealthService) Readyz(
	ctx *gin.Context,
) (rs res.ReadyzRes, err error) {
	return rs, nil
}

