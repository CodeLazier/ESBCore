package helper

import (
	"fmt"
	"os"

	MLog "github.com/RichardKnop/machinery/v1/log"
	RpcLog "github.com/smallnest/rpcx/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogParams struct {
	LogPath   string `yaml:"path"`
	LogSize   int    `yaml:"size"`
	LogMaxAge int    `yaml:"saveDays"`
	LogLevel  int    `yaml:"level"`
}

func initLogger(fileName string, size, age, level int) *zap.Logger {

	wsync := zapcore.AddSync(&lumberjack.Logger{
		Filename: fileName,
		MaxSize:  size, // megabytes
		MaxAge:   age,  // days

	})

	logConfig := zap.NewDevelopmentEncoderConfig()

	logLevel := zap.DebugLevel
	//config.EncodeLevel=zapcore.LowercaseColorLevelEncoder
	switch level {
	case 0:
		logLevel = zap.DebugLevel
	case 1:
		logLevel = zap.InfoLevel
	case 2:
		logLevel = zap.WarnLevel
	case 3:
		logLevel = zap.ErrorLevel
	case 4:
		logLevel = zap.PanicLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), wsync),
		logLevel, //调试等级
	)

	return zap.New(core, zap.AddCaller())
}

func NewAdapterLogger(fileName string, size, age, level int) *AdapterLogger {
	ad := &AdapterLogger{}
	ad.Logger = initLogger(fileName, size, age, level)
	MLog.Set(ad)
	RpcLog.SetLogger(ad)
	return ad
}

type AdapterLogger struct {
	*zap.Logger
}

func (l *AdapterLogger) Debug(v ...interface{}) {
	l.Logger.Debug(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, v...), zap.Skip())
}

func (l *AdapterLogger) Info(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Infof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Logger.Info(s, zap.Skip())
}

func (l *AdapterLogger) Warn(v ...interface{}) {
	l.Logger.Warn(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, v...), zap.Skip())
}

func (l *AdapterLogger) Error(v ...interface{}) {
	l.Logger.Error(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...), zap.Skip())
}

func (l *AdapterLogger) Fatal(v ...interface{}) {
	l.Logger.Fatal(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf(format, v...), zap.Skip())
}

func (l *AdapterLogger) Panic(v ...interface{}) {
	l.Logger.Panic(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panic(fmt.Sprint(v...), zap.Skip())
}

func (l *AdapterLogger) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *AdapterLogger) Printf(f string, v ...interface{}) {
	l.Infof(f, v...)
}

func (l *AdapterLogger) Println(v ...interface{}) {
	l.Infof("\n", v...)
}

func (l *AdapterLogger) Fatalln(v ...interface{}) {
	l.Fatalf("\n", v...)
}

func (l *AdapterLogger) Panicln(v ...interface{}) {
	l.Panicf("\n", v...)
}
