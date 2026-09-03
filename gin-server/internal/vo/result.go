package vo

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	bizErr "shack/internal/error"
)

type Result struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data"`
	RequestId string `json:"requestId"`
	TimeStamp int64  `json:"timeStamp"`
}

func Success(c *gin.Context, data any) Result {
	return Result{
		Code:      bizErr.SUCCESS,
		Msg:       "ok",
		Data:      data,
		RequestId: c.GetHeader("X-Request-ID"),
		TimeStamp: time.Now().Unix(),
	}
}

func Fail(c *gin.Context, data any, err error) Result {
	var newBizErr *bizErr.Err
	if ok := errors.As(err, &newBizErr); ok {
		return Result{
			Code:      newBizErr.Code,
			Msg:       newBizErr.Msg,
			Data:      data,
			RequestId: c.GetHeader("X-Request-ID"),
			TimeStamp: time.Now().Unix(),
		}
	}

	sysErr := bizErr.New(bizErr.SYSTEM_ERROR)
	return Result{
		Code:      sysErr.Code,
		Msg:       sysErr.Msg,
		Data:      data,
		RequestId: c.GetHeader("X-Request-ID"),
		TimeStamp: time.Now().Unix(),
	}
}
