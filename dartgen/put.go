package dartgen

import "github.com/ab36245/go-writer"

var w *writer.Writer = writer.WithPrefix("  ")

func dec(mesg string, args ...any) {
	w.Back(mesg, args...)
	if mesg != "" {
		put("")
	}
}

func inc(mesg string, args ...any) {
	w.Over(mesg, args...)
}

func put(mesg string, args ...any) {
	w.End(mesg, args...)
}
