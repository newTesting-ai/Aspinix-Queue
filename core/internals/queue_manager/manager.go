package queuemanager

import (
	"aspinix-queue/core/ds"
	"sync"
)

type Queue struct {
	id      int
	buffer  *ds.RingBuffer
	size    int
	usage   float32
	enabled bool
	dlq     bool
	mu      sync.Mutex
}

type Queues struct {
	queue []Queue
	count int
}

func InitAspinixQueue() *Queues {
	buffer := ds.InitRingBuffer(100)
	queue := Queue{
		id:      0,
		buffer:  buffer,
		size:    100,
		usage:   0.0,
		enabled: true,
		dlq:     false,
	}
	return &Queues{
		queue: []Queue{queue},
		count: 1,
	}
}

func (q *Queues) GetUsages() ([]map[string]interface{}, float32) {
	var result []map[string]interface{}
	var usage float32
	var total float32

	for _, queue := range q.queue {
		queueMap := map[string]interface{}{
			"id":      queue.id, // Assuming buffer can be represented as an array
			"usage":   queue.usage,
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

func (q *Queue) addDataToSpecificQueue(data interface{}) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	flag, count := q.buffer.InsertDataToRingBuffer(data)
	if flag {
		q.usage = float32(count) / float32(q.size)
		return true
	}
	return false
}

func (q *Queues) AddData(data interface{}) bool {
	return true
}
