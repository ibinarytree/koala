package logs

type LogInterface interface {
	Write(data *LogData)
	Close()
}
