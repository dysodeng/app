package logger

import (
	"io"
	"os"
	"time"

	"github.com/dysodeng/app/internal/config"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logFileExt = ".log"

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func init() {
	zapEncoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",                       // 结构化（json）输出：msg的key
		LevelKey:     "level",                     // 结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "time",                      // 结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      // 结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, // 将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  // 采用完整文件路径编码输出（/path/test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // 输出的时间格式
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	writer, err := logWriter()
	if err != nil {
		panic(err)
	}

	// 实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapEncoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writer)),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.InfoLevel
			}),
		),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	sugarLogger = logger.Sugar()
}

func logWriter() (io.Writer, error) {
	filename := config.LogPath + "/" + config.App.Name
	return rotateLogs.New(
		filename+".%Y-%m-%d"+logFileExt,
		rotateLogs.WithLinkName(filename+logFileExt),
		rotateLogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotateLogs.WithRotationTime(time.Hour*24), // 切割频率 24小时
	)
}

// Logger 获取结构化日志实例(性能高)
func Logger() *zap.Logger {
	return logger
}

func SugarLogger() *zap.SugaredLogger {
	return sugarLogger
}
