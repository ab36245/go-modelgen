package writer

import "github.com/ab36245/go-writer"

func From(writer *writer.Writer) GenWriter {
	return GenWriter{
		writer: writer,
	}
}

func WithPrefix(prefix string) GenWriter {
	return GenWriter{
		writer: writer.WithPrefix(prefix),
	}
}

type GenWriter struct {
	writer *writer.Writer
}

func (w GenWriter) Code() string {
	return w.writer.String()
}

func (w GenWriter) Dec(mesg string, args ...any) {
	w.writer.Back(mesg, args...)
	if mesg != "" {
		w.Put("")
	}
}

func (w GenWriter) Inc(mesg string, args ...any) {
	w.writer.Over(mesg, args...)
}

func (w GenWriter) Put(mesg string, args ...any) {
	w.writer.End(mesg, args...)
}
