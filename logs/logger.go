package logs

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"
)

var (
	defaultLogger      = NewConsoleLogger("debug")
	lm                 *LoggerMgr
	initOnce           *sync.Once = &sync.Once{}
	defaultServiceName            = "default"
)

type LoggerMgr struct {
	loggers     []LogInterface
	chanSize    int
	level       LogLevel
	logDataChan chan *LogData
	serviceName string
	wg          sync.WaitGroup
}

func initLogger(level LogLevel, chanSize int, serviceName string) {
	initOnce.Do(func() {
		lm = &LoggerMgr{
			chanSize:    chanSize,
			level:       level,
			serviceName: serviceName,
			logDataChan: make(chan *LogData, chanSize),
		}
		lm.wg.Add(1)
		go lm.run()
	})
}

func InitLogger(level LogLevel, chanSize int, serviceName string) {
	initLogger(level, chanSize, serviceName)
}

func SetLevel(level LogLevel) {
	lm.level = level
}

func (l *LoggerMgr) run() {
	for data := range l.logDataChan {
		if len(l.loggers) == 0 {
			defaultLogger.Write(data)
			continue
		}

		for _, logger := range l.loggers {
			logger.Write(data)
		}
	}

	l.wg.Done()
}

func AddLogger(logger LogInterface) {
	if lm == nil {
		initLogger(LogLevelDebug, DefaultLogChanSize, defaultServiceName)
	}

	lm.loggers = append(lm.loggers, logger)
	return
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelDebug, format, args...)
}

func Trace(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelTrace, format, args...)
}

func Access(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelAccess, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelInfo, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelWarn, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelError, format, args...)
}

func Stop() {
	close(lm.logDataChan)
	lm.wg.Wait()

	for _, logger := range lm.loggers {
		logger.Close()
	}

	//重新初始化
	initOnce = &sync.Once{}
	lm = nil
}

func writeLog(ctx context.Context, level LogLevel, format string, args ...interface{}) {

	if lm == nil {
		initLogger(LogLevelDebug, DefaultLogChanSize, defaultServiceName)
	}

	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")

	fileName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		message:     msg,
		curTime:     now,
		timeStr:     nowStr,
		level:       level,
		filename:    fileName,
		lineNo:      lineNo,
		traceId:     GetTraceId(ctx),
		serviceName: lm.serviceName,
	}

	//access日志的时候,需要把所有field拉取出来
	if level == LogLevelAccess {
		fields := getFields(ctx)
		if fields != nil {
			logData.fields = make(map[interface{}]interface{})
			fields.fields.Range(func(k, v interface{}) bool {
				logData.fields[k] = v
				return true
			})
		}
	}

	select {
	case lm.logDataChan <- logData:
	default:
		return
	}
}
