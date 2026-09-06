package systemRbac

import (
	"errors"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	"shack/internal/utils"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 用户注册-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) Register(
	ctx *gin.Context,
	r req.RegisterReq,
) (rs res.RegisterRes, err error) {
	// 参数校验
	if r.UserName == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "用户名不能为空")
	}
	if r.PassWord == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "密码不能为空")
	}

	// 检查用户名是否已存在
	var count int64
	global.GVA_DB.Model(&systemRbac.User{}).Where("username = ?", r.UserName).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "用户名已存在")
	}

	r.AuthorityId = 888 // 默认角色
	if r.AuthorityId != 0 {
		r.AuthorityId = r.AuthorityId
	}

	// 创建用户
	user := systemRbac.User{
		Username:    r.UserName,
		Password:    utils.BcryptHash(r.PassWord),
		NickName:    r.NickName,
		HeaderImg:   r.HeaderImg,
		AuthorityId: r.AuthorityId,
		Enable:      r.Enable,
		Phone:       r.Phone,
		Email:       r.Email,
		UUID:        uuid.New(),
	}
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		roleIds := r.AuthorityIds
		if len(roleIds) == 0 {
			roleIds = []uint{user.AuthorityId}
		} else {
			user.AuthorityId = roleIds[0]
			if err := tx.Model(&user).Update("authority_id", user.AuthorityId).Error; err != nil {
				return err
			}
		}

		for _, roleId := range roleIds {
			ua := systemRbac.SysUserAuthority{
				SysUserId:               user.ID,
				SysAuthorityAuthorityId: roleId,
			}
			if err := tx.Create(&ua).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建用户失败")
	}

	rs = res.RegisterRes{
		Id:          int64(user.ID),
		UUID:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      int(user.Enable),
	}
	return rs, nil
}

// SysUserService 用户服务
type SysUserService struct{}

var SysUserServiceApp = new(SysUserService)

// 用户登录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) Login(
	ctx *gin.Context,
	r req.LoginReq,
) (rs res.LoginRes, err error) {
	// 参数校验
	if r.Username == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "用户名不能为空")
	}
	if r.Password == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "密码不能为空")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").
		Where("username = ?", r.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "用户名或密码错误")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 验证密码
	if !utils.BcryptCheck(r.Password, user.Password) {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "用户名或密码错误")
	}

	// 生成 token
	token, claims, err := utils.LoginToken(&user)
	if err != nil {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "生成token失败")
	}

	// 构建返回数据
	rs = res.LoginRes{
		Uuid:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
		Token:       token,
		ExpiresAt:   claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	}

	// 构建角色信息
	for _, auth := range user.Authorities {
		rs.Authorities = append(rs.Authorities, res.LoginResAuthority{
			AuthorityId:   auth.AuthorityId,
			AuthorityName: auth.AuthorityName,
			ParentId:      *auth.ParentId,
			DefaultRouter: auth.DefaultRouter,
		})
	}

	return rs, nil
}

// 修改用户密码-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) ChangePassword(
	ctx *gin.Context,
	r req.ChangePasswordReq,
) (err error) {
	if r.Password == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "原密码不能为空")
	}
	if r.NewPassword == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "新密码不能为空")
	}

	// 从上下文获取用户ID
	userId := utils.GetUserID(ctx)
	if userId == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Select("password").Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 验证原密码
	if !utils.BcryptCheck(r.Password, user.Password) {
		return biz_err.New(biz_err.PARAM_ERROR, "原密码错误")
	}

	// 更新密码
	err = global.GVA_DB.Model(&user).Update("password", utils.BcryptHash(r.NewPassword)).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新密码失败")
	}

	return nil
}

// 获取用户列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) GetUserInfoList(
	ctx *gin.Context,
	r req.GetUserInfoListReq,
) (rs res.GetUserInfoListRes, err error) {
	// 设置默认值
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.User{})
	if r.Username != "" {
		db = db.Where("username LIKE ?", "%"+r.Username+"%")
	}
	if r.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+r.NickName+"%")
	}
	if r.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+r.Phone+"%")
	}
	if r.Email != "" {
		db = db.Where("email LIKE ?", "%"+r.Email+"%")
	}

	// 查询总数
	var total int64
	db.Count(&total)

	// 分页查询
	var users []systemRbac.User
	offset := (r.Page - 1) * r.Size
	err = db.Offset(offset).Limit(r.Size).Preload("Authorities").Preload("Authority").
		Order("id desc").Find(&users).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户列表失败")
	}

	// 构建返回数据
	rs.List = make([]res.GetUserInfoListResList, 0, len(users))
	for _, user := range users {
		item := res.GetUserInfoListResList{
			Id:          int64(user.ID),
			Uuid:        user.UUID.String(),
			UserName:    user.Username,
			NickName:    user.NickName,
			HeaderImg:   user.HeaderImg,
			AuthorityId: user.AuthorityId,
			Phone:       user.Phone,
			Email:       user.Email,
			Enable:      user.Enable,
		}

		// 主角色信息
		if user.Authority.AuthorityId != 0 {
			item.Authority = res.GetUserInfoListResListAuthority{
				AuthorityId:   user.Authority.AuthorityId,
				AuthorityName: user.Authority.AuthorityName,
				DefaultRouter: user.Authority.DefaultRouter,
			}
		}

		roleSeen := make(map[uint]struct{})
		if user.Authority.AuthorityId != 0 {
			roleSeen[user.Authority.AuthorityId] = struct{}{}
			item.Authorities = append(item.Authorities, res.GetUserInfoListResListRole{
				AuthorityId:   user.Authority.AuthorityId,
				AuthorityName: user.Authority.AuthorityName,
			})
		}
		for _, auth := range user.Authorities {
			if _, ok := roleSeen[auth.AuthorityId]; ok {
				continue
			}
			roleSeen[auth.AuthorityId] = struct{}{}
			item.Authorities = append(item.Authorities, res.GetUserInfoListResListRole{
				AuthorityId:   auth.AuthorityId,
				AuthorityName: auth.AuthorityName,
			})
		}
		item.Roles = item.Authorities

		rs.List = append(rs.List, item)
	}

	rs.Page = r.Page
	rs.Size = r.Size
	rs.Total = total
	return rs, nil
}

// 设置用户角色(单一)-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) SetUserAuthority(
	ctx *gin.Context,
	r req.SetUserAuthorityReq,
) (err error) {
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "用户ID不能为空")
	}
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 检查用户是否存在
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", r.Id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 检查角色是否存在
	var auth systemRbac.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", r.AuthorityId).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "角色不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询角色失败")
	}

	// 检查用户是否绑定该角色
	var ua systemRbac.SysUserAuthority
	err = global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", r.Id, r.AuthorityId).First(&ua).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return biz_err.New(biz_err.DB_ERROR, "查询用户角色失败")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return biz_err.New(biz_err.PARAM_ERROR, "该用户无此角色")
	}

	// 更新用户主角色
	err = global.GVA_DB.Model(&user).Update("authority_id", r.AuthorityId).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "设置角色失败")
	}

	return nil
}

// 设置用户角色(多角色)-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) SetUserAuthorities(
	ctx *gin.Context,
	r req.SetUserAuthoritiesReq,
) (err error) {
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "用户ID不能为空")
	}
	if r.AuthorityIds == nil {
		return biz_err.New(biz_err.PARAM_MISSING, "角色ID不能为空")
	}

	// 检查用户是否存在
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", r.Id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 事务处理
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有角色关联
		err = tx.Delete(&[]systemRbac.SysUserAuthority{}, "sys_user_id = ?", r.Id).Error
		if err != nil {
			return err
		}

		// 添加新角色关联
		if len(r.AuthorityIds) > 0 {
			for _, authId := range r.AuthorityIds {
				ua := systemRbac.SysUserAuthority{
					SysUserId:               uint(r.Id),
					SysAuthorityAuthorityId: authId,
				}
				err = tx.Create(&ua).Error
				if err != nil {
					return err
				}
			}

			// 更新用户主角色（使用第一个角色作为主角色）
			err = tx.Model(&user).Update("authority_id", r.AuthorityIds[0]).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "设置角色失败")
	}

	return nil
}

// 删除用户-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) DeleteUser(
	ctx *gin.Context,
	r req.DeleteUserReq,
) (err error) {
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "用户ID不能为空")
	}

	// 检查用户是否存在
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", r.Id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 事务删除用户和角色关联
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]systemRbac.SysUserAuthority{}, "sys_user_id = ?", r.Id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除用户失败")
	}

	return nil
}

// 设置用户信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) SetUserInfo(
	ctx *gin.Context,
	r req.SetUserInfoReq,
) (err error) {
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "用户ID不能为空")
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"nick_name":  r.NickName,
		"header_img": r.HeaderImg,
		"phone":      r.Phone,
		"email":      r.Email,
		"enable":     r.Enable,
	}

	err = global.GVA_DB.Model(&systemRbac.User{}).Where("id = ?", r.Id).Updates(updates).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新用户信息失败")
	}

	return nil
}

// 设置自身信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) SetSelfInfo(
	ctx *gin.Context,
	r req.SetSelfInfoReq,
) (err error) {
	userId := utils.GetUserID(ctx)
	if userId == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 检查用户是否存在
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"nick_name":  r.NickName,
		"header_img": r.HeaderImg,
		"phone":      r.Phone,
		"email":      r.Email,
	}

	// 如果传了用户名则更新
	if r.UserName != "" {
		// 检查用户名是否已被占用
		var count int64
		global.GVA_DB.Model(&systemRbac.User{}).Where("username = ? AND id != ?", r.UserName, userId).Count(&count)
		if count > 0 {
			return biz_err.New(biz_err.PARAM_ERROR, "用户名已被占用")
		}
		updates["username"] = r.UserName
	}

	err = global.GVA_DB.Model(&user).Updates(updates).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新用户信息失败")
	}

	return nil
}

// 设置用户配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) SetSelfSetting(
	ctx *gin.Context,
	r req.SetSelfSettingReq,
) (err error) {
	userId := utils.GetUserID(ctx)
	if userId == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	err = global.GVA_DB.Model(&systemRbac.User{}).Where("id = ?", userId).
		Update("origin_setting", r.OriginSetting).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新用户配置失败")
	}

	return nil
}

// 获取用户信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) GetUserInf(
	ctx *gin.Context,
	r req.GetUserInfReq,
) (rs res.GetUserInfRes, err error) {
	if r.Uuid == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "UUID不能为空")
	}

	// 解析UUID
	userUUID, err := uuid.Parse(r.Uuid)
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_FORMAT, "UUID格式错误")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").
		Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 构建返回数据
	rs = res.GetUserInfRes{
		Id:          int64(user.ID),
		Uuid:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
	}

	// 主角色��息
	if user.Authority.AuthorityId != 0 {
		rs.Authority = res.GetUserInfResAuthority{
			AuthorityId:   user.Authority.AuthorityId,
			AuthorityName: user.Authority.AuthorityName,
			DefaultRouter: user.Authority.DefaultRouter,
		}
	}

	roleSeen := make(map[uint]struct{})
	if user.Authority.AuthorityId != 0 {
		roleSeen[user.Authority.AuthorityId] = struct{}{}
		role := res.GetUserInfResAuthori{
			AuthorityId:   user.Authority.AuthorityId,
			AuthorityName: user.Authority.AuthorityName,
		}
		rs.Authorities = append(rs.Authorities, role)
		rs.Authori = append(rs.Authori, role)
	}
	for _, auth := range user.Authorities {
		if _, ok := roleSeen[auth.AuthorityId]; ok {
			continue
		}
		roleSeen[auth.AuthorityId] = struct{}{}
		role := res.GetUserInfResAuthori{
			AuthorityId:   auth.AuthorityId,
			AuthorityName: auth.AuthorityName,
		}
		rs.Authorities = append(rs.Authorities, role)
		rs.Authori = append(rs.Authori, role)
	}

	return rs, nil
}

// 通过ID获取用户信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) FindUserById(
	ctx *gin.Context,
	r req.FindUserByIdReq,
) (rs res.FindUserByIdRes, err error) {
	if r.Id == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "用户ID不能为空")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", r.Id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	rs = res.FindUserByIdRes{
		Id:          int64(user.ID),
		Uuid:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
	}

	return rs, nil
}

// 通过UUID获取用户信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) FindUserByUuid(
	ctx *gin.Context,
	r req.FindUserByUuidReq,
) (rs res.FindUserByUuidRes, err error) {
	if r.Uuid == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "UUID不能为空")
	}

	// 解析UUID
	userUUID, err := uuid.Parse(r.Uuid)
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_FORMAT, "UUID格式错误")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	rs = res.FindUserByUuidRes{
		Id:          int64(user.ID),
		Uuid:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
	}

	return rs, nil
}

// 重置用户密码-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月30日 21:32:38
func (s *SysUserService) ResetPassword(
	ctx *gin.Context,
	r req.ResetPasswordReq,
) (err error) {
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "用户ID不能为空")
	}
	if r.Password == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "密码不能为空")
	}

	// 检查用户是否存在
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", r.Id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 更新密码
	err = global.GVA_DB.Model(&user).Update("password", utils.BcryptHash(r.Password)).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "重置密码失败")
	}

	return nil
}
