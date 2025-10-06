package errors

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// 标准错误包的常用函数别名，方便使用
var (
	As     = errors.As
	Is     = errors.Is
	New    = errors.New
	Unwrap = errors.Unwrap
)

// DomainError 领域错误基础类型
type DomainError struct {
	// 错误码，用于唯一标识错误类型
	Code string
	// 错误消息，用于向用户展示
	Message string
	// 原始错误，用于错误链和调试
	Err error
	// 领域名称，标识错误所属的领域
	Domain string
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is errors.Is 判断规则
func (e *DomainError) Is(target error) bool {
	var t *DomainError
	if ok := errors.As(target, &t); !ok {
		return false
	}
	return e.Code == t.Code
}

// Wrap 包装一个错误（会修改原始错误）
func (e *DomainError) Wrap(err error) error {
	e.Err = err
	return e
}

// WrapNew 包装一个新的错误（不修改原始错误）
func (e *DomainError) WrapNew(err error) *DomainError {
	return &DomainError{
		Code:    e.Code,
		Message: e.Message,
		Domain:  e.Domain,
		Err:     err,
	}
}

func (e *DomainError) String() string {
	return e.Message
}

// Format 实现 fmt.Formatter 接口，支持自定义格式化
func (e *DomainError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// %+v 详细格式，显示完整错误链和堆栈信息
			_, _ = io.WriteString(s, e.formatVerboseWithStack())
			return
		}
		// %v 标准格式，显示简单的错误信息
		if e.Err != nil {
			_, _ = io.WriteString(s, e.Message)
			_, _ = io.WriteString(s, ": ")
			_, _ = io.WriteString(s, e.Err.Error())
			return
		}
		_, _ = io.WriteString(s, e.String())
	case 's':
		// %s 字符串格式
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		// %q 带引号的字符串格式
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

// formatVerbose 生成详细的错误信息，包括完整错误链（不含堆栈）
func (e *DomainError) formatVerbose() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("DomainError{Domain: %q, Code: %q, Message: %q", e.Domain, e.Code, e.Message))
	// 如果有包装的错误，递归显示错误链
	if e.Err != nil {
		builder.WriteString(", Err: ")
		var wrappedDomainErr *DomainError
		if errors.As(e.Err, &wrappedDomainErr) {
			builder.WriteString(wrappedDomainErr.formatVerbose())
		} else {
			builder.WriteString(fmt.Sprintf("%q", e.Err.Error()))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// formatVerboseWithStack 生成详细的错误信息，包括完整错误链和堆栈信息
func (e *DomainError) formatVerboseWithStack() string {
	var builder strings.Builder
	builder.WriteString(e.formatVerbose())
	if e.Err != nil {
		stackInfo := e.getStackTrace(e.Err)
		if stackInfo != "" {
			builder.WriteString("\n\nStack Trace:\n")
			builder.WriteString(stackInfo)
		}
	}
	return builder.String()
}

// getStackTrace 递归获取错误链中的堆栈信息
func (e *DomainError) getStackTrace(err error) string {
	if err == nil {
		return ""
	}

	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.getStackTrace(domainErr.Err)
	}
	stackStr := fmt.Sprintf("%+v", err)
	if strings.Contains(stackStr, "\n") {
		return stackStr
	}

	return ""
}

// FormatErrorChain 格式化完整的错误链为字符串
func FormatErrorChain(err error) string {
	var result strings.Builder
	current := err
	first := true

	for current != nil {
		if !first {
			result.WriteString(" -> ")
		}

		// 检查是否是 DomainError
		var domainErr *DomainError
		if errors.As(current, &domainErr) {
			if domainErr.Domain != "" {
				result.WriteString(fmt.Sprintf("[%s:%s] %s", domainErr.Domain, domainErr.Code, domainErr.Message))
			} else {
				result.WriteString(fmt.Sprintf("[%s] %s", domainErr.Code, domainErr.Message))
			}
		} else {
			result.WriteString(current.Error())
		}

		first = false
		current = errors.Unwrap(current)
	}

	return result.String()
}

// NewDomainError 创建带上下文的错误
func NewDomainError(domain, code, message string, err error) *DomainError {
	return &DomainError{
		Domain:  domain,
		Code:    code,
		Message: message,
		Err:     err,
	}
}