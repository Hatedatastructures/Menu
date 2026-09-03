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

// SysMenuService 菜单服务
type SysMenuService struct{}

var SysMenuServiceApp = new(SysMenuService)

// 获取动态菜单树-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysMenuService) GetMenuTree(
	ctx *gin.Context,
	r req.GetMenuTreeReq,
) (rs res.GetMenuTreeRes, err error) {
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 获取角色菜单
	treeMap, err := s.getMenuTreeMap(r.AuthorityId)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取菜单失败")
	}

	// 构建树形结构
	menus := treeMap[0]
	for i := 0; i < len(menus); i++ {
		s.getMenuChildren(&menus[i], treeMap)
	}

	// 构建返回数据
	rs.List = make([]res.GetMenuTreeResList, 0, len(menus))
	for _, menu := range menus {
		rs.List = append(rs.List, s.buildMenuTreeRes(menu))
	}

	return rs, nil
}

// 获取菜单树Map
func (s *SysMenuService) getMenuTreeMap(authorityId uint) (map[uint][]systemRbac.SysMenu, error) {
	treeMap := make(map[uint][]systemRbac.SysMenu)

	// 查询角色菜单关联
	var authMenus []systemRbac.SysAuthorityMenu
	err := global.GVA_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&authMenus).Error
	if err != nil {
		return nil, err
	}

	var menuIds []string
	for _, am := range authMenus {
		menuIds = append(menuIds, am.MenuId)
	}

	if len(menuIds) == 0 {
		return treeMap, nil
	}

	// 查询菜单
	var baseMenus []systemRbac.SysBaseMenu
	err = global.GVA_DB.Where("id IN (?)", menuIds).Order("sort").
		Preload("Parameters").Find(&baseMenus).Error
	if err != nil {
		return nil, err
	}

	// 查询按钮权限
	var btns []systemRbac.SysAuthorityBtn
	global.GVA_DB.Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns)

	var btnMap = make(map[uint]map[string]uint)
	for _, btn := range btns {
		if btnMap[btn.SysMenuID] == nil {
			btnMap[btn.SysMenuID] = make(map[string]uint)
		}
		if btn.SysBaseMenuBtn.ID != 0 {
			btnMap[btn.SysMenuID][btn.SysBaseMenuBtn.Name] = btn.AuthorityId
		}
	}

	// 构建菜单
	var allMenus []systemRbac.SysMenu
	for _, bm := range baseMenus {
		menu := systemRbac.SysMenu{
			SysBaseMenu: bm,
			MenuId:      bm.ID,
		}
		menu.Parameters = bm.Parameters
		menu.Btns = btnMap[bm.ID]
		allMenus = append(allMenus, menu)
	}

	for _, m := range allMenus {
		treeMap[m.ParentId] = append(treeMap[m.ParentId], m)
	}

	return treeMap, nil
}

// 获取子菜单
func (s *SysMenuService) getMenuChildren(menu *systemRbac.SysMenu, treeMap map[uint][]systemRbac.SysMenu) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		s.getMenuChildren(&menu.Children[i], treeMap)
	}
}

// 构建菜单返回
func (s *SysMenuService) buildMenuTreeRes(menu systemRbac.SysMenu) res.GetMenuTreeResList {
	item := res.GetMenuTreeResList{
		MenuId:    menu.MenuId,
		ParentId:  menu.ParentId,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Meta: res.GetMenuTreeResListMeta{
			ActiveName:     menu.Meta.ActiveName,
			KeepAlive:      menu.Meta.KeepAlive,
			DefaultMenu:    menu.Meta.DefaultMenu,
			Title:          menu.Meta.Title,
			Icon:           menu.Meta.Icon,
			CloseTab:       menu.Meta.CloseTab,
			TransitionType: menu.Meta.TransitionType,
		},
	}

	// 子菜单
	for _, child := range menu.Children {
		item.Children = append(item.Children, s.buildMenuTreeResChildren(child))
	}

	// 参数
	for _, p := range menu.Parameters {
		item.Parameters = append(item.Parameters, res.GetMenuTreeResListParameter{
			Type:  p.Type,
			Key:   p.Key,
			Value: p.Value,
		})
	}

	return item
}

// 构建子菜单返回(递归)
func (s *SysMenuService) buildMenuTreeResChildren(menu systemRbac.SysMenu) res.GetMenuTreeResListChildren {
	item := res.GetMenuTreeResListChildren{
		MenuId:    menu.MenuId,
		ParentId:  menu.ParentId,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Meta: res.GetMenuTreeResListChildrenMeta{
			ActiveName:     menu.Meta.ActiveName,
			KeepAlive:      menu.Meta.KeepAlive,
			DefaultMenu:    menu.Meta.DefaultMenu,
			Title:          menu.Meta.Title,
			Icon:           menu.Meta.Icon,
			CloseTab:       menu.Meta.CloseTab,
			TransitionType: menu.Meta.TransitionType,
		},
	}

	// 递归子菜单
	for _, child := range menu.Children {
		item.Children = append(item.Children, s.buildMenuTreeResChildren(child))
	}

	return item
}

// 获取基础菜单列表(分页)-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysMenuService) GetInfoList(
	ctx *gin.Context,
	r req.GetInfoListReq,
) (rs res.GetInfoListRes, err error) {

	// 获取菜单树
	treeMap, err := s.getBaseMenuTreeMap(r.AuthorityId)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取菜单失败")
	}

	menus := treeMap[0]
	for i := 0; i < len(menus); i++ {
		s.getBaseChildren(&menus[i], treeMap)
	}

	// 构建返回数据
	rs.List = make([]res.GetInfoListResList, 0, len(menus))
	for _, menu := range menus {
		item := res.GetInfoListResList{
			Id:        int64(menu.ID),
			ParentId:  menu.ParentId,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Meta: res.GetInfoListResListMeta{
				ActiveName:     menu.Meta.ActiveName,
				KeepAlive:      menu.Meta.KeepAlive,
				DefaultMenu:    menu.Meta.DefaultMenu,
				Title:          menu.Meta.Title,
				Icon:           menu.Meta.Icon,
				CloseTab:       menu.Meta.CloseTab,
				TransitionType: menu.Meta.TransitionType,
			},
		}

		// 递归子菜单
		for _, child := range menu.Children {
			item.Children = append(item.Children, s.buildInfoListResChildren(child))
		}

		// 参数
		for _, p := range menu.Parameters {
			item.Parameters = append(item.Parameters, res.GetInfoListResListParameter{
				Type:  p.Type,
				Key:   p.Key,
				Value: p.Value,
			})
		}

		// 按钮
		for _, b := range menu.MenuBtn {
			item.MenuBtn = append(item.MenuBtn, res.GetInfoListResListMenubtn{
				Id:   int64(b.ID),
				Name: b.Name,
				Desc: b.Desc,
			})
		}

		rs.List = append(rs.List, item)
	}

	rs.Total = int64(len(rs.List))
	return rs, nil
}

// 构建子菜单返回(递归)-InfoList用
func (s *SysMenuService) buildInfoListResChildren(menu systemRbac.SysBaseMenu) res.GetInfoListResList {
	item := res.GetInfoListResList{
		Id:        int64(menu.ID),
		ParentId:  menu.ParentId,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Meta: res.GetInfoListResListMeta{
			ActiveName:     menu.Meta.ActiveName,
			KeepAlive:      menu.Meta.KeepAlive,
			DefaultMenu:    menu.Meta.DefaultMenu,
			Title:          menu.Meta.Title,
			Icon:           menu.Meta.Icon,
			CloseTab:       menu.Meta.CloseTab,
			TransitionType: menu.Meta.TransitionType,
		},
	}

	// 递归子菜单
	for _, child := range menu.Children {
		item.Children = append(item.Children, s.buildInfoListResChildren(child))
	}

	// 参数
	for _, p := range menu.Parameters {
		item.Parameters = append(item.Parameters, res.GetInfoListResListParameter{
			Type:  p.Type,
			Key:   p.Key,
			Value: p.Value,
		})
	}

	// 按钮
	for _, b := range menu.MenuBtn {
		item.MenuBtn = append(item.MenuBtn, res.GetInfoListResListMenubtn{
			Id:   int64(b.ID),
			Name: b.Name,
			Desc: b.Desc,
		})
	}

	return item
}

// 获取基础菜单树Map
func (s *SysMenuService) getBaseMenuTreeMap(authorityId uint) (map[uint][]systemRbac.SysBaseMenu, error) {
	treeMap := make(map[uint][]systemRbac.SysBaseMenu)

	var menus []systemRbac.SysBaseMenu
	err := global.GVA_DB.Order("sort").Preload("MenuBtn").Preload("Parameters").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	for _, m := range menus {
		treeMap[m.ParentId] = append(treeMap[m.ParentId], m)
	}

	return treeMap, nil
}

// 获取子菜单(基础)
func (s *SysMenuService) getBaseChildren(menu *systemRbac.SysBaseMenu, treeMap map[uint][]systemRbac.SysBaseMenu) {
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		s.getBaseChildren(&menu.Children[i], treeMap)
	}
}

// 添加基础菜单-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysMenuService) AddBaseMenu(
	ctx *gin.Context,
	r req.AddBaseMenuReq,
) (rs res.AddBaseMenuRes, err error) {
	// 参数校验
	if r.Name == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "菜单名称不能为空")
	}
	if r.Path == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "菜单路径不能为空")
	}

	// 检查name是否重复
	var count int64
	global.GVA_DB.Model(&systemRbac.SysBaseMenu{}).Where("name = ?", r.Name).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "菜单名称已存在")
	}

	// 创建菜单
	menu := systemRbac.SysBaseMenu{
		ParentId:  r.ParentId,
		Path:      r.Path,
		Name:      r.Name,
		Hidden:    r.Hidden,
		Component: r.Component,
		Sort:      r.Sort,
		Meta: systemRbac.Meta{
			ActiveName:     r.Meta.ActiveName,
			KeepAlive:      r.Meta.KeepAlive,
			DefaultMenu:    r.Meta.DefaultMenu,
			Title:          r.Meta.Title,
			Icon:           r.Meta.Icon,
			CloseTab:       r.Meta.CloseTab,
			TransitionType: r.Meta.TransitionType,
		},
	}

	err = global.GVA_DB.Create(&menu).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建菜单失败")
	}

	rs = res.AddBaseMenuRes{
		Id:        int64(menu.ID),
		ParentId:  r.ParentId,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
	}
	return rs, nil
}

// 获取基础菜单树-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysMenuService) GetBaseMenuTree(
	ctx *gin.Context,
	r req.GetBaseMenuTreeReq,
) (rs res.GetBaseMenuTreeRes, err error) {

	treeMap, err := s.getBaseMenuTreeMap(r.AuthorityId)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取菜单失败")
	}

	menus := treeMap[0]
	for i := 0; i < len(menus); i++ {
		s.getBaseChildren(&menus[i], treeMap)
	}

	rs.List = make([]res.GetBaseMenuTreeResList, 0, len(menus))
	for _, menu := range menus {
		item := res.GetBaseMenuTreeResList{
			Id:        int64(menu.ID),
			ParentId:  menu.ParentId,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Meta: res.GetBaseMenuTreeResListMeta{
				ActiveName:     menu.Meta.ActiveName,
				KeepAlive:      menu.Meta.KeepAlive,
				DefaultMenu:    menu.Meta.DefaultMenu,
				Title:          menu.Meta.Title,
				Icon:           menu.Meta.Icon,
				CloseTab:       menu.Meta.CloseTab,
				TransitionType: menu.Meta.TransitionType,
			},
		}

		// 递归子菜单
		for _, child := range menu.Children {
			item.Children = append(item.Children, s.buildBaseMenuTreeResChildren(child))
		}

		for _, p := range menu.Parameters {
			item.Parameters = append(item.Parameters, res.GetBaseMenuTreeResListParameter{
				Type:  p.Type,
				Key:   p.Key,
				Value: p.Value,
			})
		}

		for _, b := range menu.MenuBtn {
			item.MenuBtn = append(item.MenuBtn, res.GetBaseMenuTreeResListMenubtn{
				Id:   int64(b.ID),
				Name: b.Name,
				Desc: b.Desc,
			})
		}

		rs.List = append(rs.List, item)
	}

	return rs, nil
}

// 构建子菜单返回(递归)-GetBaseMenuTree用
func (s *SysMenuService) buildBaseMenuTreeResChildren(menu systemRbac.SysBaseMenu) res.GetBaseMenuTreeResList {
	item := res.GetBaseMenuTreeResList{
		Id:        int64(menu.ID),
		ParentId:  menu.ParentId,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Meta: res.GetBaseMenuTreeResListMeta{
			ActiveName:     menu.Meta.ActiveName,
			KeepAlive:      menu.Meta.KeepAlive,
			DefaultMenu:    menu.Meta.DefaultMenu,
			Title:          menu.Meta.Title,
			Icon:           menu.Meta.Icon,
			CloseTab:       menu.Meta.CloseTab,
			TransitionType: menu.Meta.TransitionType,
		},
	}

	// 递归子菜单
	for _, child := range menu.Children {
		item.Children = append(item.Children, s.buildBaseMenuTreeResChildren(child))
	}

	// 参数
	for _, p := range menu.Parameters {
		item.Parameters = append(item.Parameters, res.GetBaseMenuTreeResListParameter{
			Type:  p.Type,
			Key:   p.Key,
			Value: p.Value,
		})
	}

	// 按钮
	for _, b := range menu.MenuBtn {
		item.MenuBtn = append(item.MenuBtn, res.GetBaseMenuTreeResListMenubtn{
			Id:   int64(b.ID),
			Name: b.Name,
			Desc: b.Desc,
		})
	}

	return item
}

// 为角色分配菜单-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysMenuService) AddMenuAuthority(
	ctx *gin.Context,
	r req.AddMenuAuthorityReq,
) (err error) {
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.First(&auth, "authority_id = ?", r.AuthorityId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 转换菜单
	var menus []systemRbac.SysBaseMenu
	for _, m := range r.Menus {
		menus = append(menus, systemRbac.SysBaseMenu{
			GVA_MODEL: global.GVA_MODEL{ID: uint(m.Id)},
			ParentId:  m.ParentId,
			Path:      m.Path,
			Name:      m.Name,
			Hidden:    m.Hidden,
			Component: m.Component,
			Sort:      m.Sort,
		})
	}

	err = global.GVA_DB.Model(&auth).Association("SysBaseMenus").Replace(&menus)
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "分配菜单失败")
	}

	return nil
}

// 获取角色菜单权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysMenuService) GetMenuAuthority(
	ctx *gin.Context,
	r req.GetMenuAuthorityReq,
) (rs res.GetMenuAuthorityRes, err error) {
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 查询角色菜单关联
	var authMenus []systemRbac.SysAuthorityMenu
	err = global.GVA_DB.Where("sys_authority_authority_id = ?", r.AuthorityId).Find(&authMenus).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询菜单失败")
	}

	var menuIds []string
	for _, am := range authMenus {
		menuIds = append(menuIds, am.MenuId)
	}

	if len(menuIds) == 0 {
		return rs, nil
	}

	// 查询菜单
	var baseMenus []systemRbac.SysBaseMenu
	err = global.GVA_DB.Where("id IN (?)", menuIds).Order("sort").Find(&baseMenus).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询菜单失败")
	}

	rs.List = make([]res.GetMenuAuthorityResList, 0, len(baseMenus))
	for _, menu := range baseMenus {
		item := res.GetMenuAuthorityResList{
			MenuId:    menu.ID,
			ParentId:  menu.ParentId,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Meta: res.GetMenuAuthorityResListMeta{
				ActiveName:     menu.Meta.ActiveName,
				KeepAlive:      menu.Meta.KeepAlive,
				DefaultMenu:    menu.Meta.DefaultMenu,
				Title:          menu.Meta.Title,
				Icon:           menu.Meta.Icon,
				CloseTab:       menu.Meta.CloseTab,
				TransitionType: menu.Meta.TransitionType,
			},
		}

		for _, p := range menu.Parameters {
			item.Parameters = append(item.Parameters, res.GetMenuAuthorityResListParameter{
				Type:  p.Type,
				Key:   p.Key,
				Value: p.Value,
			})
		}

		rs.List = append(rs.List, item)
	}

	return rs, nil
}
