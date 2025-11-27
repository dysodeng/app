package sse

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// flushRecorder 包装 ResponseRecorder 以实现 http.Flusher，并记录是否触发 Flush。
type flushRecorder struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (fr *flushRecorder) Flush() {
	fr.flushed = true
}

type multiLineJSON struct{}

func (multiLineJSON) MarshalJSON() ([]byte, error) {
	// 返回包含真实换行符的 JSON，以测试 WriteDataLines 的逐行输出
	return []byte("{\n\"a\":1,\n\"b\":2\n}"), nil
}

func TestNewSSEWriter_HeadersAndFlusher(t *testing.T) {
	// 正常：支持 Flusher
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, err := NewSSEWriter(fr)
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.Equal(t, "text/event-stream", fr.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", fr.Header().Get("Cache-Control"))
	assert.Equal(t, "keep-alive", fr.Header().Get("Connection"))
	assert.Equal(t, "no", fr.Header().Get("X-Accel-Buffering"))

	// 异常：不支持 Flusher
	nfr := &nonFlusherRecorder{rr: httptest.NewRecorder()}
	w2, err2 := NewSSEWriter(nfr)
	assert.Error(t, err2)
	assert.Nil(t, w2)
}

func TestWriter_WriteEvent(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteEvent("update", map[string]int{"x": 1}, "42")
	assert.NoError(t, err)

	expected := "id: 42\n" +
		"event: update\n" +
		"data: {\"x\":1}\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteEvent_EscapeFields(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteEvent("u<>", map[string]string{"m": "ok"}, "1\n2")
	assert.NoError(t, err)

	expected := "id: 1 2\n" +
		"event: u&lt;&gt;\n" +
		"data: {\"m\":\"ok\"}\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteData(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteData(map[string]any{"x": 1})
	assert.NoError(t, err)

	expected := "data: {\"x\":1}\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteDataWithID(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteDataWithID(map[string]any{"x": 1}, "9")
	assert.NoError(t, err)

	expected := "id: 9\n" +
		"data: {\"x\":1}\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteDataLines(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	// 使用常规对象，json.Marshal 默认生成单行
	err := w.WriteDataLines(map[string]int{"a": 1, "b": 2})
	assert.NoError(t, err)

	expected := "data: {\"a\":1,\"b\":2}\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteRetry(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteRetry(1500)
	assert.NoError(t, err)

	expected := "retry: 1500\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteID(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteID("abc\n123")
	assert.NoError(t, err)

	expected := "id: abc 123\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_WriteComment(t *testing.T) {
	fr := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w, _ := NewSSEWriter(fr)

	err := w.WriteComment("hello<world>\nline")
	assert.NoError(t, err)

	expected := ": hello&lt;world&gt; line\n\n"
	assert.Equal(t, expected, fr.Body.String())
}

func TestWriter_AndFlushVariants(t *testing.T) {
	// WriteEventAndFlush
	fr1 := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w1, _ := NewSSEWriter(fr1)
	err := w1.WriteEventAndFlush("e", map[string]any{"k": "v"}, "i")
	assert.NoError(t, err)
	assert.True(t, fr1.flushed)
	assert.Equal(t, "id: i\nevent: e\ndata: {\"k\":\"v\"}\n\n", fr1.Body.String())

	// WriteDataAndFlush
	fr2 := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w2, _ := NewSSEWriter(fr2)
	err = w2.WriteDataAndFlush(map[string]any{"a": 1})
	assert.NoError(t, err)
	assert.True(t, fr2.flushed)
	assert.Equal(t, "data: {\"a\":1}\n\n", fr2.Body.String())

	// WriteDataWithIDAndFlush
	fr3 := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w3, _ := NewSSEWriter(fr3)
	err = w3.WriteDataWithIDAndFlush(map[string]any{"a": 1}, "x")
	assert.NoError(t, err)
	assert.True(t, fr3.flushed)
	assert.Equal(t, "id: x\ndata: {\"a\":1}\n\n", fr3.Body.String())

	// WriteCommentAndFlush
	fr4 := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w4, _ := NewSSEWriter(fr4)
	err = w4.WriteCommentAndFlush("note")
	assert.NoError(t, err)
	assert.True(t, fr4.flushed)
	assert.Equal(t, ": note\n\n", fr4.Body.String())

	// WriteRetryAndFlush
	fr5 := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w5, _ := NewSSEWriter(fr5)
	err = w5.WriteRetryAndFlush(1000)
	assert.NoError(t, err)
	assert.True(t, fr5.flushed)
	assert.Equal(t, "retry: 1000\n\n", fr5.Body.String())

	// WriteIDAndFlush
	fr6 := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	w6, _ := NewSSEWriter(fr6)
	err = w6.WriteIDAndFlush("v1")
	assert.NoError(t, err)
	assert.True(t, fr6.flushed)
	assert.Equal(t, "id: v1\n\n", fr6.Body.String())
}

// nonFlusherRecorder 只实现 http.ResponseWriter，不实现 http.Flusher，用于触发 NewSSEWriter 的错误分支
type nonFlusherRecorder struct {
	rr *httptest.ResponseRecorder
}

func (n *nonFlusherRecorder) Header() http.Header {
	return n.rr.Header()
}

func (n *nonFlusherRecorder) Write(b []byte) (int, error) {
	return n.rr.Write(b)
}

func (n *nonFlusherRecorder) WriteHeader(statusCode int) {
	n.rr.WriteHeader(statusCode)
}
