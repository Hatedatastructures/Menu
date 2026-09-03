package biz_err

// Err 业务错误结构体
type Err struct {
	Code int    `json:"code"` // 错误码
	Msg  string `json:"msg"`  // 错误信息
}

func (b *Err) Error() string {
	return b.Msg
}

func New(code int, msg ...string) *Err {
	message := ""

	if len(msg) <= 0 {
		message = GetMessage(code)
	} else {
		message = msg[0]
	}

	return &Err{
		Code: code,
		Msg:  message,
	}
}
