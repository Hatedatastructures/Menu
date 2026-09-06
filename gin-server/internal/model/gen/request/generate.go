package request

type GenerateReq struct {
	FileId uint `json:"fileId" form:"fileId"` // 上传的文件ID(来自sys_rbac_files)
	FileName string `json:"fileName" form:"fileName"` // 原始文件名
}

type GetGenerateResultReq struct {
	Id uint `json:"id" form:"id"`
}
