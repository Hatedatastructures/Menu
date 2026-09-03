package systemRbac

import (
	"errors"


	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	"github.com/gin-gonic/gin"
	biz_err "shack/internal/error"
	"gorm.io/gorm"
)

// SysAuthorityBtnService 角色按钮权限服务
type SysAuthorityBtnService struct{}

var SysAuthorityBtnServiceApp = new(SysAuthorityBtnService)

// 获取角色按钮权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityBtnService) GetAuthorityBtn(
	ctx *gin.Context,
	r req.GetAuthorityBtnReq,
) (rs res.GetAuthorityBtnRes ,err error) {
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}
	if r.MenuID == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "菜单ID不能为空")
	}



	// 查询按钮权限
	var btns []systemRbac.SysAuthorityBtn
	err = global.GVA_DB.Where("authority_id = ? AND sys_menu_id = ?", r.AuthorityId,r.MenuID).Find(&btns).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询按钮权限失败")
	}

	// 构建返回数据
	rs.Selected = make([]uint, 0, len(btns))
	for _, btn := range btns {
		rs.Selected = append(rs.Selected, btn.SysBaseMenuBtnID)
	}
	return rs, nil
}

// 设置角色按钮权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityBtnService) SetAuthorityBtn(
	ctx *gin.Context,
	r req.SetAuthorityBtnReq,
) (err error) {
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}
	if r.MenuID == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "菜单ID不能为空")
	}



	// 事务处理
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有按钮权限
		err = tx.Delete(&[]systemRbac.SysAuthorityBtn{}, "authority_id = ? AND sys_menu_id = ?", r.AuthorityId,r.MenuID).Error
		if err != nil {
			return err
		}

		// 添加新按钮权限
		if r.Selected != nil && len(r.Selected) > 0 {
			var btns []systemRbac.SysAuthorityBtn
			for _, btnId := range r.Selected {
				btns = append(btns, systemRbac.SysAuthorityBtn{
					AuthorityId:      r.AuthorityId,
					SysMenuID:      r.MenuID,
					SysBaseMenuBtnID: btnId,
				})
			}
			if len(btns) > 0 {
				err = tx.Create(&btns).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "设置按钮权限失败")
	}

	return nil
}

// 检查角色按钮是否可以删除-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysAuthorityBtnService) CanRemoveAuthorityBtn(
	ctx *gin.Context,
	r req.CanRemoveAuthorityBtnReq,
) (err error) {
	if r.Id == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "按钮ID不能为空")
	}

	// 检查按钮是否被使用
	var btn systemRbac.SysAuthorityBtn
	err = global.GVA_DB.First(&btn, "sys_base_menu_btn_id = ?", r.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 可以删除
		}
		return biz_err.New(biz_err.DB_ERROR, "检查按钮失败")
	}

	return biz_err.New(biz_err.PARAM_ERROR, "此按钮正在被使用，无法删除")
}