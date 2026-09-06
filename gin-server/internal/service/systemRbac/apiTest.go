package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	biz_err "shack/internal/error"
)

type ApiTestService struct{}

// SaveApiTestLog 保存API测试请求记录-前台使用
func (s *ApiTestService) SaveApiTestLog(
	ctx *gin.Context,
	r req.SaveApiTestLogReq,
) (rs res.SaveApiTestLogRes, err error) {
	record := systemRbac.ApiTestLog{
		Method:      r.Method,
		URL:         r.Url,
		StatusCode:  r.StatusCode,
		Duration:    r.Duration,
		ReqHeaders:  r.ReqHeaders,
		ReqBody:     r.ReqBody,
		ResHeaders:  r.ResHeaders,
		ResBody:     r.ResBody,
		Params:      r.Params,
		Description: r.Description,
	}
	err = global.GVA_DB.Create(&record).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存测试记录失败")
	}
	rs = res.SaveApiTestLogRes{Id: record.ID}
	return rs, nil
}

// GetApiTestLogList 获取API测试记录列表-前台使用
func (s *ApiTestService) GetApiTestLogList(
	ctx *gin.Context,
	r req.GetApiTestLogListReq,
) (rs res.GetApiTestLogListRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 10
	}

	var total int64
	var logs []systemRbac.ApiTestLog

	db := global.GVA_DB.Model(&systemRbac.ApiTestLog{})

	if r.Keyword != "" {
		db = db.Where("url LIKE ? OR description LIKE ?", "%"+r.Keyword+"%", "%"+r.Keyword+"%")
	}

	db.Count(&total)
	err = db.Order("id desc").Limit(r.PageSize).Offset(r.PageSize * (r.Page - 1)).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询测试记录失败")
	}

	list := make([]res.GetApiTestLogListResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetApiTestLogListResList{
			Id:          log.ID,
			Method:      log.Method,
			Url:         log.URL,
			StatusCode:  log.StatusCode,
			Duration:    log.Duration,
			Description: log.Description,
			CreatedAt:   log.CreatedAt,
		})
	}

	rs = res.GetApiTestLogListRes{
		List:  list,
		Total: uint(total),
	}
	return rs, nil
}

// GetApiTestLogDetail 获取API测试记录详情-前台使用
func (s *ApiTestService) GetApiTestLogDetail(
	ctx *gin.Context,
	r req.GetApiTestLogDetailReq,
) (rs res.GetApiTestLogDetailRes, err error) {
	var record systemRbac.ApiTestLog
	err = global.GVA_DB.First(&record, r.Id).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "记录不存在")
	}

	rs = res.GetApiTestLogDetailRes{
		Id:          record.ID,
		Method:      record.Method,
		Url:         record.URL,
		StatusCode:  record.StatusCode,
		Duration:    record.Duration,
		ReqHeaders:  record.ReqHeaders,
		ReqBody:     record.ReqBody,
		ResHeaders:  record.ResHeaders,
		ResBody:     record.ResBody,
		Params:      record.Params,
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
	}
	return rs, nil
}

// DeleteApiTestLog 删除API测试记录-前台使用
func (s *ApiTestService) DeleteApiTestLog(
	ctx *gin.Context,
	r req.DeleteApiTestLogReq,
) (err error) {
	err = global.GVA_DB.Delete(&systemRbac.ApiTestLog{}, r.Id).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除测试记录失败")
	}
	return nil
}

// ClearApiTestLog 清空API测试记录-前台使用
func (s *ApiTestService) ClearApiTestLog(
	ctx *gin.Context,
) (err error) {
	err = global.GVA_DB.Unscoped().Delete(&[]systemRbac.ApiTestLog{}).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "清空测试记录失败")
	}
	return nil
}
