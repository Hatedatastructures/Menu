package initialize

import (
	"net/http"
	"os"

	"shack/docs"
	"shack/internal/global"
	"shack/internal/middleware"
	model_systemRbac "shack/internal/model/systemRbac"

	"shack/internal/router"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	// 静态文件服务 - server/statics 目录
	Router.Use(static.Serve("/statics", static.LocalFile("./statics", true)))


	systemRouter := router.RouterGroupApp.System


	Router.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GVA_CONFIG.Local.StorePath)}) // Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	Router.Use(middleware.Cors())
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	// global.GVA_LOG.Info("use middleware cors")
	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)

	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.CaptureResponseBody()).Use(middleware.EncryptResponse())
	PublicGroup.Use(middleware.CaptureResponseBody()).Use(middleware.EncryptResponse())
	apiSet := getApiSetStructCommonAsKey()

	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
	}

	{
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)               // 注册功能api路由
		systemRouter.InitJwtRouter(PrivateGroup)                            // jwt相关路由
		systemRouter.InitCasbinRouter(PrivateGroup)                         // 权限相关路由
	}

	initAutoRouter( PublicGroup,PrivateGroup,apiSet)

	//插件路由安装
	InstallPlugin(PrivateGroup, PublicGroup, Router)

	// 注册业务路由
	initBizRouter(PrivateGroup, PublicGroup)

	global.GVA_ROUTERS = Router.Routes()

	global.GVA_LOG.Info("router register success")
	return Router
}
func getApiSetStructCommonAsKey() map[string]struct{} {
	// 直接查询数据库获取所有启用的API
	var apis []model_systemRbac.Api
	err := global.GVA_DB.Where("status = ?", true).Find(&apis).Error
	if err != nil {
		return nil
	}
	apiSet := make(map[string]struct{}, len(apis))
	for _, api := range apis {
		// 使用 Name (api_comment) 作为key，与 RegisterApi 保持一致
		apiSet[api.Name] = struct{}{}
	}
	return apiSet
}