package queue

// Queue 队列接口 - 先进先出(FIFO)
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

// ArrayQueue 基于切片的数组队列实现
type ArrayQueue[T any] struct {
	items []T
}

// NewArrayQueue 创建新的数组队列
func NewArrayQueue[T any]() *ArrayQueue[T] {
	return &ArrayQueue[T]{
		items: make([]T, 0),
	}
}

// Enqueue 入队 - O(1)均摊时间复杂度
func (q *ArrayQueue[T]) Enqueue(item T) error {
	q.items = append(q.items, item)
	return nil
}

// Dequeue 出队 - O(n)时间复杂度（因为需要移动元素）
func (q *ArrayQueue[T]) Dequeue() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek 查看队列头部元素 - O(1)
func (q *ArrayQueue[T]) Peek() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}

	return q.items[0], nil
}

// IsEmpty 检查队列是否为空
func (q *ArrayQueue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size 获取队列大小
func (q *ArrayQueue[T]) Size() int {
	return len(q.items)
}

// Clear 清空队列
func (q *ArrayQueue[T]) Clear() {
	q.items = make([]T, 0)
}

// ToSlice 将队列转换为切片
func (q *ArrayQueue[T]) ToSlice() []T {
	result := make([]T, len(q.items))
	copy(result, q.items)
	return result
}

// CircularQueue 循环队列实现 - 优化出队操作
type CircularQueue[T any] struct {
	items     []T
	capacity  int
	head      int
	tail      int
	size      int
}

// NewCircularQueue 创建新的循环队列
func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	if capacity <= 0 {
		capacity = 10 // 默认容量
	}
	return &CircularQueue[T]{
		items:    make([]T, capacity),
		capacity: capacity,
		head:     0,
		tail:     0,
		size:     0,
	}
}

// Enqueue 入队 - O(1)时间复杂度
func (q *CircularQueue[T]) Enqueue(item T) error {
	if q.size == q.capacity {
		return ErrQueueFull
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return nil
}

// Dequeue 出队 - O(1)时间复杂度
func (q *CircularQueue[T]) Dequeue() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}

	item := q.items[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return item, nil
}

// Peek 查看队列头部元素 - O(1)
func (q *CircularQueue[T]) Peek() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}

	return q.items[q.head], nil
}

// IsEmpty 检查队列是否为空
func (q *CircularQueue[T]) IsEmpty() bool {
	return q.size == 0
}

// Size 获取队列大小
func (q *CircularQueue[T]) Size() int {
	return q.size
}

// Clear 清空队列
func (q *CircularQueue[T]) Clear() {
	q.head = 0
	q.tail = 0
	q.size = 0
}

// Capacity 获取队列容量
func (q *CircularQueue[T]) Capacity() int {
	return q.capacity
}

// ToSlice 将队列转换为切片
func (q *CircularQueue[T]) ToSlice() []T {
	result := make([]T, q.size)
	for i := 0; i < q.size; i++ {
		result[i] = q.items[(q.head+i)%q.capacity]
	}
	return result
}

// Resize 调整队列容量
func (q *CircularQueue[T]) Resize(newCapacity int) error {
	if newCapacity <= 0 {
		return ErrInvalidCapacity
	}

	if newCapacity <= q.size {
		return ErrCapacityTooSmall
	}

	newItems := make([]T, newCapacity)
	for i := 0; i < q.size; i++ {
		newItems[i] = q.items[(q.head+i)%q.capacity]
	}

	q.items = newItems
	q.capacity = newCapacity
	q.head = 0
	q.tail = q.size % newCapacity
	return nil
}

// PriorityQueue 优先级队列接口
type PriorityQueue[T any] interface {
	Queue[T]
	// EnqueueWithPriority 带优先级入队
	EnqueueWithPriority(item T, priority int) error
}

// HeapPriorityQueue 基于堆的优先级队列实现
type HeapPriorityQueue[T any] struct {
	items    []priorityItem[T]
	compare  func(a, b T) bool
}

type priorityItem[T any] struct {
	value    T
	priority int
}

// NewPriorityQueue 创建新的优先级队列
// compare: 返回true表示a的优先级高于b
func NewPriorityQueue[T any](compare func(a, b T) bool) *HeapPriorityQueue[T] {
	return &HeapPriorityQueue[T]{
		items:   make([]priorityItem[T], 0),
		compare: compare,
	}
}

// Enqueue 带默认优先级入队
func (q *HeapPriorityQueue[T]) Enqueue(item T) error {
	return q.EnqueueWithPriority(item, 0)
}

// EnqueueWithPriority 带优先级入队 - O(log n)
func (q *HeapPriorityQueue[T]) EnqueueWithPriority(item T, priority int) error {
	q.items = append(q.items, priorityItem[T]{
		value:    item,
		priority: priority,
	})
	q.heapifyUp(len(q.items) - 1)
	return nil
}

// Dequeue 出队 - O(log n)
func (q *HeapPriorityQueue[T]) Dequeue() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}

	root := q.items[0]
	last := len(q.items) - 1
	q.items[0] = q.items[last]
	q.items = q.items[:last]

	if len(q.items) > 0 {
		q.heapifyDown(0)
	}

	return root.value, nil
}

// Peek 查看队列头部元素 - O(1)
func (q *HeapPriorityQueue[T]) Peek() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrQueueEmpty
	}

	return q.items[0].value, nil
}

// IsEmpty 检查队列是否为空
func (q *HeapPriorityQueue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size 获取队列大小
func (q *HeapPriorityQueue[T]) Size() int {
	return len(q.items)
}

// Clear 清空队列
func (q *HeapPriorityQueue[T]) Clear() {
	q.items = make([]priorityItem[T], 0)
}

// ToSlice 将队列转换为切片
func (q *HeapPriorityQueue[T]) ToSlice() []T {
	result := make([]T, len(q.items))
	for i, item := range q.items {
		result[i] = item.value
	}
	return result
}

// heapifyUp 向上堆化
func (q *HeapPriorityQueue[T]) heapifyUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if !q.shouldSwap(parent, index) {
			break
		}
		q.items[parent], q.items[index] = q.items[index], q.items[parent]
		index = parent
	}
}

// heapifyDown 向下堆化
func (q *HeapPriorityQueue[T]) heapifyDown(index int) {
	last := len(q.items) - 1

	for {
		left := 2*index + 1
		right := 2*index + 2
		largest := index

		if left <= last && q.shouldSwap(largest, left) {
			largest = left
		}
		if right <= last && q.shouldSwap(largest, right) {
			largest = right
		}

		if largest == index {
			break
		}

		q.items[index], q.items[largest] = q.items[largest], q.items[index]
		index = largest
	}
}

// shouldSwap 判断是否需要交换
// 返回true表示需要交换（child应该上浮/parent应该下沉）
func (q *HeapPriorityQueue[T]) shouldSwap(parent, child int) bool {
	// 先比较优先级（数值越大优先级越高）
	if q.items[parent].priority != q.items[child].priority {
		// 如果child的优先级更高，需要交换
		return q.items[parent].priority < q.items[child].priority
	}
	// 优先级相同时，使用比较函数
	// compare返回true表示第一个参数优先级更高
	// 所以如果child的值优先级更高，需要交换
	return q.compare(q.items[child].value, q.items[parent].value)
}
