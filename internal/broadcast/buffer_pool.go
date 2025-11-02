package broadcast

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

func GetBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func PutBuffer(buf *bytes.Buffer) {

	if buf.Cap() <= 64*1024 {
		buf.Reset()
		bufferPool.Put(buf)
	}
}

func WithBuffer(fn func(*bytes.Buffer) error) error {
	buf := GetBuffer()
	defer PutBuffer(buf)
	return fn(buf)
}
