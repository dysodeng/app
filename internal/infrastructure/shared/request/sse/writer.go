package sse

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
)

// Writer 用于向客户端发送 Server-Sent Events（SSE）数据流。
// 封装了 http.ResponseWriter 与 http.Flusher，提供线程安全、带缓冲刷新、
// 符合 W3C SSE 规范的事件写入能力。
//
// 主要特性：
//  1. 自动设置 SSE 所需的响应头（Content-Type: text/event-stream 等）。
//  2. 提供多种写入方法：事件（event）、数据（data）、注释（comment）、
//     重连间隔（retry）、事件 ID（id）等，支持链式调用与立即刷新。
//  3. 所有写入方法最终都会调用底层 ResponseWriter，并可通过 Flush() 立即
//     将缓冲数据推送到客户端，实现低延迟实时推送。
//  4. 支持带 ID 的消息，用于客户端断线后通过 Last-Event-ID 恢复数据流。
//  5. 支持多行 data 的规范拆分，确保 JSON 中的换行不会破坏 SSE 格式。
//
// 使用示例：
//
//	w, err := sse.NewSSEWriter(rw)
//	if err != nil {
//	    http.Error(rw, err.Error(), http.StatusInternalServerError)
//	    return
//	}
//	defer w.Flush() // 确保最后残留数据立即推送
//
//	// 发送一个带事件名与 ID 的 JSON 消息
//	if err := w.WriteEvent("message", map[string]string{"msg": "hello"}, "1"); err != nil {
//	    log.Printf("write event failed: %v", err)
//	    return
//	}
//
//	// 立即刷新到客户端
//	w.Flush()
//
// 注意事项：
//   - 必须在响应开始前调用 NewSSEWriter，否则无法设置响应头。
//   - 写入完成后务必调用 Flush()，避免数据滞留在缓冲区。
//   - 若客户端不支持 Flusher，NewSSEWriter 会返回错误，需提前处理。
type Writer struct {
	writer  http.ResponseWriter
	flusher http.Flusher
}

// NewSSEWriter 创建并初始化一个 SSE Writer 实例。
// 该函数会检查传入的 http.ResponseWriter 是否支持 http.Flusher 接口，
// 若不支持则返回错误；若支持，则设置 SSE 所需的响应头（如 Content-Type、Cache-Control 等），
// 并返回一个可用于向客户端推送 Server-Sent Events 的 Writer。
func NewSSEWriter(w http.ResponseWriter) (*Writer, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("http.ResponseWriter flusher not supported")
	}

	// SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	return &Writer{
		writer:  w,
		flusher: flusher,
	}, nil
}

// Write 向底层 http.ResponseWriter 写入原始字节数据，返回实际写入的字节数及可能发生的错误。
// 本方法直接透传至 http.ResponseWriter.Write，可用于发送任意自定义格式的 SSE 帧。
// 注意：写入后数据仍可能驻留缓冲区，如需立即推送到客户端，应手动调用 Flush()。
func (w *Writer) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

// Flush 刷新缓冲区，将所有待发送数据立即推送到客户端。
// 本方法应在每次写入完成后调用，确保客户端及时接收事件。
func (w *Writer) Flush() {
	w.flusher.Flush()
}

// WriteEvent 写入标准SSE事件
func (w *Writer) WriteEvent(event string, payload any, id string) error {
	var builder strings.Builder
	if id != "" {
		builder.WriteString(fmt.Sprintf("id: %s\n", escapeSSEFieldValue(id)))
	}

	builder.WriteString(fmt.Sprintf("event: %s\n", escapeSSEFieldValue(event)))

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		builder.WriteString("data: ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	_, err = w.Write(helper.StringToBytes(builder.String()))
	return err
}

// WriteData 仅写入data
func (w *Writer) WriteData(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var builder strings.Builder
	builder.WriteString("data: ")
	builder.Write(data)
	builder.WriteString("\n\n")

	_, err = w.Write(helper.StringToBytes(builder.String()))
	return err
}

// WriteDataWithID 仅写入data且可选id
func (w *Writer) WriteDataWithID(payload any, id string) error {
	var builder strings.Builder
	if id != "" {
		builder.WriteString(fmt.Sprintf("id: %s\n", escapeSSEFieldValue(id)))
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		builder.WriteString("data: ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	_, err = w.Write(helper.StringToBytes(builder.String()))
	return err
}

// WriteComment 注释心跳或备注（SSE注释行），每条消息以空行结束
func (w *Writer) WriteComment(comment string) error {
	var builder strings.Builder
	builder.WriteString(": ")
	builder.WriteString(escapeSSEFieldValue(comment))
	builder.WriteString("\n\n")
	_, err := w.Write(helper.StringToBytes(builder.String()))
	return err
}

// WriteRetry 设置浏览器自动重连的间隔（毫秒）
func (w *Writer) WriteRetry(ms int) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("retry: %d\n\n", ms))
	_, err := w.Write(helper.StringToBytes(builder.String()))
	return err
}

// WriteID 仅写入id一条消息（用于同步Last-Event-ID游标）
func (w *Writer) WriteID(id string) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("id: %s\n\n", escapeSSEFieldValue(id)))
	_, err := w.Write(helper.StringToBytes(builder.String()))
	return err
}

// WriteDataLines 写入多行data（当JSON包含换行时更规范）
func (w *Writer) WriteDataLines(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var builder strings.Builder
	for _, line := range splitLines(string(data)) {
		builder.WriteString("data: ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("\n")
	_, err = w.Write(helper.StringToBytes(builder.String()))
	return err
}

func (w *Writer) WriteEventAndFlush(event string, payload any, id string) error {
	err := w.WriteEvent(event, payload, id)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (w *Writer) WriteDataAndFlush(payload any) error {
	err := w.WriteData(payload)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (w *Writer) WriteDataWithIDAndFlush(payload any, id string) error {
	err := w.WriteDataWithID(payload, id)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (w *Writer) WriteCommentAndFlush(comment string) error {
	err := w.WriteComment(comment)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (w *Writer) WriteRetryAndFlush(ms int) error {
	err := w.WriteRetry(ms)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (w *Writer) WriteIDAndFlush(id string) error {
	err := w.WriteID(id)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (w *Writer) WriteDataLinesAndFlush(payload any) error {
	err := w.WriteDataLines(payload)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

// escapeSSEFieldValue 转义 SSE 字段值中的特殊字符，防止破坏 SSE 格式和 XSS 攻击
// 注意：这里只转义字段值，不影响 SSE 协议本身的换行符
func escapeSSEFieldValue(value string) string {
	// 将字段值中的换行符替换为空格，防止破坏 SSE 的行结构
	// SSE 格式要求每个字段占一行，字段值中不能包含换行符
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")

	// HTML 转义，防止 XSS 攻击
	value = html.EscapeString(value)

	return value
}
