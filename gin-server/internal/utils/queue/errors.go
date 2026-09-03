package queue

import "errors"

// 队列相关错误定义
var (
	// ErrQueueEmpty 队列为空错误
	ErrQueueEmpty = errors.New("queue is empty")

	// ErrQueueFull 队列已满错误（仅循环队列）
	ErrQueueFull = errors.New("queue is full")

	// ErrInvalidCapacity 无效的容量错误
	ErrInvalidCapacity = errors.New("invalid capacity")

	// ErrCapacityTooSmall 容量太小错误
	ErrCapacityTooSmall = errors.New("capacity is too small")
)
