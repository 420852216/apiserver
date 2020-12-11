package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Log *zap.Logger

func ginHook(fileName string) *lumberjack.Logger {
	ginHook := &lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxSize:    128,     // 每个日志文件保存的大小 单位:M
		MaxAge:     7,       // 文件最多保存多少天
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		Compress:   false,   // 是否压缩
	}
	return ginHook
}

func sqlHook(fileName string) *lumberjack.Logger {
	sqlHook := &lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxSize:    128,     // 每个日志文件保存的大小 单位:M
		MaxAge:     7,       // 文件最多保存多少天
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		Compress:   false,   // 是否压缩
	}
	return sqlHook
}
func ginLogFormat() *zapcore.EncoderConfig{
	return &zapcore.EncoderConfig{
		MessageKey:     "msg",
	}
}

func sqlLogFormat() *zapcore.EncoderConfig{
	return &zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		//EncodeLevel: zapcore.CapitalColorLevelEncoder, //加入颜色
		//EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeTime:     timeFormat, //自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func InitLogger(debug bool) {
	ginHook := ginHook("./logs/all.log")
	sqlHook := sqlHook("./logs/sql.log")

	encoderConfig := zapcore.EncoderConfig{}

	// 设置日志级别
	//atomicLevel := zap.NewAtomicLevel()
	//atomicLevel.SetLevel(zap.DebugLevel)

	var writes = []zapcore.WriteSyncer{zapcore.AddSync(ginHook)}
	// 如果是开发环境，同时在控制台上也输出
	if debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	ginCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...), //可以配置多个输出方式
		zap.LevelEnablerFunc(infoLevel),
		//atomicLevel,
	)
	sqlCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(*sqlLogFormat()),
		zapcore.AddSync(sqlHook),
		zap.LevelEnablerFunc(debugLevel),
	)
	//zapcore.NewTee()可以承载多个core
	rootCore:=zapcore.NewTee(
		ginCore,
		sqlCore,
		)


	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	//field1 := zap.Fields(zap.String("appName", "name"))
	// 构造日志
	//Log = zap.New(rootCore, caller, development, field1)

	Log = zap.New(rootCore, caller, development)
}

func debugLevel(lvl zapcore.Level) bool {
	return lvl == zapcore.DebugLevel
}

func infoLevel(lvl zapcore.Level) bool {
	return lvl >= zapcore.InfoLevel
}

func timeFormat(t time.Time, enc zapcore.PrimitiveArrayEncoder){
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}