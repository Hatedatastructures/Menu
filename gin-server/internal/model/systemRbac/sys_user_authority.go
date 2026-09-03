package systemRbac

// SysUserAuthority 是 sysUser 和 sysAuthority 的连接表
type SysUserAuthority struct {
	SysUserId               uint `gorm:"column:sys_user_id;primaryKey;index"`
	SysAuthorityAuthorityId uint `gorm:"column:sys_authority_authority_id;primaryKey;index"`
}

func (s *SysUserAuthority) TableName() string {
	return "sys_rbac_user_authority"
}
