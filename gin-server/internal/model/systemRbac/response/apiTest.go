package response


import "time"

type SaveApiTestLogRes struct {
	Id uint `json:"id"` // 记录ID
}

type GetApiTestLogListRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Keyword string `json:"keyword"`
	List []GetApiTestLogListResList `json:"list"` // []
	Total uint `json:"total"` // 总数
}

type GetApiTestLogListResList struct {
	Id uint `json:"id"` // 记录ID
	Method string `json:"method"` // 请求方法
	Url string `json:"url"` // 请求URL
	StatusCode int `json:"statusCode"` // 响应状态码
	Duration int `json:"duration"` // 耗时(毫秒)
	Description string `json:"description"` // 描述
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

type GetApiTestLogDetailRes struct {
	Id uint `json:"id"` // 记录ID
	Method string `json:"method"` // 请求方法
	Url string `json:"url"` // 请求URL
	StatusCode int `json:"statusCode"` // 响应状态码
	Duration int `json:"duration"` // 耗时(毫秒)
	ReqHeaders string `json:"reqHeaders"` // 请求头JSON
	ReqBody string `json:"reqBody"` // 请求体
	ResHeaders string `json:"resHeaders"` // 响应头JSON
	ResBody string `json:"resBody"` // 响应体
	Params string `json:"params"` // 请求参数JSON
	Description string `json:"description"` // 描述
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}
