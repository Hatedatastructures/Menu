package menu

import (
	"github.com/gin-gonic/gin"
	res "shack/internal/model/menu/response"
)

type HealthService struct{}

func (s *HealthService) Healthz(ctx *gin.Context) (rs res.HealthzRes, err error) {
	return res.HealthzRes{Status: "ok"}, nil
}

func (s *HealthService) Readyz(ctx *gin.Context) (rs res.ReadyzRes, err error) {
	return res.ReadyzRes{Status: "ready"}, nil
}
