package queue

import (
	"testing"
)

// TestArrayQueue_BasicOperations 测试数组队列基本操作
func TestArrayQueue_BasicOperations(t *testing.T) {
	q := NewArrayQueue[int]()

	// 测试初始状态
	if !q.IsEmpty() {
		t.Error("新队列应该为空")
	}
	if q.Size() != 0 {
		t.Error("新队列大小应该为0")
	}

	// 测试入队
	err := q.Enqueue(1)
	if err != nil {
		t.Errorf("入队失败: %v", err)
	}
	err = q.Enqueue(2)
	err = q.Enqueue(3)

	if q.Size() != 3 {
		t.Errorf("队列大小应该为3，实际为%d", q.Size())
	}

	// 测试Peek
	val, err := q.Peek()
	if err != nil {
		t.Errorf("Peek失败: %v", err)
	}
	if val != 1 {
		t.Errorf("Peek应该返回1，实际返回%d", val)
	}
	if q.Size() != 3 {
		t.Error("Peek不应该改变队列大小")
	}

	// 测试出队
	val, err = q.Dequeue()
	if err != nil {
		t.Errorf("出队失败: %v", err)
	}
	if val != 1 {
		t.Errorf("第一次出队应该返回1，实际返回%d", val)
	}

	val, err = q.Dequeue()
	if err != nil {
		t.Errorf("出队失败: %v", err)
	}
	if val != 2 {
		t.Errorf("第二次出队应该返回2，实际返回%d", val)
	}

	if q.Size() != 1 {
		t.Errorf("队列大小应该为1，实际为%d", q.Size())
	}

	// 测试清空
	q.Clear()
	if !q.IsEmpty() {
		t.Error("清空后队列应该为空")
	}
}

// TestArrayQueue_EmptyQueue 测试空队列操作
func TestArrayQueue_EmptyQueue(t *testing.T) {
	q := NewArrayQueue[string]()

	// 测试空队列出队
	_, err := q.Dequeue()
	if err != ErrQueueEmpty {
		t.Errorf("空队列出队应该返回ErrQueueEmpty，实际返回: %v", err)
	}

	// 测试空队列Peek
	_, err = q.Peek()
	if err != ErrQueueEmpty {
		t.Errorf("空队列Peek应该返回ErrQueueEmpty，实际返回: %v", err)
	}
}

// TestArrayQueue_ToSlice 测试转换为切片
func TestArrayQueue_ToSlice(t *testing.T) {
	q := NewArrayQueue[int]()
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	slice := q.ToSlice()
	if len(slice) != 3 {
		t.Errorf("切片长度应该为3，实际为%d", len(slice))
	}
	if slice[0] != 10 || slice[1] != 20 || slice[2] != 30 {
		t.Errorf("切片内容不正确: %v", slice)
	}

	// 修改切片不应该影响队列
	slice[0] = 999
	val, _ := q.Peek()
	if val != 10 {
		t.Error("修改切片不应该影响队列")
	}
}

// TestCircularQueue_BasicOperations 测试循环队列基本操作
func TestCircularQueue_BasicOperations(t *testing.T) {
	q := NewCircularQueue[int](5)

	// 测试入队和出队
	for i := 1; i <= 5; i++ {
		err := q.Enqueue(i)
		if err != nil {
			t.Errorf("入队失败: %v", err)
		}
	}

	if q.Size() != 5 {
		t.Errorf("队列大小应该为5，实际为%d", q.Size())
	}

	// 测试满队列入队
	err := q.Enqueue(6)
	if err != ErrQueueFull {
		t.Errorf("满队列入队应该返回ErrQueueFull，实际返回: %v", err)
	}

	// 测试出队
	val, err := q.Dequeue()
	if err != nil {
		t.Errorf("出队失败: %v", err)
	}
	if val != 1 {
		t.Errorf("第一次出队应该返回1，实际返回%d", val)
	}

	// 出队后应该可以继续入队
	err = q.Enqueue(6)
	if err != nil {
		t.Errorf("出队后入队失败: %v", err)
	}

	// 验证队列顺序
	expected := []int{2, 3, 4, 5, 6}
	slice := q.ToSlice()
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("第%d个元素应该为%d，实际为%d", i, expected[i], v)
		}
	}
}

// TestCircularQueue_CircularBehavior 测试循环行为
func TestCircularQueue_CircularBehavior(t *testing.T) {
	q := NewCircularQueue[int](3)

	// 入队3个元素
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	// 出队2个元素
	q.Dequeue() // 1
	q.Dequeue() // 2

	// 再入队2个元素
	q.Enqueue(4)
	q.Enqueue(5)

	// 验证队列顺序
	slice := q.ToSlice()
	expected := []int{3, 4, 5}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("第%d个元素应该为%d，实际为%d", i, expected[i], v)
		}
	}
}

// TestCircularQueue_Resize 测试调整容量
func TestCircularQueue_Resize(t *testing.T) {
	q := NewCircularQueue[int](3)

	q.Enqueue(1)
	q.Enqueue(2)

	// 扩容
	err := q.Resize(5)
	if err != nil {
		t.Errorf("扩容失败: %v", err)
	}
	if q.Capacity() != 5 {
		t.Errorf("容量应该为5，实际为%d", q.Capacity())
	}

	// 缩容
	err = q.Resize(2)
	if err != ErrCapacityTooSmall {
		t.Errorf("缩容到小于当前大小应该返回ErrCapacityTooSmall，实际返回: %v", err)
	}

	// 正常缩容
	q.Dequeue() // 移除一个元素
	err = q.Resize(2)
	if err != nil {
		t.Errorf("缩容失败: %v", err)
	}
	if q.Capacity() != 2 {
		t.Errorf("容量应该为2，实际为%d", q.Capacity())
	}

	// 无效容量
	err = q.Resize(0)
	if err != ErrInvalidCapacity {
		t.Errorf("无效容量应该返回ErrInvalidCapacity，实际返回: %v", err)
	}
}

// TestPriorityQueue_BasicOperations 测试优先级队列基本操作（最大堆）
func TestPriorityQueue_BasicOperations(t *testing.T) {
	// 创建最大堆（数值越大优先级越高）
	q := NewPriorityQueue(func(a, b int) bool {
		return a > b
	})

	// 入队
	q.Enqueue(5)
	q.Enqueue(1)
	q.Enqueue(10)
	q.Enqueue(3)
	q.Enqueue(8)

	if q.Size() != 5 {
		t.Errorf("队列大小应该为5，实际为%d", q.Size())
	}

	// 出队应该按优先级从高到低
	expected := []int{10, 8, 5, 3, 1}
	for i, exp := range expected {
		val, err := q.Dequeue()
		if err != nil {
			t.Errorf("第%d次出队失败: %v", i+1, err)
		}
		if val != exp {
			t.Errorf("第%d次出队应该返回%d，实际返回%d", i+1, exp, val)
		}
	}

	// 队列应该为空
	if !q.IsEmpty() {
		t.Error("所有元素出队后队列应该为空")
	}
}

// TestPriorityQueue_WithPriority 测试带优先级的入队
func TestPriorityQueue_WithPriority(t *testing.T) {
	q := NewPriorityQueue(func(a, b string) bool {
		return a > b
	})

	// 相同值，不同优先级
	q.EnqueueWithPriority("low", 1)
	q.EnqueueWithPriority("high", 10)
	q.EnqueueWithPriority("medium", 5)

	// 应该按优先级出队
	val, _ := q.Dequeue()
	if val != "high" {
		t.Errorf("应该先出队high优先级，实际出队: %s", val)
	}

	val, _ = q.Dequeue()
	if val != "medium" {
		t.Errorf("应该次出队medium优先级，实际出队: %s", val)
	}

	val, _ = q.Dequeue()
	if val != "low" {
		t.Errorf("应该最后出队low优先级，实际出队: %s", val)
	}
}

// TestPriorityQueue_MinHeap 测试最小堆
func TestPriorityQueue_MinHeap(t *testing.T) {
	// 创建最小堆（数值越小优先级越高）
	q := NewPriorityQueue(func(a, b int) bool {
		return a < b
	})

	q.Enqueue(5)
	q.Enqueue(1)
	q.Enqueue(10)
	q.Enqueue(3)

	// 出队应该从小到大
	expected := []int{1, 3, 5, 10}
	for i, exp := range expected {
		val, err := q.Dequeue()
		if err != nil {
			t.Errorf("第%d次出队失败: %v", i+1, err)
		}
		if val != exp {
			t.Errorf("第%d次出队应该返回%d，实际返回%d", i+1, exp, val)
		}
	}
}

// TestPriorityQueue_Peek 测试优先级队列Peek
func TestPriorityQueue_Peek(t *testing.T) {
	q := NewPriorityQueue(func(a, b int) bool {
		return a > b
	})

	q.Enqueue(3)
	q.Enqueue(9)
	q.Enqueue(5)

	// Peek应该返回最高优先级的元素
	val, err := q.Peek()
	if err != nil {
		t.Errorf("Peek失败: %v", err)
	}
	if val != 9 {
		t.Errorf("Peek应该返回9，实际返回%d", val)
	}

	// Peek不应该移除元素
	if q.Size() != 3 {
		t.Error("Peek不应该改变队列大小")
	}
}

// TestPriorityQueue_StringTypes 测试字符串类型优先级队列
func TestPriorityQueue_StringTypes(t *testing.T) {
	q := NewPriorityQueue(func(a, b string) bool {
		return a < b // 字典序小的优先级高
	})

	names := []string{"charlie", "alice", "bob", "david"}
	for _, name := range names {
		q.Enqueue(name)
	}

	// 应该按字典序出队
	expected := []string{"alice", "bob", "charlie", "david"}
	for i, exp := range expected {
		val, err := q.Dequeue()
		if err != nil {
			t.Errorf("第%d次出队失败: %v", i+1, err)
		}
		if val != exp {
			t.Errorf("第%d次出队应该返回%s，实际返回%s", i+1, exp, val)
		}
	}
}

// TestQueue_Interface 测试接口一致性
func TestQueue_Interface(t *testing.T) {
	var q Queue[int]

	// 测试ArrayQueue实现了接口
	q = NewArrayQueue[int]()
	q.Enqueue(1)
	val, _ := q.Dequeue()
	if val != 1 {
		t.Error("ArrayQueue接口实现有问题")
	}

	// 测试CircularQueue实现了接口
	q = NewCircularQueue[int](5)
	q.Enqueue(2)
	val, _ = q.Dequeue()
	if val != 2 {
		t.Error("CircularQueue接口实现有问题")
	}
}

// BenchmarkArrayQueue_Enqueue 数组队列入队性能测试
func BenchmarkArrayQueue_Enqueue(b *testing.B) {
	q := NewArrayQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

// BenchmarkArrayQueue_Dequeue 数组队列出队性能测试
func BenchmarkArrayQueue_Dequeue(b *testing.B) {
	q := NewArrayQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

// BenchmarkCircularQueue_Enqueue 循环队列入队性能测试
func BenchmarkCircularQueue_Enqueue(b *testing.B) {
	q := NewCircularQueue[int](100000)
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

// BenchmarkCircularQueue_Dequeue 循环队列出队性能测试
func BenchmarkCircularQueue_Dequeue(b *testing.B) {
	q := NewCircularQueue[int](100000)
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

// BenchmarkPriorityQueue_Enqueue 优先级队列入队性能测试
func BenchmarkPriorityQueue_Enqueue(b *testing.B) {
	q := NewPriorityQueue(func(a, b int) bool {
		return a > b
	})
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

// BenchmarkPriorityQueue_Dequeue 优先级队列出队性能测试
func BenchmarkPriorityQueue_Dequeue(b *testing.B) {
	q := NewPriorityQueue(func(a, b int) bool {
		return a > b
	})
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

// 注意：以下队列实现不是并发安全的
// 如需并发安全的队列，请添加互斥锁保护
