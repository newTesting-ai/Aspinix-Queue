package ds

import (
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

func (qh QueueHeap) Len() int           { return len(qh) }
func (qh QueueHeap) Less(i, j int) bool { return qh[i].usage < qh[j].usage }
func (qh QueueHeap) Swap(i, j int) {
	qh[i], qh[j] = qh[j], qh[i]
	qh[i].id = i
	qh[j].id = j
}

func NewQueue(id int) *Queue {
	buffer := InitRingBuffer(100)
	return &Queue{
		id:      id,
		buffer:  buffer,
		size:    4,
		usage:   0.0,
		enabled: true,
		dlq:     false,
	}
}

func Information(queue *Queue) map[string]interface{} {
	return map[string]interface{}{
		"id":      queue.id,
		"buffer":  queue.buffer,
		"size":    queue.size,
		"usage":   queue.usage,
		"enabled": queue.enabled,
		"dlq":     queue.dlq,
	}
}

// Getters
func (q *Queue) ID() int                        { return q.id }
func (q *Queue) Size() int                      { return q.size }
func (q *Queue) Usage() float32                 { return q.usage }
func (q *Queue) UpdateUsage(usage float32) bool { q.usage = usage; return true }
func (q *Queue) Enabled() bool                  { return q.enabled }
func (q *Queue) DLQ() bool                      { return q.dlq }

func (q *Queue) AddDataToSpecificQueue(data interface{}) (bool, float32) {
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
