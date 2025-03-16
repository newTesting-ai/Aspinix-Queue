package ds

import (
	"container/heap"
	"sync"
)

type Queue struct {
	id      int
	buffer  *RingBuffer
	size    int
	usage   float32
	enabled bool
	dlq     bool
	mu      sync.Mutex
}

type QueueHeap []*Queue

type Heap struct {
	heap  QueueHeap
	count int
}

func InitHeap() *Heap {
	buffer := InitRingBuffer(100)
	queue := Queue{
		id:      0,
		buffer:  buffer,
		size:    100,
		usage:   0.0,
		enabled: true,
		dlq:     false,
	}
	localHeap := &Heap{
		heap:  []*Queue{&queue},
		count: 0,
	}
	heap.Init(&localHeap.heap)
	return localHeap
}

func (qh QueueHeap) Len() int           { return len(qh) }
func (qh QueueHeap) Less(i, j int) bool { return qh[i].usage < qh[j].usage }
func (qh QueueHeap) Swap(i, j int) {
	qh[i], qh[j] = qh[j], qh[i]
	qh[i].id = i
	qh[j].id = j
}

func (q *Queue) addDataToSpecificQueue(data interface{}) (bool, float32) {
	q.mu.Lock()
	defer q.mu.Unlock()

	flag, count := q.buffer.InsertDataToRingBuffer(data)
	if flag {
		usage := float32(count) / float32(q.size)
		return true, usage
	}
	return false, 0.0
}

func (h *QueueHeap) Push(x interface{}) {
	n := len(*h)
	queue := x.(*Queue)
	queue.id = n
	*h = append(*h, queue)
}

func (h *QueueHeap) Pop() interface{} {
	old := *h
	n := len(old)
	queue := old[n-1]
	queue.id = -1
	*h = old[0 : n-1]
	return queue
}

func (h *Heap) addNewQueue() {
	n := len(h.heap)
	buffer := InitRingBuffer(100)
	queue := &Queue{
		id:      n,
		buffer:  buffer,
		size:    100,
		usage:   0.0,
		enabled: true,
		dlq:     false,
	}
	heap.Push(&h.heap, queue)
	h.count++
}
func (h *Heap) removeQueue() bool {
	if len(h.heap) == 0 {
		return false
	}
	heap.Pop(&h.heap)
	h.count--
	return true
}

func (q *Heap) updateQueueUsage(queue *Queue, newUsage float32) {
	heap.Remove(&q.heap, queue.id) // Remove from heap
	queue.usage = newUsage         // Update value
	heap.Push(&q.heap, queue)      // Reinsert into heap
}

func (h *Heap) AddData(data interface{}) bool {
	if h.heap[0].usage > 75.0 {
		h.addNewQueue()
	}
	flag, usage := h.heap[0].addDataToSpecificQueue(data)
	if flag {
		h.updateQueueUsage(h.heap[0], usage)
		return true
	}
	return false
}
