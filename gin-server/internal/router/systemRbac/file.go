//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type FileRouter struct{}

func (s *FileRouter) InitFileRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	fileRouter := Router.Group("")
	utils.RegisterApi(fileRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/files/upload", "文件上传-后台使用", fileApi.FileUploadHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/files/uploads", "多文件上传-后台使用", fileApi.MultiFileUploadHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/files/:id", "获取文件详情-后台使用", fileApi.GetFileDetailHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/files/:id/download", "文件下载-后台使用", fileApi.FileDownloadHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/files", "文件列表-后台使用", fileApi.GetFileListHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/files/:id", "删除文件-后台使用", fileApi.DeleteFileHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/files/batch", "批量删除文件-后台使用", fileApi.BatchDeleteFilesHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/files/:id", "更新文件名-后台使用", fileApi.UpdateFileNameHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/files/types", "文件分类统计-后台使用", fileApi.GetFileTypeStatsHandler),
	)
	return fileRouter
}
