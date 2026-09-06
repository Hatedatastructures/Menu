package response

type FileUploadRes struct {
	Id uint `json:"id"` // 文件ID
	FileName string `json:"fileName"` // 原始文件名
	NewName string `json:"newName"` // 存储文件名
	FilePath string `json:"filePath"` // 本地存储路径
	FileUrl string `json:"fileUrl"` // 访问URL
	FileExt string `json:"fileExt"` // 文件扩展名
	MimeType string `json:"mimeType"` // MIME类型
	FileType uint `json:"fileType"` // 文件分类
	Size int64 `json:"size"` // 文件大小(字节)
	Storage string `json:"storage"` // 存储方式
	CreatedAt string `json:"createdAt"` // 创建时间
}

type MultiFileUploadRes struct {
	List []MultiFileUploadResList `json:"list"` // []
	SuccessCount int `json:"successCount"` // 成功数量
	FailCount int `json:"failCount"` // 失败数量
}

type MultiFileUploadResList struct {
	Id uint `json:"id"` // 文件ID
	FileName string `json:"fileName"` // 原始文件名
	NewName string `json:"newName"` // 存储文件名
	FilePath string `json:"filePath"` // 本地存储路径
	FileUrl string `json:"fileUrl"` // 访问URL
	FileExt string `json:"fileExt"` // 文件扩展名
	MimeType string `json:"mimeType"` // MIME类型
	FileType uint `json:"fileType"` // 文件分类
	Size int64 `json:"size"` // 文件大小(字节)
	Storage string `json:"storage"` // 存储方式
}

type GetFileDetailRes struct {
	Id uint `json:"id"` // 文件ID
	FileName string `json:"fileName"` // 原始文件名
	NewName string `json:"newName"` // 存储文件名
	FilePath string `json:"filePath"` // 本地存储路径
	FileUrl string `json:"fileUrl"` // 访问URL
	FileExt string `json:"fileExt"` // 文件扩展名
	MimeType string `json:"mimeType"` // MIME类型
	FileType uint `json:"fileType"` // 文件分类
	Size int64 `json:"size"` // 文件大小(字节)
	Storage string `json:"storage"` // 存储方式
	CreatedBy uint `json:"createdBy"` // 上传人ID
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type FileDownloadRes struct {
	File []byte `json:"file"` // 文件二进制数据
}

type GetFileListRes struct {

	FileType int `json:"fileType"`
	Storage string `json:"storage"`
	CreatedBy uint `json:"createdBy"`
	Keyword string `json:"keyword"`
	List []GetFileListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
	Page int `json:"page"` // 当前页码
	Size int `json:"size"` // 每页数量
}

type GetFileListResList struct {
	Id uint `json:"id"` // 文件ID
	FileName string `json:"fileName"` // 原始文件名
	NewName string `json:"newName"` // 存储文件名
	FileUrl string `json:"fileUrl"` // 访问URL
	FileExt string `json:"fileExt"` // 文件扩展名
	FileType uint `json:"fileType"` // 文件分类
	Size int64 `json:"size"` // 文件大小(字节)
	Storage string `json:"storage"` // 存储方式
	CreatedBy uint `json:"createdBy"` // 上传人ID
	CreatedAt string `json:"createdAt"` // 创建时间
}

type BatchDeleteFilesRes struct {
	SuccessCount int `json:"successCount"` // 成功删除数量
	FailCount int `json:"failCount"` // 删除失败数量
}

type UpdateFileNameRes struct {
	Id uint `json:"id"` // 文件ID
	FileName string `json:"fileName"` // 更新后的文件名
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetFileTypeStatsRes struct {
	List []GetFileTypeStatsResList `json:"list"` // []
	TotalCount int64 `json:"totalCount"` // 文件总数
	TotalSize int64 `json:"totalSize"` // 总存储大小(字节)
}

type GetFileTypeStatsResList struct {
	FileType uint `json:"fileType"` // 文件分类(1图片/2文档/3视频/4音频/5其他)
	Count int64 `json:"count"` // 文件数量
	TotalSize int64 `json:"totalSize"` // 总大小(字节)
}
