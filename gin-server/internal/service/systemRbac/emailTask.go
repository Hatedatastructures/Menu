package systemRbac

import (
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmailTaskService struct{}

// 创建邮件任务-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) CreateEmailTask(
	c *gin.Context,
	r req.CreateEmailTaskReq,
) (err error) {
	task := systemRbac.EmailTask{
		Name:         r.Name,
		Description:  r.Description,
		CronExpr:     r.CronExpr,
		TemplateCode: r.TemplateCode,
		TargetType:   r.TargetType,
		Status:       r.Status,
	}

	err = global.GVA_DB.Create(&task).Error
	if err != nil {
		global.GVA_LOG.Error("创建邮件任务失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	return nil
}

// 更新邮件任务-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) UpdateEmailTask(
	c *gin.Context,
	r req.UpdateEmailTaskReq,
) (err error) {
	var task systemRbac.EmailTask
	err = global.GVA_DB.First(&task, r.Id).Error
	if err != nil {
		global.GVA_LOG.Error("任务不存在", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	err = global.GVA_DB.Model(&task).Updates(map[string]interface{}{
		"name":          r.Name,
		"description":   r.Description,
		"cron_expr":     r.CronExpr,
		"template_code": r.TemplateCode,
		"target_type":   r.TargetType,
	}).Error

	if err != nil {
		global.GVA_LOG.Error("更新邮件任务失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	return nil
}

// 删除邮件任务-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) DeleteEmailTask(
	c *gin.Context,
	r req.DeleteEmailTaskReq,
) (err error) {
	err = global.GVA_DB.Delete(&systemRbac.EmailTask{}, r.Id).Error
	if err != nil {
		global.GVA_LOG.Error("删除邮件任务失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}
	return nil
}

// 获取邮件任务列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) GetEmailTaskList(
	c *gin.Context,
	r req.GetEmailTaskListReq,
) (rs res.GetEmailTaskListRes, err error) {
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	db := global.GVA_DB.Model(&systemRbac.EmailTask{})

	if r.Status != "" {
		db = db.Where("status = ?", r.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var tasks []systemRbac.EmailTask
	offset := (page - 1) * pageSize
	err = db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&tasks).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	list := make([]res.GetEmailTaskListResList, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, res.GetEmailTaskListResList{
			Id:          int64(task.ID),
			Name:        task.Name,
			CronExpr:    task.CronExpr,
			Status:      task.Status,
			TotalRuns:   task.TotalRuns,
			SuccessRate: calculateSuccessRate(task.SuccessRuns, task.TotalRuns),
		})
	}

	rs = res.GetEmailTaskListRes{
		Page:     page,
		PageSize: pageSize,
		Status:   r.Status,
		Total:    total,
		List:     list,
	}

	return rs, nil
}

// 启动邮件任务-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) StartEmailTask(
	c *gin.Context,
	r req.StartEmailTaskReq,
) (err error) {
	err = global.GVA_DB.Model(&systemRbac.EmailTask{}).
		Where("id = ?", r.Id).
		Update("status", "active").Error
	if err != nil {
		global.GVA_LOG.Error("启动邮件任务失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}
	return nil
}

// 停止邮件任务-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) StopEmailTask(
	c *gin.Context,
	r req.StopEmailTaskReq,
) (err error) {
	err = global.GVA_DB.Model(&systemRbac.EmailTask{}).
		Where("id = ?", r.Id).
		Update("status", "inactive").Error
	if err != nil {
		global.GVA_LOG.Error("停止邮件任务失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}
	return nil
}

// 手动执行邮件任务-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTaskService) RunEmailTask(
	c *gin.Context,
	r req.RunEmailTaskReq,
) (rs res.RunEmailTaskRes, err error) {
	// TODO: 执行任务逻辑
	executionId := fmt.Sprintf("exec_%d", time.Now().Unix())

	rs = res.RunEmailTaskRes{
		ExecutionId: executionId,
	}

	return rs, nil
}

func calculateSuccessRate(success, total int) string {
	if total == 0 {
		return "0%"
	}
	rate := float64(success) / float64(total) * 100
	return fmt.Sprintf("%.1f%%", rate)
}

