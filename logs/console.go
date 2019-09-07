package logs

import (
	"os"
)

type ConsoleLogger struct {
	level LogLevel
}

func NewConsoleLogger(logLevel string) (log LogInterface) {
	level := getLogLevel(logLevel)
	log = &ConsoleLogger{
		level: level,
	}
	return
}

func (c *ConsoleLogger) SetLevel(level LogLevel) {
	if level < LogLevelDebug || level > LogLevelError {
		level = LogLevelDebug
	}

	c.level = level
}

func (c *ConsoleLogger) Write(data *LogData) {

	color := getLevelColor(data.level)
	text := color.Add(string(data.Bytes()))
	os.Stdout.Write([]byte(text))
}

func (c *ConsoleLogger) Close() {

}
