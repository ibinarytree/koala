package logs

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"
)

var (
	defaultOutputer    = NewConsoleOutputer()
	lm                 *LoggerMgr
	initOnce           *sync.Once = &sync.Once{}
	defaultServiceName            = "default"
)

type LoggerMgr struct {
	outputers   []Outputer
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
	if chanSize <= 0 {
		chanSize = DefaultLogChanSize
	}

	initLogger(level, chanSize, serviceName)
}

func SetLevel(level LogLevel) {
	lm.level = level
}

func (l *LoggerMgr) run() {

	for data := range l.logDataChan {
		if len(l.outputers) == 0 {
			defaultOutputer.Write(data)
			continue
		}

		for _, outputer := range l.outputers {
			outputer.Write(data)
		}
	}

	l.wg.Done()
}

func AddOutputer(ouputer Outputer) {
	if lm == nil {
		initLogger(LogLevelDebug, DefaultLogChanSize, defaultServiceName)
	}

	lm.outputers = append(lm.outputers, ouputer)
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

	for _, outputer := range lm.outputers {
		outputer.Close()
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
			logData.fields = fields
		}
	}

	select {
	case lm.logDataChan <- logData:
	default:
		return
	}
}
