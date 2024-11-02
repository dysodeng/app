package qrcode

import "bytes"

// bufferWriter 实现
type bufferWriter struct {
	*bytes.Buffer
}

func (buf bufferWriter) Write(p []byte) (int, error) {
	return buf.Buffer.Write(p)
}

func (buf bufferWriter) Close() error {
	return nil
}
