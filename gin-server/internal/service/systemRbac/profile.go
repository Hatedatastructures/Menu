package systemRbac

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	"shack/internal/utils"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfileService struct{}

// 获取个人信息-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月08日 16:05:12
func (s *ProfileService) GetProfile(
	ctx context.Context,
) (rs res.GetProfileRes, err error) {
	// 从context获取gin.Context
	c, ok := ctx.(*gin.Context)
	if !ok {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
	}

	// 获取当前用户ID
	userId := utils.GetUserID(c)
	if userId == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询用户信息
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.USER_NOT_FOUND, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户信息失败")
	}

	// 构建返回数据
	rs = res.GetProfileRes{
		Uuid:        user.UUID.String(),
		Username:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		Phone:       user.Phone,
		Email:       user.Email,
		AuthorityId: user.AuthorityId,
		Enable:      user.Enable,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// 更新个人信息-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月08日 16:05:12
func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	r req.UpdateProfileReq,
) (rs res.UpdateProfileRes, err error) {
	// 从context获取gin.Context
	c, ok := ctx.(*gin.Context)
	if !ok {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
	}

	// 获取当前用户ID
	userId := utils.GetUserID(c)
	if userId == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.USER_NOT_FOUND, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 构建更新数据
	updates := make(map[string]interface{})

	if r.NickName != "" {
		updates["nick_name"] = r.NickName
	}
	if r.Phone != "" {
		updates["phone"] = r.Phone
	}
	if r.Email != "" {
		updates["email"] = r.Email
	}

	// 如果没有要更新的字段
	if len(updates) == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请提供要更新的信息")
	}

	// 更新用户信息
	err = global.GVA_DB.Model(&user).Updates(updates).Error
	if err != nil {
		global.GVA_LOG.Error("更新用户信息失败", zap.Error(err))
		return rs, biz_err.New(biz_err.USER_UPDATE_FAILED, "更新用户信息失败")
	}

	rs = res.UpdateProfileRes{
		Id: user.ID,
	}

	return rs, nil
}

// 修改密码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月08日 16:05:12
func (s *ProfileService) UpdatePassword(
	ctx context.Context,
	r req.UpdatePassword1Req,
) (err error) {
	// 参数校验
	if r.OldPassword == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "旧密码不能为空")
	}
	if r.NewPassword == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "新密码不能为空")
	}
	if r.ConfirmPassword == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "确认密码不能为空")
	}
	if r.NewPassword != r.ConfirmPassword {
		return biz_err.New(biz_err.PARAM_ERROR, "两次输入的密码不一致")
	}

	// 从context获取gin.Context
	c, ok := ctx.(*gin.Context)
	if !ok {
		return biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
	}

	// 获取当前用户ID
	userId := utils.GetUserID(c)
	if userId == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Select("id", "password").Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.USER_NOT_FOUND, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 验证旧密码
	if !utils.BcryptCheck(r.OldPassword, user.Password) {
		return biz_err.New(biz_err.INVALID_CREDENTIALS, "原密码错误")
	}

	// 加密新密码
	hashedPassword := utils.BcryptHash(r.NewPassword)

	// 更新密码
	err = global.GVA_DB.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		global.GVA_LOG.Error("更新密码失败", zap.Error(err))
		return biz_err.New(biz_err.PASSWORD_RESET_FAILED, "更新密码失败")
	}

	return nil
}

// 修改头像-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月08日 16:05:12
func (s *ProfileService) UpdateAvatar(
	ctx context.Context,
	r req.UpdateAvatarReq,
) (rs res.UpdateAvatarRes, err error) {
	// 参数校验
	if r.HeaderImg == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "头像URL不能为空")
	}

	// 从context获取gin.Context
	c, ok := ctx.(*gin.Context)
	if !ok {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
	}

	// 获取当前用户ID
	userId := utils.GetUserID(c)
	if userId == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.USER_NOT_FOUND, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 更新头像
	err = global.GVA_DB.Model(&user).Update("header_img", r.HeaderImg).Error
	if err != nil {
		global.GVA_LOG.Error("更新头像失败", zap.Error(err))
		return rs, biz_err.New(biz_err.USER_UPDATE_FAILED, "更新头像失败")
	}

	rs = res.UpdateAvatarRes{
		HeaderImg: r.HeaderImg,
	}

	return rs, nil
}

// 上传头像-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月08日 16:05:12
func (s *ProfileService) UploadAvatar(
	ctx context.Context,
	r req.UploadAvatarReq,
) (rs res.UploadAvatarRes, err error) {
	// 从context获取gin.Context
	c, ok := ctx.(*gin.Context)
	if !ok {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
	}

	// 获取当前用户ID
	userId := utils.GetUserID(c)
	if userId == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 从gin.Context获取文件
	file, err := c.FormFile("file")
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请选择要上传的文件")
	}

	// 验证文件类型（只允许图片）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "只支持上传图片文件（jpg、jpeg、png、gif、webp）")
	}

	// 验证文件大小（最大5MB）
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "图片大小不能超过5MB")
	}

	// 创建上传目录
	uploadDir := filepath.Join("statics", "avatars")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		global.GVA_LOG.Error("创建上传目录失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "创建上传目录失败")
	}

	// 生成唯一文件名
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(uploadDir, newFileName)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		global.GVA_LOG.Error("打开上传文件失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "打开上传文件失败")
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(savePath)
	if err != nil {
		global.GVA_LOG.Error("创建目标文件失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "创建目标文件失败")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		global.GVA_LOG.Error("保存文件失败", zap.Error(err))
		// 删除已创建的文件
		os.Remove(savePath)
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "保存文件失败")
	}

	// 构建访问URL
	fileURL := "/statics/avatars/" + newFileName

	// 更新用户头像
	var user systemRbac.User
	err = global.GVA_DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.USER_NOT_FOUND, "用户不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 删除旧头像文件（如果不是默认头像）
	if user.HeaderImg != "" && !strings.HasPrefix(user.HeaderImg, "http") {
		oldAvatarPath := filepath.Join("statics", strings.TrimPrefix(user.HeaderImg, "/statics/"))
		if err := os.Remove(oldAvatarPath); err != nil {
			// 只记录日志，不中断流程
			global.GVA_LOG.Warn("删除旧头像失败", zap.String("path", oldAvatarPath), zap.Error(err))
		}
	}

	// 更新数据库中的头像URL
	err = global.GVA_DB.Model(&user).Update("header_img", fileURL).Error
	if err != nil {
		global.GVA_LOG.Error("更新头像URL失败", zap.Error(err))
		// 删除已上传的文件
		os.Remove(savePath)
		return rs, biz_err.New(biz_err.DB_ERROR, "更新头像URL失败")
	}

	rs = res.UploadAvatarRes{
		Url:      fileURL,
		Filename: newFileName,
		Size:     file.Size,
	}

	return rs, nil
}
