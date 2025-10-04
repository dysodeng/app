package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logFileExt = ".log"

var _zapLogger *zap.Logger

func newZapLogger(debug bool) {
	zapEncoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",                       // 结构化（json）输出：msg的key
		LevelKey:    "level",                     // 结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:     "ts",                        // 结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:   "file",                      // 结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel: zapcore.CapitalLevelEncoder, // 将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) { // 采用文件相对路径编码输出
			_, file, line, ok := runtime.Caller(6)
			if !ok {
				encoder.AppendString(caller.TrimmedPath())
				return
			}
			encoder.AppendString(fmt.Sprintf("%s:%d", file, line))
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // 输出的时间格式
			enc.AppendString(t.Format(time.DateTime))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	fileWriter, err := logFileWriter()
	if err != nil {
		panic(err)
	}

	var cores []zapcore.Core
	if debug {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(zapEncoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)),
			zapcore.DebugLevel,
		))
	} else {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(zapEncoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter)),
			zapcore.InfoLevel,
		))
	}

	// 实现多个输出
	core := zapcore.NewTee(cores...)

	_zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func logFileWriter() (io.Writer, error) {
	filename := config.LogPath + "/app"
	return rotateLogs.New(
		filename+".%Y-%m-%d"+logFileExt,
		rotateLogs.WithLinkName(filename+logFileExt),
		rotateLogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotateLogs.WithRotationTime(time.Hour*24), // 切割频率 24小时
	)
}
