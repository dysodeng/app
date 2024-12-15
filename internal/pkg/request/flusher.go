package request

import (
	"net/http"
)

// Flusher is a wrapper for http.ResponseWriter with Flush method.
type Flusher struct {
	flusher        http.Flusher
	responseWriter http.ResponseWriter
}

func NewFlusher(writer http.ResponseWriter, httpFlusher http.Flusher) *Flusher {
	return &Flusher{
		flusher:        httpFlusher,
		responseWriter: writer,
	}
}

func (f *Flusher) writer(body string) (int, error) {
	return f.responseWriter.Write([]byte(body))
}

func (f *Flusher) flush() {
	f.flusher.Flush()
}

func (f *Flusher) Writer(body string) (int, error) {
	return f.writer(body)
}

func (f *Flusher) Flush() {
	f.flusher.Flush()
}

func (f *Flusher) WriterWithFlush(body string) (n int, err error) {
	n, err = f.writer(body)
	if err != nil {
		return
	}
	f.flusher.Flush()
	return
}
