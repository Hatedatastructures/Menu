package systemRbac

import (
	"strconv"

	"github.com/gin-gonic/gin"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"
	biz_err "shack/internal/error"
)

type FileApi struct{}

// FileUploadHandler
// @Tags systemRbacfileApi
// @Summary FileUploadHandler 文件上传-后台使用
// @Description FileUploadHandler 文件上传-后台使用
// @Param file formData file true "文件"
// @Param fileType formData int false "文件分类(1图片/2文档/3视频/4音频/5其他)"
// @Success 200 {object} vo.Result{data=_.FileUploadRes}
// @Router /api/v1/files/upload [POST]
func (s *FileApi) FileUploadHandler(c *gin.Context) {
	var req req.FileUploadReq
	// formData 参数
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := fileService.FileUpload(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// MultiFileUploadHandler
// @Tags systemRbacfileApi
// @Summary MultiFileUploadHandler 多文件上传-后台使用
// @Description MultiFileUploadHandler 多文件上传-后台使用
// @Param files formData file true "文件数组"
// @Param fileType formData int false "文件分类(1图片/2文档/3视频/4音频/5其他)"
// @Success 200 {object} vo.Result{data=_.MultiFileUploadRes}
// @Router /api/v1/files/uploads [POST]
func (s *FileApi) MultiFileUploadHandler(c *gin.Context) {
	var req req.MultiFileUploadReq
	// formData 参数
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := fileService.MultiFileUpload(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetFileDetailHandler
// @Tags systemRbacfileApi
// @Summary GetFileDetailHandler 获取文件详情-后台使用
// @Description GetFileDetailHandler 获取文件详情-后台使用
// @Param data body req.GetFileDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetFileDetailRes}
// @Router /api/v1/files/:id [GET]
func (s *FileApi) GetFileDetailHandler(c *gin.Context) {
	var req req.GetFileDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := fileService.GetFileDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// FileDownloadHandler
// @Tags systemRbacfileApi
// @Summary FileDownloadHandler 文件下载-后台使用
// @Description FileDownloadHandler 文件下载-后台使用
// @Param data body req.FileDownloadReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.FileDownloadRes}
// @Router /api/v1/files/:id/download [GET]
func (s *FileApi) FileDownloadHandler(c *gin.Context) {
	var req req.FileDownloadReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	_, err := fileService.FileDownload(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	// FileDownload 方法内部已经通过 c.File() 发送了文件，这里不需要再返回JSON
}

// GetFileListHandler
// @Tags systemRbacfileApi
// @Summary GetFileListHandler 文件列表-后台使用
// @Description GetFileListHandler 文件列表-后台使用
// @Param data body req.GetFileListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetFileListRes}
// @Router /api/v1/files [GET]
func (s *FileApi) GetFileListHandler(c *gin.Context) {
	var req req.GetFileListReq
	// query 参数
	{
		val := c.Query("page")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Page = parsed

		}
	}
	{
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Size = parsed

		}
	}
	{
		val := c.Query("fileType")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.FileType = parsed

		}
	}
	{
		val := c.Query("storage")
		if val != "" {
			req.Storage = val

		}
	}
	{
		val := c.Query("createdBy")
		if val != "" {
			parsed, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.CreatedBy = uint(parsed)

		}
	}
	{
		val := c.Query("keyword")
		if val != "" {
			req.Keyword = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := fileService.GetFileList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteFileHandler
// @Tags systemRbacfileApi
// @Summary DeleteFileHandler 删除文件-后台使用
// @Description DeleteFileHandler 删除文件-后台使用
// @Param data body req.DeleteFileReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/v1/files/:id [DELETE]
func (s *FileApi) DeleteFileHandler(c *gin.Context) {
	var req req.DeleteFileReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := fileService.DeleteFile(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// BatchDeleteFilesHandler
// @Tags systemRbacfileApi
// @Summary BatchDeleteFilesHandler 批量删除文件-后台使用
// @Description BatchDeleteFilesHandler 批量删除文件-后台使用
// @Param data body req.BatchDeleteFilesReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.BatchDeleteFilesRes}
// @Router /api/v1/files/batch [DELETE]
func (s *FileApi) BatchDeleteFilesHandler(c *gin.Context) {
	var req req.BatchDeleteFilesReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := fileService.BatchDeleteFiles(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateFileNameHandler
// @Tags systemRbacfileApi
// @Summary UpdateFileNameHandler 更新文件名-后台使用
// @Description UpdateFileNameHandler 更新文件名-后台使用
// @Param data body req.UpdateFileNameReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateFileNameRes}
// @Router /api/v1/files/:id [PUT]
func (s *FileApi) UpdateFileNameHandler(c *gin.Context) {
	var req req.UpdateFileNameReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = uint(parsed)

	}
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = uint(parsed)

	}
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := fileService.UpdateFileName(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetFileTypeStatsHandler
// @Tags systemRbacfileApi
// @Summary GetFileTypeStatsHandler 文件分类统计-后台使用
// @Description GetFileTypeStatsHandler 文件分类统计-后台使用
// @Success 200 {object} vo.Result{data=_.GetFileTypeStatsRes}
// @Router /api/v1/files/types [GET]
func (s *FileApi) GetFileTypeStatsHandler(c *gin.Context) {
	data, err := fileService.GetFileTypeStats(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}