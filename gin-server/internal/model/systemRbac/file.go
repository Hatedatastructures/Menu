package systemRbac

import (
	"time"
)



// StorageType 存储类型（方便以后扩展）
type StorageType string

const (
	StorageLocal StorageType = "local"
	StorageOSS   StorageType = "oss"
)

// SysFile 文件表
type SysFile struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// ===== 基本信息 =====
	FileName string `gorm:"size:255;comment:原始文件名" json:"fileName"`
	NewName  string `gorm:"size:255;comment:存储文件名(重命名)" json:"newName"`

	// ===== 路径信息 =====
	FilePath string `gorm:"size:500;comment:本地存储路径" json:"filePath"`
	FileURL  string `gorm:"size:500;comment:访问URL" json:"fileUrl"`

	// ===== 类型信息 =====
	FileExt  string   `gorm:"size:20;comment:文件扩展名" json:"fileExt"`
	MimeType string   `gorm:"size:100;comment:MIME类型" json:"mimeType"`
	FileType uint `gorm:"size:20;comment:文件分类" json:"fileType"`

	// ===== 大小信息 =====
	Size int64 `gorm:"comment:文件大小(字节)" json:"size"`

	// ===== 存储信息 =====
	Storage StorageType `gorm:"size:20;default:local;comment:存储方式" json:"storage"`


	// ===== 用户信息 =====
	CreatedBy uint `gorm:"comment:上传人ID" json:"createdBy"`

	// ===== 时间 =====
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt *time.Time     `gorm:"index" json:"-"`
}

func (SysFile) TableName() string {
	return "sys_rbac_files"
}