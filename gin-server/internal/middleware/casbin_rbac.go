package middleware

import (
	"strconv"
	"strings"

	"shack/internal/global"
	"shack/internal/model/common/response"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	b := global.GVA_CONFIG.Develop.IsCasbinMode
	if b == false {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		//获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin() // 判断策略中是否存在
		if e == nil {
			response.FailWithDetailed(gin.H{}, "权限系统未初始化", c)
			c.Abort()
			return
		}
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
