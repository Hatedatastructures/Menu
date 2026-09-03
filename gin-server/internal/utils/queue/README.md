# Queue 通用队列包

这是一个功能完善的Go语言队列实现包，提供了多种队列类型和接口设计。

## 特性

- ✅ **泛型支持** - 使用Go 1.18+泛型，支持任意数据类型
- ✅ **接口设计** - 统一的Queue接口，易于扩展和替换实现
- ✅ **多种实现** - 数组队列、循环队列、优先级队列
- ✅ **类型安全** - 编译时类型检查，避免运行时错误
- ✅ **高性能** - 循环队列入队出队均O(1)时间复杂度
- ✅ **完善测试** - 单元测试覆盖率100%，包含性能测试

## 队列类型

### 1. ArrayQueue（数组队列）

基于切片的简单队列实现，适合大多数场景。

**时间复杂度：**
- 入队：O(1) 均摊
- 出队：O(n) - 需要移动元素
- 查看头部：O(1)

**适用场景：**
- 队列大小适中
- 出队操作不频繁
- 需要动态扩容

```go
q := queue.NewArrayQueue[int]()
q.Enqueue(1)
q.Enqueue(2)
q.Enqueue(3)

val, _ := q.Dequeue() // 1
val, _ := q.Peek()    // 2
```

### 2. CircularQueue（循环队列）

固定容量的循环队列，性能最优。

**时间复杂度：**
- 入队：O(1)
- 出队：O(1)
- 查看头部：O(1)

**适用场景：**
- 需要高性能的入队出队
- 队列大小相对固定
- 内存使用敏感

```go
q := queue.NewCircularQueue[string](100) // 容量100
q.Enqueue("first")
q.Enqueue("second")

val, _ := q.Dequeue() // "first"

// 调整容量
q.Resize(200)
```

### 3. PriorityQueue（优先级队列）

基于堆的优先级队列，支持自定义优先级比较。

**时间复杂度：**
- 入队：O(log n)
- 出队：O(log n)
- 查看头部：O(1)

**适用场景：**
- 需要按优先级处理元素
- 任务调度
- 事件处理系统

```go
// 最大堆（数值越大优先级越高）
q := queue.NewPriorityQueue(func(a, b int) bool {
    return a > b
})

q.Enqueue(5)
q.Enqueue(10)
q.Enqueue(3)

val, _ := q.Dequeue() // 10

// 带优先级入队
q.EnqueueWithPriority("low priority task", 1)
q.EnqueueWithPriority("high priority task", 10)
```

## 接口定义

所有队列都实现了统一的`Queue[T]`接口：

```go
type Queue[T any] interface {
    // Enqueue 入队 - 将元素添加到队列尾部
    Enqueue(item T) error

    // Dequeue 出队 - 从队列头部移除并返回元素
    Dequeue() (T, error)

    // Peek 查看队列头部元素但不移除
    Peek() (T, error)

    // IsEmpty 检查队列是否为空
    IsEmpty() bool

    // Size 获取队列大小
    Size() int

    // Clear 清空队列
    Clear()

    // ToSlice 将队列转换为切片
    ToSlice() []T
}
```

## 使用示例

### 基本使用

```go
package main

import (
    "fmt"
    "shack/internal/utils/queue"
)

func main() {
    // 创建队列
    q := queue.NewArrayQueue[string]()

    // 入队
    q.Enqueue("task1")
    q.Enqueue("task2")
    q.Enqueue("task3")

    // 查看队列大小
    fmt.Printf("队列大小: %d\n", q.Size()) // 3

    // 查看头部元素
    head, _ := q.Peek()
    fmt.Printf("头部元素: %s\n", head) // "task1"

    // 出队
    for !q.IsEmpty() {
        item, _ := q.Dequeue()
        fmt.Printf("处理: %s\n", item)
    }
}
```

### 邮件任务队列

```go
// 定义邮件任务
type EmailTask struct {
    To      string
    Subject string
    Content string
}

// 创建任务队列
taskQueue := queue.NewCircularQueue[*EmailTask](1000)

// 添加任务
taskQueue.Enqueue(&EmailTask{
    To:      "user@example.com",
    Subject: "Welcome",
    Content: "Hello!",
})

// 工作协程处理任务
for !taskQueue.IsEmpty() {
    task, _ := taskQueue.Dequeue()
    sendEmail(task)
}
```

### 优先级任务调度

```go
type Job struct {
    Name     string
    Priority int
}

// 创建优先级队列（按优先级排序）
jobQueue := queue.NewPriorityQueue(func(a, b Job) bool {
    return a.Priority > b.Priority
})

jobQueue.Enqueue(Job{Name: "low", Priority: 1})
jobQueue.Enqueue(Job{Name: "high", Priority: 10})
jobQueue.Enqueue(Job{Name: "medium", Priority: 5})

// 会按优先级顺序处理：high -> medium -> low
for !jobQueue.IsEmpty() {
    job, _ := jobQueue.Dequeue()
    fmt.Printf("执行任务: %s (优先级: %d)\n", job.Name, job.Priority)
}
```

## 错误处理

队列可能返回以下错误：

```go
import "shack/internal/utils/queue"

// 队列为空时进行出队或查看操作
_, err := q.Dequeue()
if err == queue.ErrQueueEmpty {
    // 处理空队列情况
}

// 循环队列已满时入队
err := q.Enqueue(item)
if err == queue.ErrQueueFull {
    // 处理队列已满情况
}

// 无效的容量调整
err := q.Resize(-1)
if err == queue.ErrInvalidCapacity {
    // 处理无效容量
}
```

## 性能基准测试结果

在12th Gen Intel(R) Core(TM) i5-12500H上的测试结果：

```
ArrayQueue:
- Enqueue: 12.87 ns/op
- Dequeue:  1.08 ns/op

CircularQueue:
- Enqueue:  0.50 ns/op  (最快)
- Dequeue:  0.45 ns/op  (最快)

PriorityQueue:
- Enqueue: 454.40 ns/op
- Dequeue: 407.10 ns/op
```

## 注意事项

### 并发安全

⚠️ **当前实现不是并发安全的**

如果在多个goroutine中同时使用队列，需要添加额外的同步机制：

```go
import "sync"

type SafeQueue[T any] struct {
    mu sync.Mutex
    q  Queue[T]
}

func (sq *SafeQueue[T]) Enqueue(item T) error {
    sq.mu.Lock()
    defer sq.mu.Unlock()
    return sq.q.Enqueue(item)
}

// 类似地包装其他方法...
```

### 内存考虑

- **ArrayQueue**：动态扩容，但出队时可能有内存浪费
- **CircularQueue**：固定容量，内存使用高效
- **PriorityQueue**：基于堆，内存使用与元素数量成正比

## 运行测试

```bash
# 运行所有测试
go test ./internal/utils/queue/...

# 运行测试并显示覆盖率
go test -cover ./internal/utils/queue/...

# 运行性能测试
go test -bench=. -benchmem ./internal/utils/queue/...
```

## 许可证

本项目的一部分，遵循项目主许可证。

## 作者

- Auth: shack
- Github: https://github.com/hishack
