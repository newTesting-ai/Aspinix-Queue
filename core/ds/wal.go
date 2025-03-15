package ds

import (
	"os"
	"sync"
)

type WAL struct {
	file *os.File
	mu   sync.Mutex
}

func NewWal(filePath string) (*WAL, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{
		file: file,
	}, nil
}

func (wal *WAL) appendLogToWAL(log string) error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	_, err := wal.file.WriteString(log + "\n")
	return err
}
