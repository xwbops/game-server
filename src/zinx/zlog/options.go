package zlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// Level 定义日志级别和日志选项。
// 定义日志级别
type Level uint8

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

const (
	FmtEmptySeparate = ""
)

//在日志输出时，要通过对比开关级别和输出级别的大小，来决定是否输出，所以日志级别Level要定义成方便比较的数值类型。几乎所有的日志包都是用常量计数器iota来定义日志级别。
//定义：定义日志级别和日志选项。
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

//另外，因为要在日志输出中，输出可读的日志级别（例如输出INFO而不是1），所以需要有Level到Level Name的映射LevelNameMapping，LevelNameMapping会在格式化时用到。
//日志级别和字符串名称映射
var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO":
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

type options struct {
	output        io.Writer
	level         Level //日志级别
	stdLevel      Level
	formatter     Formatter //输出格式 text/json
	disableCaller bool      //是否开启文件名和行号。
}

//为了灵活地设置日志的选项，你可以通过选项模式，来对日志选项进行设置：
type Option func(o *options)

//为了灵活地设置日志的选项，你可以通过选项模式，来对日志选项进行设置：
func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.output == nil {
		o.output = os.Stderr
	}
	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}
	return
}

//具有选项模式的日志包，可通过以下方式，来动态地修改日志的选项：
//zlog.SetOptions(zlog.WithLevel(zlog.DebugLevel))
//WithOutput（output io.Writer）：设置输出位置。
func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

//WithLevel（level Level）：设置输出级别。
func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}
func WithStdLevel(level Level) Option {
	return func(o *options) {
		o.stdLevel = level
	}
}

//WithFormatter（formatter Formatter）：设置输出格式。
func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

//WithDisableCaller（caller bool）：设置是否打印文件名和行号。
func WithDisableCaller(caller bool) Option {
	return func(o *options) {
		o.disableCaller = caller
	}
}
