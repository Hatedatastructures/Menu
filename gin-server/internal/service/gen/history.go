package gen

import (
	"encoding/json"

	"shack/internal/global"
	"shack/internal/model/gen"
	req "shack/internal/model/gen/request"
	res "shack/internal/model/gen/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HistoryService struct{}

// GetHistoryList 获取生成历史列表
func (s *HistoryService) GetHistoryList(
	ctx *gin.Context,
	r req.GetHistoryListReq,
) (rs res.GetHistoryListRes, err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR)
	}

	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := global.GVA_DB.Model(&gen.GenHistory{}).Where("user_id = ?", userID)

	if r.Status != "" {
		db = db.Where("status = ?", r.Status)
	}

	var total int64
	db.Count(&total)

	var histories []gen.GenHistory
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询历史记录失败")
	}

	var list []res.GetHistoryListResList
	for _, h := range histories {
		list = append(list, res.GetHistoryListResList{
			Id:            h.ID,
			FileName:      h.FileName,
			FileId:        h.FileID,
			QuestionCount: h.QuestionCount,
			Status:        h.Status,
			CreatedAt:     h.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetHistoryListRes{
		Page:     page,
		PageSize: pageSize,
		Status:   r.Status,
		List:     list,
		Total:    total,
	}
	return rs, nil
}

// GetHistoryDetail 获取历史详情(含答案解析)
func (s *HistoryService) GetHistoryDetail(
	ctx *gin.Context,
	r req.GetHistoryDetailReq,
) (rs res.GetHistoryDetailRes, err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR)
	}

	var history gen.GenHistory
	err = global.GVA_DB.Where("id = ? AND user_id = ?", r.Id, userID).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "记录不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询记录失败")
	}

	rs = res.GetHistoryDetailRes{
		Id:            history.ID,
		FileName:      history.FileName,
		FileId:        history.FileID,
		QuestionCount: history.QuestionCount,
		Status:        history.Status,
		CreatedAt:     history.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 解析Result JSON
	if len(history.Result) > 0 {
		var result map[string]interface{}
		if err := json.Unmarshal(history.Result, &result); err == nil {
			rs.Result = result
		}
	}

	return rs, nil
}

// DeleteHistory 删除历史记录
func (s *HistoryService) DeleteHistory(
	ctx *gin.Context,
	r req.DeleteHistoryReq,
) (err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR)
	}

	result := global.GVA_DB.Where("id = ? AND user_id = ?", r.Id, userID).Delete(&gen.GenHistory{})
	if result.Error != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除记录失败")
	}
	if result.RowsAffected == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "记录不存在")
	}
	return nil
}
