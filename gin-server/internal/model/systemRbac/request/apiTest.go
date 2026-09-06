package request

type SaveApiTestLogReq struct {
	Method string `json:"method" form:"method"` // 请求方法
	Url string `json:"url" form:"url"` // 请求URL
	StatusCode int `json:"statusCode" form:"statusCode"` // 响应状态码
	Duration int `json:"duration" form:"duration"` // 耗时(毫秒)
	ReqHeaders string `json:"reqHeaders" form:"reqHeaders"` // 请求头JSON
	ReqBody string `json:"reqBody" form:"reqBody"` // 请求体
	ResHeaders string `json:"resHeaders" form:"resHeaders"` // 响应头JSON
	ResBody string `json:"resBody" form:"resBody"` // 响应体
	Params string `json:"params" form:"params"` // 请求参数JSON
	Description string `json:"description" form:"description"` // 描述
}

type GetApiTestLogListReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Keyword string `json:"keyword" form:"keyword"`
}

type GetApiTestLogDetailReq struct {
	Id int `json:"id" form:"id"`
}

type DeleteApiTestLogReq struct {
	Id int `json:"id" form:"id"`
}
