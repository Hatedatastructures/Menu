package systemRbac

import (
	"time"
)

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	ID              uint            `gorm:"primarykey"`                                                         // 主键ID
	AuthorityId     uint            `json:"authorityId" gorm:"column:authority_id;unique;comment:角色ID;size:90"` // 角色ID (业务主键)
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                                   // 角色名
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`                                      // 父角色ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_rbac_data_authority_id;foreignKey:AuthorityId;joinForeignKey:sys_authority_authority_id;references:AuthorityId;joinReferences:data_authority_authority_id"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_rbac_authority_menus;foreignKey:AuthorityId;joinForeignKey:sys_authority_authority_id;references:ID;joinReferences:sys_base_menu_id"`
	Users           []User          `json:"-" gorm:"many2many:sys_rbac_user_authority;foreignKey:AuthorityId;joinForeignKey:sys_authority_authority_id;references:ID;joinReferences:sys_user_id"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
}

func (SysAuthority) TableName() string {
	return "sys_rbac_authorities"
}
