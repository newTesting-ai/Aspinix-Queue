package queuemanager

import (
	"aspinix-queue/core/ds"
	"container/heap"
)

type Heap struct {
	heap          ds.QueueHeap
	reserved      ds.QueueHeap
	count         int
	reservedCount int
}

func InitHeap() *Heap {
	queue := ds.NewQueue(0)
	reservedQueue := ds.NewQueue(0)

	localHeap := &Heap{
		heap:          []*ds.Queue{queue},
		reserved:      []*ds.Queue{reservedQueue},
		count:         0,
		reservedCount: 0,
	}
	heap.Init(&localHeap.heap)
	return localHeap
}

func (h *Heap) addNewQueue(priority bool) {
	if priority {
		n := len(h.reserved)
		queue := ds.NewQueue(n)
		heap.Push(&h.reserved, queue)
		h.reservedCount++
	} else {
		n := len(h.heap)
		queue := ds.NewQueue(n)
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

func (q *Heap) updateQueueUsage(queue *ds.Queue, newUsage float32, priority bool) {
	if priority {
		heap.Remove(&q.reserved, queue.ID()) // Remove from heap
		queue.UpdateUsage(newUsage)          // Update value
		heap.Push(&q.reserved, queue)        // Reinsert into heap
	} else {
		heap.Remove(&q.heap, queue.ID()) // Remove from heap
		queue.UpdateUsage(newUsage)      // Update value
		heap.Push(&q.heap, queue)        // Reinsert into heap
	}
}

func (h *Heap) handleDataAddition(heap *ds.Queue, data interface{}, priority bool) bool {
	if heap.Usage() > 75.0 {
		h.addNewQueue(priority)
		return h.AddData(data, priority)
	}
	flag, usage := heap.AddDataToSpecificQueue(data)
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

	for _, queueHeap := range q.heap {
		queue := ds.Information(queueHeap)
		queueMap := map[string]interface{}{
			"id":      queue["id"], // ✅ Access map values using keys
			"usage":   queue["usage"],
			"size":    queue["size"],
			"enabled": queue["enabled"],
			"dlq":     queue["dlq"],
		}
		result = append(result, queueMap)
		usageValue, usageOk := queue["usage"].(float32) // or float64, depending on original type
		sizeValue, sizeOk := queue["size"].(int)        // or float32 if size was stored as float

		if usageOk && sizeOk {
			usage += usageValue * float32(sizeValue)
			total += float32(sizeValue)
		}
	}

	for _, queueHeap := range q.reserved {
		queue := ds.Information(queueHeap)
		queueMap := map[string]interface{}{
			"id":      queue["id"], // ✅ Access map values using keys
			"usage":   queue["usage"],
			"size":    queue["size"],
			"enabled": queue["enabled"],
			"dlq":     queue["dlq"],
		}
		result = append(result, queueMap)
		usageValue, usageOk := queue["usage"].(float32) // or float64, depending on original type
		sizeValue, sizeOk := queue["size"].(int)        // or float32 if size was stored as float

		if usageOk && sizeOk {
			usage += usageValue * float32(sizeValue)
			total += float32(sizeValue)
		}
	}

	usage = usage / total

	return result, usage
}
