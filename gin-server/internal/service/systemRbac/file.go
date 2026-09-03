package systemRbac

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"shack/internal/global"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	"shack/internal/model/systemRbac"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// FileTypeMap 文件类型映射
	FileTypeImage   = 1 // 图片
	FileTypeDoc     = 2 // 文档
	FileTypeVideo   = 3 // 视频
	FileTypeAudio   = 4 // 音频
	FileTypeOther   = 5 // 其他

	// RedisKeyFileStats 文件统计缓存Key
	RedisKeyFileStats = "file:stats:"
	RedisExpireStats  = 300 // 5分钟
)

// FileService 文件服务
type FileService struct{}

// getFileTypeByExt 根据扩展名获取文件类型
func getFileTypeByExt(ext string) uint {
	ext = strings.ToLower(ext)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico"}
	docExts := []string{".doc", ".docx", ".pdf", ".txt", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".rtf"}
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm"}
	audioExts := []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma"}

	for _, e := range imageExts {
		if ext == e {
			return FileTypeImage
		}
	}
	for _, e := range docExts {
		if ext == e {
			return FileTypeDoc
		}
	}
	for _, e := range videoExts {
		if ext == e {
			return FileTypeVideo
		}
	}
	for _, e := range audioExts {
		if ext == e {
			return FileTypeAudio
		}
	}
	return FileTypeOther
}

// getFileSavePath 获取文件保存路径
func getFileSavePath(fileType uint) string {
	typeNames := map[uint]string{
		FileTypeImage: "images",
		FileTypeDoc:   "documents",
		FileTypeVideo: "videos",
		FileTypeAudio: "audios",
		FileTypeOther: "others",
	}
	return typeNames[fileType]
}

// saveFileToDisk 保存文件到磁盘
func saveFileToDisk(file *multipart.FileHeader, subDir string) (string, string, error) {
	// 创建上传目录
	uploadDir := filepath.Join("statics", subDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", "", err
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	newName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(uploadDir, newName)

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(savePath)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	// 复制内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", "", err
	}

	// 返回相对路径和访问URL
	relativePath := filepath.Join(subDir, newName)
	fileURL := "/statics/" + strings.ReplaceAll(relativePath, "\\", "/")

	return relativePath, fileURL, nil
}

// deleteFileFromDisk 删除磁盘文件
func deleteFileFromDisk(filePath string) error {
	fullPath := filepath.Join("statics", filePath)
	return os.Remove(fullPath)
}

// getCurrentUserID 从上下文获取当前用户ID
func getCurrentUserID(c *gin.Context) uint {
	if user, exists := c.Get("userId"); exists {
		if id, ok := user.(uint); ok {
			return id
		}
	}
	return 0
}

// FileUpload 单文件上传
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) FileUpload(
	c *gin.Context,
	r req.FileUploadReq,
) (rs res.FileUploadRes, err error) {
	if r.File == nil {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请选择文件")
	}

	// 自动判断文件类型
	fileType := uint(r.FileType)
	if fileType == 0 {
		ext := filepath.Ext(r.File.Filename)
		fileType = getFileTypeByExt(ext)
	}

	// 保存文件
	subDir := getFileSavePath(fileType)
	filePath, fileURL, err := saveFileToDisk(r.File, subDir)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "文件保存失败")
	}

	// 获取文件扩展名和MIME类型
	ext := strings.ToLower(filepath.Ext(r.File.Filename))
	mimeType := r.File.Header.Get("Content-Type")

	// 创建数据库记录
	userID := getCurrentUserID(c)
	file := systemRbac.SysFile{
		FileName:  r.File.Filename,
		NewName:   filepath.Base(filePath),
		FilePath:  filePath,
		FileURL:   fileURL,
		FileExt:   ext,
		MimeType:  mimeType,
		FileType:  fileType,
		Size:      r.File.Size,
		Storage:   systemRbac.StorageLocal,
		CreatedBy: userID,
	}

	if err := global.GVA_DB.Create(&file).Error; err != nil {
		// 删除已上传的文件
		deleteFileFromDisk(filePath)
		return rs, biz_err.New(biz_err.DB_ERROR, "保存文件记录失败")
	}

	// 清除统计缓存
	clearFileStatsCache()

	rs = res.FileUploadRes{
		Id:        file.ID,
		FileName:  file.FileName,
		NewName:   file.NewName,
		FilePath:  file.FilePath,
		FileUrl:   file.FileURL,
		FileExt:   file.FileExt,
		MimeType:  file.MimeType,
		FileType:  file.FileType,
		Size:      file.Size,
		Storage:   string(file.Storage),
		CreatedAt: file.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// MultiFileUpload 多文件上传
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) MultiFileUpload(
	c *gin.Context,
	r req.MultiFileUploadReq,
) (rs res.MultiFileUploadRes, err error) {
	if len(r.Files) == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请选择文件")
	}

	var successList []res.MultiFileUploadResList
	successCount := 0
	failCount := 0
	userID := getCurrentUserID(c)

	for _, file := range r.Files {
		if file == nil {
			failCount++
			continue
		}

		// 自动判断文件类型
		fileType := uint(r.FileType)
		if fileType == 0 {
			ext := filepath.Ext(file.Filename)
			fileType = getFileTypeByExt(ext)
		}

		// 保存文件
		subDir := getFileSavePath(fileType)
		filePath, fileURL, err := saveFileToDisk(file, subDir)
		if err != nil {
			failCount++
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		mimeType := file.Header.Get("Content-Type")

		// 创建数据库记录
		f := systemRbac.SysFile{
			FileName:  file.Filename,
			NewName:   filepath.Base(filePath),
			FilePath:  filePath,
			FileURL:   fileURL,
			FileExt:   ext,
			MimeType:  mimeType,
			FileType:  fileType,
			Size:      file.Size,
			Storage:   systemRbac.StorageLocal,
			CreatedBy: userID,
		}

		if err := global.GVA_DB.Create(&f).Error; err != nil {
			deleteFileFromDisk(filePath)
			failCount++
			continue
		}

		successList = append(successList, res.MultiFileUploadResList{
			Id:        f.ID,
			FileName:  f.FileName,
			NewName:   f.NewName,
			FilePath:  f.FilePath,
			FileUrl:   f.FileURL,
			FileExt:   f.FileExt,
			MimeType:  f.MimeType,
			FileType:  f.FileType,
			Size:      f.Size,
			Storage:   string(f.Storage),
		})
		successCount++
	}

	// 清除统计缓存
	if successCount > 0 {
		clearFileStatsCache()
	}

	rs = res.MultiFileUploadRes{
		List:         successList,
		SuccessCount: successCount,
		FailCount:    failCount,
	}

	return rs, nil
}

// GetFileDetail 获取文件详情
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) GetFileDetail(
	c *gin.Context,
	r req.GetFileDetailReq,
) (rs res.GetFileDetailRes, err error) {
	var file systemRbac.SysFile
	if err := global.GVA_DB.First(&file, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "文件不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "获取文件信息失败")
	}

	rs = res.GetFileDetailRes{
		Id:        file.ID,
		FileName:  file.FileName,
		NewName:   file.NewName,
		FilePath:  file.FilePath,
		FileUrl:   file.FileURL,
		FileExt:   file.FileExt,
		MimeType:  file.MimeType,
		FileType:  file.FileType,
		Size:      file.Size,
		Storage:   string(file.Storage),
		CreatedBy: file.CreatedBy,
		CreatedAt: file.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: file.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// FileDownload 文件下载
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) FileDownload(
	c *gin.Context,
	r req.FileDownloadReq,
) (res.FileDownloadRes, error) {
	var file systemRbac.SysFile
	if err := global.GVA_DB.First(&file, r.Id).Error; err != nil {
		return res.FileDownloadRes{}, biz_err.New(biz_err.PARAM_ERROR, "文件不存在")
	}

	// 设置响应头
	fullPath := filepath.Join("statics", file.FilePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "public")

	// 使用文件路径返回
	c.File(fullPath)

	return res.FileDownloadRes{}, nil
}

// GetFileList 文件列表
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) GetFileList(
	c *gin.Context,
	r req.GetFileListReq,
) (rs res.GetFileListRes, err error) {
	// 默认分页参数
	page := r.Page
	if page < 1 {
		page = 1
	}
	size := r.Size
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.SysFile{})

	// 条件过滤
	if r.FileType > 0 {
		db = db.Where("file_type = ?", r.FileType)
	}
	if r.Storage != "" {
		db = db.Where("storage = ?", r.Storage)
	}
	if r.CreatedBy > 0 {
		db = db.Where("created_by = ?", r.CreatedBy)
	}
	if r.Keyword != "" {
		db = db.Where("file_name LIKE ?", "%"+r.Keyword+"%")
	}

	// 查询总数
	var total int64
	db.Count(&total)

	// 分页查询
	var files []systemRbac.SysFile
	offset := (page - 1) * size
	if err := db.Order("created_at DESC").Offset(offset).Limit(size).Find(&files).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取文件列表失败")
	}

	// 转换结果
	var list []res.GetFileListResList
	for _, f := range files {
		list = append(list, res.GetFileListResList{
			Id:        f.ID,
			FileName:  f.FileName,
			NewName:   f.NewName,
			FileUrl:   f.FileURL,
			FileExt:   f.FileExt,
			FileType:  f.FileType,
			Size:      f.Size,
			Storage:   string(f.Storage),
			CreatedBy: f.CreatedBy,
			CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetFileListRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}

	return rs, nil
}

// DeleteFile 删除文件
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) DeleteFile(
	c *gin.Context,
	r req.DeleteFileReq,
) error {
	var file systemRbac.SysFile
	if err := global.GVA_DB.First(&file, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "文件不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "获取文件信息失败")
	}

	// 删除磁盘文件
	if err := deleteFileFromDisk(file.FilePath); err != nil {
		// 记录日志但继续删除数据库记录
		global.GVA_LOG.Warn("删除磁盘文件失败", zap.String("path", file.FilePath))
	}

	// 删除数据库记录
	if err := global.GVA_DB.Delete(&file).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除文件记录失败")
	}

	// 清除统计缓存
	clearFileStatsCache()

	return nil
}

// BatchDeleteFiles 批量删除文件
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) BatchDeleteFiles(
	c *gin.Context,
	r req.BatchDeleteFilesReq,
) (rs res.BatchDeleteFilesRes, err error) {
	if len(r.Ids) == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请选择要删除的文件")
	}

	var files []systemRbac.SysFile
	if err := global.GVA_DB.Where("id IN ?", r.Ids).Find(&files).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取文件信息失败")
	}

	successCount := 0
	failCount := 0

	for _, file := range files {
		// 删除磁盘文件
		if err := deleteFileFromDisk(file.FilePath); err != nil {
			global.GVA_LOG.Warn("删除磁盘文件失败", zap.String("path", file.FilePath))
		}

		// 删除数据库记录
		if err := global.GVA_DB.Delete(&file).Error; err != nil {
			failCount++
			continue
		}
		successCount++
	}

	// 清除统计缓存
	if successCount > 0 {
		clearFileStatsCache()
	}

	rs = res.BatchDeleteFilesRes{
		SuccessCount: successCount,
		FailCount:    failCount,
	}

	return rs, nil
}

// UpdateFileName 更新文件名
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) UpdateFileName(
	c *gin.Context,
	r req.UpdateFileNameReq,
) (rs res.UpdateFileNameRes, err error) {
	if r.FileName == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "文件名不能为空")
	}

	var file systemRbac.SysFile
	if err := global.GVA_DB.First(&file, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "文件不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "获取文件信息失败")
	}

	// 更新文件名
	if err := global.GVA_DB.Model(&file).Update("file_name", r.FileName).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "更新文件名失败")
	}

	rs = res.UpdateFileNameRes{
		Id:        file.ID,
		FileName:  r.FileName,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// GetFileTypeStats 获取文件分类统计
// Auth: shack
// Time: 2026年04月02日
func (s *FileService) GetFileTypeStats(
	c *gin.Context,
) (rs res.GetFileTypeStatsRes, err error) {
	// 尝试从Redis获取缓存
	cacheKey := RedisKeyFileStats + "all"
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(context.Background(), cacheKey).Result()
		if err == nil && cached != "" {
			// 解析缓存数据 (实际生产中可使用JSON反序列化)
			_ = cached
			// 此处可实现缓存命中后直接返回
		}
	}

	// 从数据库查询统计
	type StatResult struct {
		FileType uint
		Count    int64
		TotalSize int64
	}

	var stats []StatResult
	if err := global.GVA_DB.Model(&systemRbac.SysFile{}).
		Select("file_type, COUNT(*) as count, COALESCE(SUM(size), 0) as total_size").
		Group("file_type").
		Order("file_type").
		Find(&stats).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取文件统计失败")
	}

	// 计算总计
	var totalCount int64
	var totalSize int64

	var list []res.GetFileTypeStatsResList
	for _, s := range stats {
		list = append(list, res.GetFileTypeStatsResList{
			FileType:  s.FileType,
			Count:     s.Count,
			TotalSize: s.TotalSize,
		})
		totalCount += s.Count
		totalSize += s.TotalSize
	}

	rs = res.GetFileTypeStatsRes{
		List:       list,
		TotalCount: totalCount,
		TotalSize:  totalSize,
	}

	// 缓存结果
	if global.GVA_REDIS != nil {
		// 实际生产中可序列化后缓存
		global.GVA_REDIS.Set(context.Background(), cacheKey, "1", RedisExpireStats*time.Second)
	}

	return rs, nil
}

// clearFileStatsCache 清除文件统计缓存
func clearFileStatsCache() {
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(context.Background(), RedisKeyFileStats+"all")
	}
}