package request

import (
	"shack/internal/model/common/request"
	"shack/internal/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
