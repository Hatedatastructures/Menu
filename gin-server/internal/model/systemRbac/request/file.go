package request

import (
	"mime/multipart"
)

type FileUploadReq struct {
	File *multipart.FileHeader `form:"file"` // 文件对象
	FileType int `json:"fileType" form:"fileType"` // 文件分类(1图片/2文档/3视频/4音频/5其他)
}

type MultiFileUploadReq struct {
	Files []*multipart.FileHeader `form:"files"` // 文件对象数组
	FileType int `json:"fileType" form:"fileType"` // 文件分类(1图片/2文档/3视频/4音频/5其他)
}

type GetFileDetailReq struct {
	Id uint `json:"id" form:"id"`
}

type FileDownloadReq struct {
	Id uint `json:"id" form:"id"`
}

type GetFileListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	FileType int `json:"fileType" form:"fileType"`
	Storage string `json:"storage" form:"storage"`
	CreatedBy uint `json:"createdBy" form:"createdBy"`
	Keyword string `json:"keyword" form:"keyword"`
}

type DeleteFileReq struct {
	Id uint `json:"id" form:"id"`
}

type BatchDeleteFilesReq struct {
	Ids []uint `json:"ids" form:"ids"` // 文件ID数组
}

type UpdateFileNameReq struct {
	Id uint `json:"id" form:"id"`
	FileName string `json:"fileName" form:"fileName"` // 新的文件名
}
