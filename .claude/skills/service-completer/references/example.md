示例(具体情况具体定,实际开发时候,模块名可能不同):

package system

import (
	"shack/internal/internal/global"
	"shack/internal/internal/model/system"
	req "shack/internal/internal/model/system/request"
	"shack/internal/internal/model/system/response"
	"github.com/gin-gonic/gin"
  biz_err "shack/internal/internal/error"
)

type SysOperationRecordService struct {
}


func (s *SysOperationRecordService) DeleteSysOperationRecordByIds(c *gin.Context, r req.DeleteSysOperationRecordByIdsReq) (res response.DeleteSysOperationRecordByIdsRes, err error) {
	err = global.GVA_DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", r.Ids).Error
	res = response.DeleteSysOperationRecordByIdsRes{
		Success: err == nil,
	}
	return res, nil
}


func (s *SysOperationRecordService) DeleteSysOperationRecord(c *gin.Context, r req.DeleteSysOperationRecordReq) (res response.DeleteSysOperationRecordRes, err error) {
	var record system.SysOperationRecord
	record.ID = r.ID
	err = global.GVA_DB.Delete(&record).Error
	res = response.DeleteSysOperationRecordRes{
		Success: err == nil,
	}
	return res, nil
}