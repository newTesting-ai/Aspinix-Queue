package ds

import "sync"

type RingBuffer struct {
	data  []interface{}
	size  int
	head  int
	tail  int
	count int
	mu    sync.Mutex
}

func InitRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:  make([]interface{}, size),
		size:  size,
		count: 0,
	}
}

func (r *RingBuffer) InsertDataToRingBuffer(data interface{}) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == r.size {
		return false, r.count
	}

	r.data[r.tail] = data
	r.tail = (r.tail + 1) % r.size
	r.count++
	return true, r.count
}

func (r *RingBuffer) GetDataFromRingBuffer() (interface{}, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		return nil, false
	}

	data := r.data[r.head]
	r.head = (r.head + 1) % r.size
	r.count--
	return data, true
}

func (r *RingBuffer) GetBufferSize() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.count
}
