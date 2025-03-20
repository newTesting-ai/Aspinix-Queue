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
	heap          QueueHeap
	reserved      QueueHeap
	count         int
	reservedCount int
}

func InitHeap() *Heap {
	buffer := InitRingBuffer(100)
	queue := Queue{
		id:      0,
		buffer:  buffer,
		size:    4,
		usage:   0.0,
		enabled: true,
		dlq:     false,
	}
	reservedBuffer := InitRingBuffer(100)
	reservedQueue := Queue{
		id:      0,
		buffer:  reservedBuffer,
		size:    4,
		usage:   0.0,
		enabled: true,
		dlq:     false,
	}
	localHeap := &Heap{
		heap:          []*Queue{&queue},
		reserved:      []*Queue{&reservedQueue},
		count:         0,
		reservedCount: 0,
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
		usage := float32(count) * 100.0 / float32(q.size)
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

func (h *Heap) addNewQueue(priority bool) {
	if priority {
		n := len(h.reserved)
		buffer := InitRingBuffer(100)
		queue := &Queue{
			id:      n,
			buffer:  buffer,
			size:    4,
			usage:   0.0,
			enabled: true,
			dlq:     false,
		}
		heap.Push(&h.reserved, queue)
		h.reservedCount++
	} else {
		n := len(h.heap)
		buffer := InitRingBuffer(100)
		queue := &Queue{
			id:      n,
			buffer:  buffer,
			size:    4,
			usage:   0.0,
			enabled: true,
			dlq:     false,
		}
		heap.Push(&h.heap, queue)
		h.count++
	}
}
func (h *Heap) removeQueue() bool {
	if len(h.heap) == 0 {
		return false
	}
	heap.Pop(&h.heap)
	h.count--
	return true
}

func (q *Heap) updateQueueUsage(queue *Queue, newUsage float32, priority bool) {
	if priority {
		heap.Remove(&q.reserved, queue.id) // Remove from heap
		queue.usage = newUsage             // Update value
		heap.Push(&q.reserved, queue)      // Reinsert into heap
	} else {
		heap.Remove(&q.heap, queue.id) // Remove from heap
		queue.usage = newUsage         // Update value
		heap.Push(&q.heap, queue)      // Reinsert into heap
	}
}

func (h *Heap) handleDataAddition(heap *Queue, data interface{}, priority bool) bool {
	if heap.usage > 75.0 {
		h.addNewQueue(priority)
		return h.AddData(data, priority)
	}
	flag, usage := heap.addDataToSpecificQueue(data)
	if flag {
		h.updateQueueUsage(heap, usage, priority)
		return true
	}
	return false
}

func (h *Heap) AddData(data interface{}, priority bool) bool {
	var flag bool
	if priority {
		flag = h.handleDataAddition(h.reserved[0], data, priority)
	} else {
		flag = h.handleDataAddition(h.heap[0], data, priority)
	}
	return flag
}

func (q *Heap) GetUsages() ([]map[string]interface{}, float32) {
	var result []map[string]interface{}
	var usage float32
	var total float32

	for _, queue := range q.heap {
		queueMap := map[string]interface{}{
			"id":      queue.id, // Assuming buffer can be represented as an array
			"usage":   queue.usage,
			"size":    queue.size,
			"enabled": queue.enabled,
			"dlq":     queue.dlq,
		}
		result = append(result, queueMap)
		usage += queue.usage * float32(queue.size)
		total += float32(queue.size)
	}

	for _, queue := range q.reserved {
		queueMap := map[string]interface{}{
			"id":      queue.id, // Assuming buffer can be represented as an array
			"usage":   queue.usage,
			"size":    queue.size,
			"enabled": queue.enabled,
			"dlq":     queue.dlq,
		}
		result = append(result, queueMap)
		usage += queue.usage * float32(queue.size)
		total += float32(queue.size)
	}

	usage = usage / total

	return result, usage
}
