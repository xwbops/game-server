package zlog

import (
	"bytes"
	"runtime"
	"strings"
	"time"
)

//将日志输出到支持的输出中。
type Entry struct {
	logger *logger
	Buffer *bytes.Buffer
	Map    map[string]interface{}
	Level  Level
	Time   time.Time
	File   string
	Line   int
	Func   string
	Format string
	Args   []interface{}
}

func entry(logger *logger) *Entry {
	return &Entry{
		logger: logger,
		Buffer: new(bytes.Buffer),
		Map:    make(map[string]interface{}, 5),
	}
}

func (e *Entry) write(level Level, format string, args ...interface{}) {
	if e.logger.opt.level > level {
		return
	}
	e.Time = time.Now()
	e.Level = level
	e.Format = format
	e.Args = args
	if !e.logger.opt.disableCaller {
		if pc, file, line, ok := runtime.Caller(2); !ok {
			e.File = "???"
			e.Func = "???"
		} else {
			e.File, e.Line, e.Func = file, line, runtime.FuncForPC(pc).Name() //runtime.Caller() 来获取文件名和行号，调用runtime.Caller() 时，要注意传入正确的栈深度。
			e.Func = e.Func[strings.LastIndex(e.Func, "/")+1:]
		}
	}
	e.format()
	e.writer()
	e.release()
}
func (e *Entry) format() {
	_ = e.logger.opt.formatter.Format(e)
}

//即可将日志写入到指定的位置中
func (e *Entry) writer() {
	e.logger.mu.Lock()
	_, _ = e.logger.opt.output.Write(e.Buffer.Bytes())
	e.logger.mu.Unlock()
}

//调用release()方法来清空缓存和对象池
func (e *Entry) release() {
	e.Args, e.Line, e.File, e.Format, e.Func = nil, 0, "", "", ""
	e.Buffer.Reset()
	e.logger.entryPool.Put(e)
}
