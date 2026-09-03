package systemRbac

import (
	"shack/internal/service/systemRbac"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigEncryptApi struct{}

// GetEncryptConfigHandler 获取加密配置(带Redis缓存)
func (s *ConfigEncryptApi) GetEncryptConfigHandler(c *gin.Context) {
	data, err := (&systemRbac.ConfigEncryptService{}).GetSecurityConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}