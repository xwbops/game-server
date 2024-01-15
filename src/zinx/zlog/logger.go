package zlog

import (
	"io"
	"sync"
	"unsafe"
)

var std = New()

//设置日志信息
type logger struct {
	opt       *options
	mu        sync.Mutex
	entryPool *sync.Pool
}

func New(opts ...Option) *logger {
	logger := &logger{opt: initOptions(opts...)}
	logger.entryPool = &sync.Pool{New: func() interface{} { return entry(logger) }}
	return logger
}
func (l *logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, opt := range opts {
		opt(l.opt)
	}
}
func (l *logger) Writer() io.Writer {
	return l
}
func (l *logger) Write(data []byte) (int, error) {
	l.entry().write(l.opt.stdLevel, FmtEmptySeparate, *(*string)(unsafe.Pointer(&data)))
}
func (l *logger) entry() *Entry {
	return l.entryPool.Get().(*Entry)
}
func StdLogger() *logger {
	return std
}
func SetOptions(opts ...Option) {
	std.SetOptions(opts...)
}
func Writer() io.Writer {
	return std
}
func (l *logger) Debug(args ...interface{}) {
	l.entry.Wr
}
