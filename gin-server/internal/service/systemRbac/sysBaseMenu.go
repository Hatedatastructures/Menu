package systemRbac

import (
	"errors"
	"strconv"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	"github.com/gin-gonic/gin"
	biz_err "shack/internal/error"
	"gorm.io/gorm"
)

// SysBaseMenuService 基础菜单服务
type SysBaseMenuService struct{}

var SysBaseMenuServiceApp = new(SysBaseMenuService)

// 删除基础菜单-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysBaseMenuService) DeleteBaseMenu(
	ctx *gin.Context,
	r req.DeleteBaseMenuReq,
) (err error) {
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "菜单ID不能为空")
	}

	// 检查是否有子菜单
	var childCount int64
	global.GVA_DB.Model(&systemRbac.SysBaseMenu{}).Where("parent_id = ?", r.Id).Count(&childCount)
	if childCount > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "此菜单存在子菜单，无法删除")
	}

	// 检查菜单是否存在
	var menu systemRbac.SysBaseMenu
	err = global.GVA_DB.First(&menu, r.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "菜单不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询菜单失败")
	}

	// 检查是否有角色使用此菜单作为首页
	var routerCount int64
	global.GVA_DB.Model(&systemRbac.SysAuthority{}).Where("default_router = ?", menu.Name).Count(&routerCount)
	if routerCount > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "此菜单有角色作为首页，无法删除")
	}

	// 事务删除
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&menu).Error; err != nil {
			return err
		}
		// 删除参数
		if err := tx.Where("sys_base_menu_id = ?", r.Id).Delete(&[]systemRbac.SysBaseMenuParameter{}).Error; err != nil {
			return err
		}
		// 删除按钮
		if err := tx.Where("sys_base_menu_id = ?", r.Id).Delete(&[]systemRbac.SysBaseMenuBtn{}).Error; err != nil {
			return err
		}
		// 删除角色菜单关联
		if err := tx.Where("sys_base_menu_id = ?", strconv.Itoa(r.Id)).Delete(&[]systemRbac.SysAuthorityMenu{}).Error; err != nil {
			return err
		}
		// 删除按钮权限
		if err := tx.Where("sys_menu_id = ?", r.Id).Delete(&[]systemRbac.SysAuthorityBtn{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除菜单失败")
	}

	return nil
}

// 更新基础菜单-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysBaseMenuService) UpdateBaseMenu(
	ctx *gin.Context,
	r req.UpdateBaseMenuReq,
) (rs res.UpdateBaseMenuRes, err error) {
	if r.Id == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "菜单ID不能为空")
	}

	// 查询菜单
	var menu systemRbac.SysBaseMenu
	err = global.GVA_DB.Where("id = ?", r.Id).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "菜单不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询菜单失败")
	}

	// 检查name是否重复
	if r.Name != "" && r.Name != menu.Name {
		var count int64
		global.GVA_DB.Model(&systemRbac.SysBaseMenu{}).Where("id <> ? AND name = ?", r.Id, r.Name).Count(&count)
		if count > 0 {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "菜单名称已存在")
		}
	}

	// 构建更新数据
	updates := map[string]interface{}{}
	if r.ParentId != 0 {

			updates["parent_id"] = r.ParentId

	}
	if r.Path != "" {
		updates["path"] = r.Path
	}
	if r.Name != "" {
		updates["name"] = r.Name
	}
	updates["hidden"] = r.Hidden
	if r.Component != "" {
		updates["component"] = r.Component
	}
	updates["sort"] = r.Sort
	updates["keep_alive"] = r.Meta.KeepAlive
	updates["transition_type"] = r.Meta.TransitionType
	updates["close_tab"] = r.Meta.CloseTab
	updates["default_menu"] = r.Meta.DefaultMenu
	updates["active_name"] = r.Meta.ActiveName
	updates["title"] = r.Meta.Title
	updates["icon"] = r.Meta.Icon

	err = global.GVA_DB.Model(&menu).Updates(updates).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "更新菜单失败")
	}

	rs.Id = r.Id
	return rs, nil
}

// 获取基础菜单详情-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysBaseMenuService) GetBaseMenuById(
	ctx *gin.Context,
	r req.GetBaseMenuByIdReq,
) (rs res.GetBaseMenuByIdRes, err error) {
	if r.Id == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "菜单ID不能为空")
	}

	// 查询菜单
	var menu systemRbac.SysBaseMenu
	err = global.GVA_DB.Preload("MenuBtn").Preload("Parameters").Where("id = ?", r.Id).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "菜单不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询菜单失败")
	}

	// 构建返回数据
	rs = res.GetBaseMenuByIdRes{
		Id:        int64(menu.ID),
		ParentId:  menu.ParentId,
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:   menu.Hidden,
		Component: menu.Component,
		Sort:     menu.Sort,
		Meta: res.GetBaseMenuByIdResMeta{
			ActiveName:     menu.Meta.ActiveName,
			KeepAlive:    menu.Meta.KeepAlive,
			DefaultMenu:   menu.Meta.DefaultMenu,
			Title:      menu.Meta.Title,
			Icon:      menu.Meta.Icon,
			CloseTab:   menu.Meta.CloseTab,
			TransitionType: menu.Meta.TransitionType,
		},
	}

	// 参数
	for _, p := range menu.Parameters {
		rs.Parameters = append(rs.Parameters, res.GetBaseMenuByIdResParameter{
			Id:    int64(p.ID),
			Type:  p.Type,
			Key:   p.Key,
			Value: p.Value,
		})
	}

	// 按钮
	for _, b := range menu.MenuBtn {
		rs.MenuBtn = append(rs.MenuBtn, res.GetBaseMenuByIdResMenubtn{
			Id:   int64(b.ID),
			Name: b.Name,
			Desc: b.Desc,
		})
	}

	return rs, nil
}