package sse

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientConnect_ParseStream(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SSE 必要响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher := w.(http.Flusher)

		// 注释心跳（应被忽略）
		_, _ = w.Write([]byte(": keepalive\n\n"))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond)

		// 事件1：默认事件名 message，无 id/retry
		_, _ = w.Write([]byte("data: {\"x\":1}\n\n"))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond)

		// 事件2：显式事件名 + id + retry
		_, _ = w.Write([]byte("event: update\n"))
		_, _ = w.Write([]byte("id: abc\n"))
		_, _ = w.Write([]byte("retry: 2000\n"))
		_, _ = w.Write([]byte("data: {\"a\":1}\n\n"))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond)

		// 事件3：payload 内含 event 字段，覆盖事件名（单行 JSON，解析更稳定）
		_, _ = w.Write([]byte("data: {\"event\":\"override\",\"y\":2}\n\n"))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond)
	}))
	defer ts.Close()

	client := NewSSEClient()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	evCh, errCh, stop, err := client.Connect(ctx, ts.URL)
	assert.NoError(t, err)
	defer stop()

	var events []Event
	deadline := time.After(300 * time.Millisecond)
	for {
		select {
		case ev, ok := <-evCh:
			if !ok {
				goto DONE
			}
			events = append(events, ev)
			if len(events) == 3 {
				goto DONE
			}
		case e := <-errCh:
			assert.NoError(t, e)
		case <-deadline:
			goto DONE
		}
	}
DONE:

	assert.Len(t, events, 3)

	// 事件1：默认 message
	assert.Equal(t, "", events[0].ID)
	assert.Equal(t, "message", events[0].Event)
	assert.Equal(t, 0, events[0].Retry)
	var p0 map[string]int
	assert.NoError(t, json.Unmarshal(events[0].Data, &p0))
	assert.Equal(t, 1, p0["x"])

	// 事件2：显式
	assert.Equal(t, "abc", events[1].ID)
	assert.Equal(t, "update", events[1].Event)
	assert.Equal(t, 2000, events[1].Retry)
	var p1 map[string]int
	assert.NoError(t, json.Unmarshal(events[1].Data, &p1))
	assert.Equal(t, 1, p1["a"])

	// 事件3：payload 覆盖事件名
	assert.Equal(t, "", events[2].ID)
	assert.Equal(t, "override", events[2].Event)
	var p2 map[string]any
	assert.NoError(t, json.Unmarshal(events[2].Data, &p2))
	assert.Equal(t, float64(2), p2["y"])
}

func TestClientConnect_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("wrong method"))
			return
		}
		if r.Header.Get("X-Test") != "ok" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("missing header"))
			return
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "hello" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad body"))
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher := w.(http.Flusher)

		_, _ = w.Write([]byte("id: z1\n"))
		_, _ = w.Write([]byte("data: {\"ok\":true}\n\n"))
		flusher.Flush()
	}))
	defer ts.Close()

	client := NewSSEClient()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	evCh, errCh, stop, err := client.Connect(
		ctx,
		ts.URL,
		WithMethod(http.MethodPost),
		WithHeaders(map[string]string{"X-Test": "ok"}),
		WithBody(strings.NewReader("hello")),
	)
	assert.NoError(t, err)
	defer stop()

	select {
	case ev := <-evCh:
		assert.Equal(t, "z1", ev.ID)
		var p map[string]bool
		assert.NoError(t, json.Unmarshal(ev.Data, &p))
		assert.True(t, p["ok"])
	case e := <-errCh:
		assert.NoError(t, e)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
}

func TestClientConnect_ErrorStatus(t *testing.T) {
	// 400：返回 body 文本作为错误
	ts400 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad param"))
	}))
	defer ts400.Close()

	client := NewSSEClient()
	_, _, _, err := client.Connect(context.Background(), ts400.URL)
	assert.Error(t, err)
	assert.Equal(t, "bad param", err.Error())

	// 500：unexpected status
	ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts500.Close()

	_, _, _, err = client.Connect(context.Background(), ts500.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status: 500")
}
