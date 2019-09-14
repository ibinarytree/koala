package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileOutputerOptions struct {
	filename      string
	lastSplitHour int
}

type FileOutputer struct {
	file       *os.File
	accessFile *os.File
	option     *FileOutputerOptions
}

func NewFileOutputer(filename string) (Outputer, error) {

	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	option := &FileOutputerOptions{
		filename: filename,
	}

	log := &FileOutputer{
		option: option,
	}

	err = log.init()
	return log, err
}

func (f *FileOutputer) getCurFilename() (curFilename, originFilename string) {

	now := time.Now()
	curFilename = fmt.Sprintf("%s.%04d%02d%02d%02d", f.option.filename,
		now.Year(), now.Month(), now.Day(), now.Hour())
	originFilename = f.option.filename
	return
}

func (f *FileOutputer) getCurAccessFilename() (curAccessFilename, originAccessFilename string) {

	now := time.Now()
	curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d%02d", f.option.filename,
		now.Year(), now.Month(), now.Day(), now.Hour())

	originAccessFilename = fmt.Sprintf("%s.acccess", f.option.filename)
	return
}

func (f *FileOutputer) initFile(filename, originFilename string) (file *os.File, err error) {

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		err = fmt.Errorf("open faile %s failed, err:%v", filename, err)
		return
	}

	os.Symlink(filename, originFilename)
	return
}

func (f *FileOutputer) init() (err error) {

	curFilename, originFilename := f.getCurFilename()
	f.file, err = f.initFile(curFilename, originFilename)
	if err != nil {
		return
	}

	accessFilename, originAccessFilename := f.getCurAccessFilename()
	f.accessFile, err = f.initFile(accessFilename, originAccessFilename)
	if err != nil {
		return
	}

	f.option.lastSplitHour = time.Now().Hour()
	return
}

func (f *FileOutputer) checkSplitFile(curTime time.Time) {

	hour := curTime.Hour()
	if hour == f.option.lastSplitHour {
		return
	}

	f.init()
}

func (f *FileOutputer) Write(data *LogData) {
	f.checkSplitFile(data.curTime)

	file := f.file
	if data.level == LogLevelAccess {
		file = f.accessFile
	}

	file.Write(data.Bytes())
}

func (f *FileOutputer) Close() {
	f.file.Close()
	f.accessFile.Close()
}
