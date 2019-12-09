package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var TaskLog *logrus.Logger

// 记录行号的hook
type TaskHook struct {
}

func (hook TaskHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook TaskHook) Fire(entry *logrus.Entry) error {
	goroutineName,ok := entry.Data["goroutine"]
	delete(entry.Data, "goroutine")
	if ok {
		entry.Message = fmt.Sprintf("goroutine[%v]: %s", goroutineName, entry.Message)

	}
	return nil
}

func init() {

	TaskLog = logrus.New()
	TaskLog.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	TaskLog.SetLevel(logrus.InfoLevel)
	TaskLog.AddHook(&TaskHook{})

}
