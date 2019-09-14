package logs

type Outputer interface {
	Write(data *LogData)
	Close()
}
