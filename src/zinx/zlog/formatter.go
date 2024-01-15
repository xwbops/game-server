package zlog

type Formatter interface {
	Format(entry *Entry) error
}
