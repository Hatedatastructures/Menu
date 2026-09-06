package systemRbac

import (
	"errors"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SysAuthorityService 角色服务
type SysAuthorityService struct{}

var SysAuthorityServiceApp = new(SysAuthorityService)

// 创建角色-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) CreateAuthority(
	ctx *gin.Context,
	r req.CreateAuthorityReq,
) (rs res.CreateAuthorityRes, err error) {
	// 参数校验
	if r.AuthorityName == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色名不能为空")
	}

	// 检查角色是否已存在 (仅在指定了 ID 时检查)
	if r.AuthorityId != 0 {
		var count int64
		global.GVA_DB.Model(&systemRbac.SysAuthority{}).Where("authority_id = ?", r.AuthorityId).Count(&count)
		if count > 0 {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "角色ID已存在")
		}
	}

	// 创建角色
	auth := systemRbac.SysAuthority{
		AuthorityId:   r.AuthorityId,
		AuthorityName: r.AuthorityName,
		ParentId:      &r.ParentId,
		DefaultRouter: r.DefaultRouter,
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&auth).Error; err != nil {
			return err
		}
		// 默认绑定仪表盘菜单 (从数据库查询现有菜单)
		var dashboardMenu systemRbac.SysBaseMenu
		if err = tx.Where("path = ?", "dashboard").First(&dashboardMenu).Error; err == nil {
			// 找到则绑定
			auth.SysBaseMenus = []systemRbac.SysBaseMenu{dashboardMenu}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到则跳过 (可能没有初始化菜单)
			return nil
		} else {
			return err
		}
		if err = tx.Model(&auth).Association("SysBaseMenus").Replace(&auth.SysBaseMenus); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建角色失败")
	}

	var respParentId uint
	if auth.ParentId != nil {
		respParentId = *auth.ParentId
	}

	rs = res.CreateAuthorityRes{
		AuthorityId:   auth.AuthorityId,
		AuthorityName: auth.AuthorityName,
		ParentId:      respParentId,
		DefaultRouter: auth.DefaultRouter,
	}
	return rs, nil
}

// 复制角色-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) CopyAuthority(
	ctx *gin.Context,
	r req.CopyAuthorityReq,
) (rs res.CopyAuthorityRes, err error) {
	// 参数校验
	if r.Authority.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "新角色ID不能为空")
	}
	if r.OldAuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "旧角色ID不能为空")
	}

	// 检查新角色是否已存在
	var count int64
	global.GVA_DB.Model(&systemRbac.SysAuthority{}).Where("authority_id = ?", r.Authority.AuthorityId).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "新角色ID已存在")
	}

	// 查询旧角色
	var oldAuth systemRbac.SysAuthority
	err = global.GVA_DB.Preload("SysBaseMenus").Where("authority_id = ?", r.OldAuthorityId).First(&oldAuth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "旧角色不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 创建新角色
	newAuth := systemRbac.SysAuthority{
		AuthorityId:   r.Authority.AuthorityId,
		AuthorityName: r.Authority.AuthorityName,
		ParentId:      &r.Authority.ParentId,
		DefaultRouter: r.Authority.DefaultRouter,
	}

	// 复制菜单权限
	if len(oldAuth.SysBaseMenus) > 0 {
		newAuth.SysBaseMenus = oldAuth.SysBaseMenus
	}

	err = global.GVA_DB.Create(&newAuth).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "复制角色失败")
	}

	// 复制按钮权限
	var btns []systemRbac.SysAuthorityBtn
	global.GVA_DB.Find(&btns, "authority_id = ?", r.OldAuthorityId)
	if len(btns) > 0 {
		for i := range btns {
			btns[i].AuthorityId = r.Authority.AuthorityId
		}
		global.GVA_DB.Create(&btns)
	}

	rs = res.CopyAuthorityRes{
		AuthorityId:   r.Authority.AuthorityId,
		AuthorityName: newAuth.AuthorityName,
		ParentId:      r.Authority.ParentId,
		DefaultRouter: newAuth.DefaultRouter,
	}
	return rs, nil
}

// 更新角色-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) UpdateAuthority(
	ctx *gin.Context,
	r req.UpdateAuthorityReq,
) (rs res.UpdateAuthorityRes, err error) {
	// 参数校验
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", r.AuthorityId).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 更新角色
	if r.AuthorityName != "" {
		auth.AuthorityName = r.AuthorityName
	}
	if r.DefaultRouter != "" {
		auth.DefaultRouter = r.DefaultRouter
	}

	// 更新父角色ID (允许设置为 0，即根角色)
	var parentId *uint
	if r.ParentId != 0 {
		parentId = &r.ParentId
	}
	auth.ParentId = parentId

	err = global.GVA_DB.Save(&auth).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "更新角色失败")
	}

	var respParentId uint
	if auth.ParentId != nil {
		respParentId = *auth.ParentId
	}

	rs = res.UpdateAuthorityRes{
		AuthorityId:   r.AuthorityId,
		AuthorityName: auth.AuthorityName,
		ParentId:      respParentId,
		DefaultRouter: auth.DefaultRouter,
	}
	return rs, nil
}

// 删除角色-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) DeleteAuthority(
	ctx *gin.Context,
	r req.DeleteAuthorityReq,
) (err error) {
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 检查角色是否存在
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Preload("Users").Where("authority_id = ?", r.AuthorityId).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 检查是否有用户使用
	if len(auth.Users) > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "此角色有用户使用，无法删除")
	}

	// 检查是否有子角色
	var childCount int64
	global.GVA_DB.Model(&systemRbac.SysAuthority{}).Where("parent_id = ?", r.AuthorityId).Count(&childCount)
	if childCount > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "此角色存在子角色，无法删除")
	}

	// 事务删除
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Model(&auth).Association("SysBaseMenus").Clear(); err != nil {
			return err
		}
		// 删除数据权限关联
		if err := tx.Model(&auth).Association("DataAuthorityId").Clear(); err != nil {
			return err
		}
		// 删除用户角色关联
		if err := tx.Delete(&[]systemRbac.SysUserAuthority{}, "sys_authority_authority_id = ?", r.AuthorityId).Error; err != nil {
			return err
		}
		// 删除按钮权限
		if err := tx.Where("authority_id = ?", r.AuthorityId).Delete(&[]systemRbac.SysAuthorityBtn{}).Error; err != nil {
			return err
		}
		// 删除角色
		if err := tx.Delete(&auth).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除角色失败")
	}

	return nil
}

// 获取角色列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) GetAuthorityInfoList(
	ctx *gin.Context,
	r req.GetAuthorityInfoListReq,
) (rs res.GetAuthorityInfoListRes, err error) {

	var auth systemRbac.SysAuthority
	_ = global.GVA_DB.First(&auth, "authority_id = ?", r.AuthorityId)

	// 查询直接子角色
	var authorities []systemRbac.SysAuthority
	db := global.GVA_DB.Model(&systemRbac.SysAuthority{}).Preload("DataAuthorityId")

	if global.GVA_CONFIG.System.UseStrictAuth {
		// 严格树形结构
		if auth.ParentId == nil || *auth.ParentId == 0 {
			db = db.Where("authority_id = ? OR parent_id = ?", r.AuthorityId, r.AuthorityId)
		} else {
			db = db.Where("parent_id = ?", r.AuthorityId)
		}
	} else {
		db = db.Where("parent_id = ?", 0)
	}

	// 过滤掉无效的 authority_id = 0 记录
	db = db.Where("authority_id != ?", 0)

	err = db.Find(&authorities).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询角色列表失败")
	}

	// 递归查找子角色
	for i := range authorities {
		s.findChildrenAuthority(&authorities[i])
	}

	// 构建返回数据
	rs.List = make([]res.GetAuthorityInfoListResList, 0, len(authorities))
	for _, authority := range authorities {
		item := res.GetAuthorityInfoListResList{
			AuthorityId:   authority.AuthorityId,
			AuthorityName: authority.AuthorityName,
			ParentId:      *authority.ParentId,
			DefaultRouter: authority.DefaultRouter,
		}

		// 数据权限
		for _, da := range authority.DataAuthorityId {
			item.DataAuthorityId = append(item.DataAuthorityId, res.GetAuthorityInfoListResListDataauthorityid{
				AuthorityId:   da.AuthorityId,
				AuthorityName: da.AuthorityName,
			})
		}

		// 子角色
		for _, child := range authority.Children {
			item.Children = append(item.Children, res.GetAuthorityInfoListResListChildren{
				AuthorityId:   child.AuthorityId,
				AuthorityName: child.AuthorityName,
				ParentId:      *child.ParentId,
				DefaultRouter: child.DefaultRouter,
			})
		}

		rs.List = append(rs.List, item)
	}

	rs.AuthorityId = r.AuthorityId
	return rs, nil
}

// 递归查找子角色（带去重）
func (s *SysAuthorityService) findChildrenAuthority(auth *systemRbac.SysAuthority) error {
	return s.findChildrenAuthorityWithVisited(auth, make(map[uint]bool))
}

// 带已访问标记的递归查找子角色
func (s *SysAuthorityService) findChildrenAuthorityWithVisited(auth *systemRbac.SysAuthority, visited map[uint]bool) error {
	// 防止无限递归
	if visited[auth.AuthorityId] {
		return nil
	}
	visited[auth.AuthorityId] = true

	err := global.GVA_DB.Preload("DataAuthorityId").
		Where("parent_id = ?", auth.AuthorityId).Find(&auth.Children).Error
	if err != nil {
		return err
	}
	for i := range auth.Children {
		s.findChildrenAuthorityWithVisited(&auth.Children[i], visited)
	}
	return nil
}

// 获取角色结构列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) GetStructAuthorityList(
	ctx *gin.Context,
	r req.GetStructAuthorityListReq,
) (rs res.GetStructAuthorityListRes, err error) {
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	list, err := s.getAuthorityTreeList(r.AuthorityId)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询角色结构失败")
	}

	// 直接使用返回的列表
	rs.List = list
	return rs, nil
}

// 获取角色树形列表（带去重）
func (s *SysAuthorityService) getAuthorityTreeList(authorityId uint) ([]uint, error) {
	return s.getAuthorityTreeListWithVisited(authorityId, make(map[uint]bool))
}

// 带已访问标记的递归获取角色树形列表
func (s *SysAuthorityService) getAuthorityTreeListWithVisited(authorityId uint, visited map[uint]bool) ([]uint, error) {
	// 防止无限递归
	if visited[authorityId] {
		return nil, nil
	}
	visited[authorityId] = true

	var list []uint
	var auth systemRbac.SysAuthority
	err := global.GVA_DB.First(&auth, "authority_id = ?", authorityId).Error
	if err != nil {
		return nil, err
	}

	var authorities []systemRbac.SysAuthority
	err = global.GVA_DB.Preload("DataAuthorityId").Where("parent_id = ?", authorityId).Find(&authorities).Error
	if err != nil {
		return nil, err
	}

	for k := range authorities {
		list = append(list, authorities[k].AuthorityId)
		childrenList, err := s.getAuthorityTreeListWithVisited(authorities[k].AuthorityId, visited)
		if err == nil {
			list = append(list, childrenList...)
		}
	}

	if auth.ParentId == nil || *auth.ParentId == 0 {
		list = append(list, authorityId)
	}
	return list, nil
}

// 获取角色信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) GetAuthorityInfo(
	ctx *gin.Context,
	r req.GetAuthorityInfoReq,
) (rs res.GetAuthorityInfoRes, err error) {
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Preload("DataAuthorityId").Preload("SysBaseMenus").
		Where("authority_id = ?", r.AuthorityId).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 构建返回数据
	rs = res.GetAuthorityInfoRes{
		AuthorityId:   r.AuthorityId,
		AuthorityName: auth.AuthorityName,
		ParentId:      *auth.ParentId,
		DefaultRouter: auth.DefaultRouter,
	}

	// 数据权限
	for _, da := range auth.DataAuthorityId {
		rs.DataAuthorityId = append(rs.DataAuthorityId, res.GetAuthorityInfoResDataauthorityid{
			AuthorityId:   da.AuthorityId,
			AuthorityName: da.AuthorityName,
		})
	}

	// 菜单权限
	for _, menu := range auth.SysBaseMenus {
		rs.SysBaseMenus = append(rs.SysBaseMenus, res.GetAuthorityInfoResSysbasemenu{
			Id:        int64(menu.ID),
			ParentId:  menu.ParentId,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Meta: res.GetAuthorityInfoResSysbasemenuMeta{
				ActiveName:     menu.Meta.ActiveName,
				KeepAlive:      menu.Meta.KeepAlive,
				DefaultMenu:    menu.Meta.DefaultMenu,
				Title:          menu.Meta.Title,
				Icon:           menu.Meta.Icon,
				CloseTab:       menu.Meta.CloseTab,
				TransitionType: menu.Meta.TransitionType,
			},
		})
	}

	return rs, nil
}

// 设置角色数据权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) SetDataAuthority(
	ctx *gin.Context,
	r req.SetDataAuthorityReq,
) (err error) {
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Preload("DataAuthorityId").First(&auth, "authority_id = ?", r.AuthorityId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 解析数据权限ID
	var dataAuths []systemRbac.SysAuthority
	if r.DataAuthorityId != nil {
		daId := r.DataAuthorityId
		if err == nil {
			var da systemRbac.SysAuthority
			_ = global.GVA_DB.First(&da, "authority_id = ?", daId)
			dataAuths = append(dataAuths, da)
		}
	}

	err = global.GVA_DB.Model(&auth).Association("DataAuthorityId").Replace(&dataAuths)
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "设置数据权限失败")
	}

	return nil
}

// 设置角色菜单权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) SetMenuAuthority(
	ctx *gin.Context,
	r req.SetMenuAuthorityReq,
) (err error) {
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Preload("SysBaseMenus").First(&auth, "authority_id = ?", r.AuthorityId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	idsMap := make(map[uint]struct{})
	for _, m := range r.SysBaseMenus {
		if m.Id <= 0 {
			continue
		}
		idsMap[uint(m.Id)] = struct{}{}
	}

	ids := make([]uint, 0, len(idsMap))
	for id := range idsMap {
		ids = append(ids, id)
	}

	var menus []systemRbac.SysBaseMenu
	if len(ids) > 0 {
		if err := global.GVA_DB.Where("id IN ?", ids).Find(&menus).Error; err != nil {
			return biz_err.New(biz_err.DB_ERROR, "查询菜单失败")
		}
		if len(menus) != len(ids) {
			return biz_err.New(biz_err.PARAM_ERROR, "存在不存在的菜单ID")
		}
	}

	err = global.GVA_DB.Model(&auth).Association("SysBaseMenus").Replace(&menus)
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "设置菜单权限失败")
	}

	return nil
}

// 获取父角色ID-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityService) GetParentAuthorityID(
	ctx *gin.Context,
	r req.GetParentAuthorityIDReq,
) (rs res.GetParentAuthorityIDRes, err error) {
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", r.AuthorityId).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	rs.ParentId = *auth.ParentId
	return rs, nil
}
