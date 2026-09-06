package systemRbac

import (
	"context"
	"fmt"
	"sync"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	"shack/internal/utils/queue"

	"go.uber.org/zap"
)

// EmailTask 邮件任务
type EmailTask struct {
	ID       uint
	To       []string
	Cc       []string
	Bcc      []string
	Subject  string
	Content  string
	IsHTML   bool
	Priority int
}

// EmailQueue 邮件队列处理器
type EmailQueue struct {
	queue      *queue.HeapPriorityQueue[*EmailTask]
	workerNum  int
	mailer     *Mailer
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.Mutex
	running    bool
}

// NewEmailQueue 创建邮件队列处理器
func NewEmailQueue(workerNum int) (*EmailQueue, error) {
	// 获取邮件配置
	config, err := GetEmailConfigFromDB()
	if err != nil {
		return nil, fmt.Errorf("获取邮件配置失败: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建优先级队列（优先级高的先处理）
	q := queue.NewPriorityQueue(func(a, b *EmailTask) bool {
		return a.Priority > b.Priority
	})

	eq := &EmailQueue{
		queue:     q,
		workerNum: workerNum,
		mailer:    NewMailer(*config),
		ctx:       ctx,
		cancel:    cancel,
		running:   false,
	}

	return eq, nil
}

// Start 启动队列处理器
func (eq *EmailQueue) Start() error {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if eq.running {
		return fmt.Errorf("邮件队列已经在运行")
	}

	eq.running = true

	// 启动worker协程
	for i := 0; i < eq.workerNum; i++ {
		eq.wg.Add(1)
		go eq.worker(i)
	}

	global.GVA_LOG.Info("邮件队列处理器已启动", zap.Int("worker数量", eq.workerNum))

	return nil
}

// Stop 停止队列处理器
func (eq *EmailQueue) Stop() {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if !eq.running {
		return
	}

	eq.running = false
	eq.cancel()

	// 等待所有worker完成
	eq.wg.Wait()

	global.GVA_LOG.Info("邮件队列处理器已停止")
}

// Push 推入邮件任务到队列
func (eq *EmailQueue) Push(task *EmailTask) error {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if !eq.running {
		return fmt.Errorf("邮件队列未运行")
	}

	err := eq.queue.Enqueue(task)
	if err != nil {
		return fmt.Errorf("推入队列失败: %w", err)
	}

	global.GVA_LOG.Info("邮件任务已推入队列",
		zap.Uint("id", task.ID),
		zap.Strings("to", task.To),
		zap.Int("priority", task.Priority),
		zap.Int("queueSize", eq.queue.Size()),
	)

	return nil
}

// worker worker协程处理邮件任务
func (eq *EmailQueue) worker(workerID int) {
	defer eq.wg.Done()

	global.GVA_LOG.Info("邮件worker已启动", zap.Int("workerID", workerID))

	for {
		select {
		case <-eq.ctx.Done():
			global.GVA_LOG.Info("邮件worker正在退出", zap.Int("workerID", workerID))
			return

		default:
			// 从队列中获取任务
			task, err := eq.queue.Dequeue()
			if err != nil {
				// 队列为空，稍后重试
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 处理邮件任务
			eq.processTask(workerID, task)
		}
	}
}

// processTask 处理邮件任务
func (eq *EmailQueue) processTask(workerID int, task *EmailTask) {
	global.GVA_LOG.Info("开始处理邮件任务",
		zap.Int("workerID", workerID),
		zap.Uint("id", task.ID),
		zap.Strings("to", task.To),
	)

	// 发送邮件
	err := eq.mailer.SendEmail(task.To, task.Cc, task.Bcc, task.Subject, task.Content, task.IsHTML)

	// 更新数据库状态
	var emailLog systemRbac.EmailLog
	err = global.GVA_DB.Where("id = ?", task.ID).First(&emailLog).Error
	if err == nil {
		if err != nil {
			global.GVA_DB.Model(&emailLog).Updates(map[string]interface{}{
				"status":    "failed",
				"error_msg": err.Error(),
			})
			global.GVA_LOG.Error("邮件发送失败",
				zap.Int("workerID", workerID),
				zap.Uint("id", task.ID),
				zap.Error(err),
			)
		} else {
			global.GVA_DB.Model(&emailLog).Update("status", "sent")
			global.GVA_LOG.Info("邮件发送成功",
				zap.Int("workerID", workerID),
				zap.Uint("id", task.ID),
				zap.Strings("to", task.To),
			)
		}
	}
}

// GetQueueSize 获取队列大小
func (eq *EmailQueue) GetQueueSize() int {
	return eq.queue.Size()
}

// IsRunning 检查队列是否在运行
func (eq *EmailQueue) IsRunning() bool {
	eq.mu.Lock()
	defer eq.mu.Unlock()
	return eq.running
}

// ReloadConfig 重新加载邮件配置
func (eq *EmailQueue) ReloadConfig() error {
	config, err := GetEmailConfigFromDB()
	if err != nil {
		return fmt.Errorf("获取邮件配置失败: %w", err)
	}

	eq.mailer = NewMailer(*config)
	global.GVA_LOG.Info("邮件配置已重新加载")

	return nil
}

// 全局邮件队列实例
var globalEmailQueue *EmailQueue

// InitEmailQueue 初始化邮件队列
func InitEmailQueue(workerNum int) error {
	var err error
	globalEmailQueue, err = NewEmailQueue(workerNum)
	if err != nil {
		return err
	}

	return globalEmailQueue.Start()
}

// StopEmailQueue 停止邮件队列
func StopEmailQueue() {
	if globalEmailQueue != nil {
		globalEmailQueue.Stop()
	}
}

// GetEmailQueue 获取全局邮件队列实例
func GetEmailQueue() *EmailQueue {
	return globalEmailQueue
}
