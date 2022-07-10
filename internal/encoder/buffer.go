package encoder

import "sync"

type buffer struct {
	b []byte
}

var bufferPool = sync.Pool{New: func() any {
	return &buffer{b: make([]byte, 0, 1024)}
}}

func newBuffer() *buffer {
	return bufferPool.Get().(*buffer)
}

func freeBuffer(buf *buffer) {
	buf.b = buf.b[:0]
	bufferPool.Put(buf)
}
