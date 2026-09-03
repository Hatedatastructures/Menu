package request

import (
	"shack/internal/model/common/request"
	"shack/internal/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
