package systemRbac

import (
	"shack/internal/model/common"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// API模型
type Api struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Path      string         `gorm:"size:100;not null;uniqueIndex:idx_api" json:"path"`
	Method    string         `gorm:"size:10;not null;uniqueIndex:idx_api" json:"method"`
	Status    bool           `gorm:"default:true" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 自定义表名
func (Api) TableName() string {
	return "sys_rbac_apis"
}

// Router 路由模型 - 用于前端路由展示
type Router struct {
	ID            int      `json:"id,omitempty" uri:"id"`
	Name          string   `json:"name,omitempty"`
	Icon          string   `json:"icon,omitempty"`
	Path          string   `json:"path,omitempty"`
	ComponentPath string   `json:"componentPath,omitempty"`
	ParentID      int      `json:"parentId,omitempty"`
	RouterOrder   int      `json:"routerOrder,omitempty"`
	Children      []Router `json:"children,omitempty" gorm:"-"`
	Hidden        *bool    `json:"hidden"`
	Required      *bool    `json:"required"`
}

func (Router) TableName() string {
	return "sys_rbac_routers"
}

// RegisterApiParam API注册参数
type RegisterApiParam struct {
	ApiMethod  common.HttpType
	ApiUrl     string
	ApiComment string
	Handle     func(ctx *gin.Context)
}
