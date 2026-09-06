package initialize

import (
	"shack/internal/global"
	"shack/internal/service/systemRbac"

	"go.uber.org/zap"
)

// EmailQueue 初始化邮件队列
func EmailQueue() {
	// Worker数量，默认3个
	workerNum := 3

	// 初始化并启动邮件队列
	err := systemRbac.InitEmailQueue(workerNum)
	if err != nil {
		global.GVA_LOG.Error("初始化邮件队列失败", zap.Error(err))
		return
	}

	global.GVA_LOG.Info("邮件队列初始化成功", zap.Int("worker数量", workerNum))
}
