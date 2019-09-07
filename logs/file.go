package logs

import (
	"fmt"
	"os"
	"time"
)

type FileLoggerOptions struct {
	filename      string
	lastSplitHour int
}

type FileLogger struct {
	file   *os.File
	option *FileLoggerOptions
}

func NewFileLogger(filename string) (LogInterface, error) {

	option := &FileLoggerOptions{
		filename: filename,
	}

	log := &FileLogger{
		option: option,
	}

	err := log.init()
	return log, err
}

func (f *FileLogger) getFilename() (filename string) {

	now := time.Now()
	filename = fmt.Sprintf("%s.%04d%02d%02d%02d", f.option.filename,
		now.Year(), now.Month(), now.Day(), now.Hour())
	return
}

func (f *FileLogger) init() (err error) {

	filename := f.getFilename()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("open faile %s failed, err:%v", filename, err)
	}

	f.file = file
	os.Symlink(filename, f.option.filename)
	f.option.lastSplitHour = time.Now().Hour()

	return
}

func (f *FileLogger) checkSplitFile(curTime time.Time) {

	hour := curTime.Hour()
	if hour == f.option.lastSplitHour {
		return
	}

	f.init()
}

func (f *FileLogger) Write(data *LogData) {
	f.checkSplitFile(data.curTime)
	f.file.Write(data.Bytes())
}

func (f *FileLogger) Close() {
	f.file.Close()
}
