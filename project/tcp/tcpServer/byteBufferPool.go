package tcpServer

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ByteBufferPool struct {
	freeBuffers    chan []byte
	bufferSize     int
	maxPoolSize    int
	totalAllocated int64
}

var byteBufferPool *ByteBufferPool
var poolOnce sync.Once

func GetBufferPool() *ByteBufferPool {
	poolOnce.Do(func() {
		byteBufferPool = newByteBufferPool(100, 5*1024)
	})
	return byteBufferPool
}

func newByteBufferPool(maxPoolSize, bufferSize int) *ByteBufferPool {
	return &ByteBufferPool{
		freeBuffers: make(chan []byte, maxPoolSize),
		bufferSize:  bufferSize,
		maxPoolSize: maxPoolSize,
	}
}

func (p *ByteBufferPool) GetBuffer() []byte {
	fmt.Println("try GetBuffer ...")
	if atomic.LoadInt64(&p.totalAllocated) < int64(p.maxPoolSize) {
		if atomic.AddInt64(&p.totalAllocated, 1) <= int64(p.maxPoolSize) {
			fmt.Printf("GetBuffer success, currentPoolSize [%d]\n", p.bufferSize)
			return make([]byte, p.bufferSize)
		}
		atomic.AddInt64(&p.totalAllocated, -1)
	}

	// 否则，从池中取（会阻塞直到有 buffer 归还）
	return <-p.freeBuffers
}

func (p *ByteBufferPool) ReturnBuffer(buf []byte) {
	if cap(buf) != p.bufferSize {
		return
	}

	buf = buf[:0]
	buf = append(buf, make([]byte, p.bufferSize)...)

	select {
	case p.freeBuffers <- buf:

	default:
		// 池已满（理论上不会发生，因为 freeBuffers 容量 = maxPoolSize）
		// 此时 buffer 会被 GC 回收
	}
}
