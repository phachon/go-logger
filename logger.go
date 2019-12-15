package go_logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Version = "v1.2"

const (
	LOGGER_LEVEL_EMERGENCY = iota
	LOGGER_LEVEL_ALERT
	LOGGER_LEVEL_CRITICAL
	LOGGER_LEVEL_ERROR
	LOGGER_LEVEL_WARNING
	LOGGER_LEVEL_NOTICE
	LOGGER_LEVEL_INFO
	LOGGER_LEVEL_DEBUG
)

type adapterLoggerFunc func() LoggerAbstract

type LoggerAbstract interface {
	Name() string
	Init(config Config) error
	Write(loggerMsg *loggerMessage) error
	Flush()
}

var adapters = make(map[string]adapterLoggerFunc)

var levelStringMapping = map[int]string{
	LOGGER_LEVEL_EMERGENCY: "Emergency",
	LOGGER_LEVEL_ALERT:     "Alert",
	LOGGER_LEVEL_CRITICAL:  "Critical",
	LOGGER_LEVEL_ERROR:     "Error",
	LOGGER_LEVEL_WARNING:   "Warning",
	LOGGER_LEVEL_NOTICE:    "Notice",
	LOGGER_LEVEL_INFO:      "Info",
	LOGGER_LEVEL_DEBUG:     "Debug",
}

var defaultLoggerMessageFormat = "%millisecond_format% [%level_string%] %body%"

//Register logger adapter
func Register(adapterName string, newLog adapterLoggerFunc) {
	if adapters[adapterName] != nil {
		panic("logger: logger adapter " + adapterName + " already registered!")
	}
	if newLog == nil {
		panic("logger: logger adapter " + adapterName + " is nil!")
	}

	adapters[adapterName] = newLog
}

type Logger struct {
	lock        sync.Mutex          //sync lock
	outputs     []*outputLogger     // outputs loggers
	msgChan     chan *loggerMessage // message channel
	synchronous bool                // is sync
	wait        sync.WaitGroup      // process wait
	signalChan  chan string
}

type outputLogger struct {
	Name  string
	Level int
	LoggerAbstract
}

type loggerMessage struct {
	Timestamp         int64  `json:"timestamp"`
	TimestampFormat   string `json:"timestamp_format"`
	Millisecond       int64  `json:"millisecond"`
	MillisecondFormat string `json:"millisecond_format"`
	Level             int    `json:"level"`
	LevelString       string `json:"level_string"`
	Body              string `json:"body"`
	File              string `json:"file"`
	Line              int    `json:"line"`
	Function          string `json:"function"`
}

//new logger
//return logger
func NewLogger() *Logger {
	logger := &Logger{
		outputs:     []*outputLogger{},
		msgChan:     make(chan *loggerMessage, 10),
		synchronous: true,
		wait:        sync.WaitGroup{},
		signalChan:  make(chan string, 1),
	}
	//default adapter console
	logger.attach("console", LOGGER_LEVEL_DEBUG, &ConsoleConfig{})

	return logger
}

//start attach a logger adapter
//param : adapterName console | file | database | ...
//return : error
func (logger *Logger) Attach(adapterName string, level int, config Config) error {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	return logger.attach(adapterName, level, config)
}

//attach a logger adapter after lock
//param : adapterName console | file | database | ...
//return : error
func (logger *Logger) attach(adapterName string, level int, config Config) error {
	for _, output := range logger.outputs {
		if output.Name == adapterName {
			printError("logger: adapter " + adapterName + "already attached!")
		}
	}
	logFun, ok := adapters[adapterName]
	if !ok {
		printError("logger: adapter " + adapterName + "is nil!")
	}
	adapterLog := logFun()
	err := adapterLog.Init(config)
	if err != nil {
		printError("logger: adapter " + adapterName + " init failed, error: " + err.Error())
	}

	output := &outputLogger{
		Name:           adapterName,
		Level:          level,
		LoggerAbstract: adapterLog,
	}

	logger.outputs = append(logger.outputs, output)
	return nil
}

//start attach a logger adapter
//param : adapterName console | file | database | ...
//return : error
func (logger *Logger) Detach(adapterName string) error {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	return logger.detach(adapterName)
}

//detach a logger adapter after lock
//param : adapterName console | file | database | ...
//return : error
func (logger *Logger) detach(adapterName string) error {
	outputs := []*outputLogger{}
	for _, output := range logger.outputs {
		if output.Name == adapterName {
			continue
		}
		outputs = append(outputs, output)
	}
	logger.outputs = outputs
	return nil
}

//set logger level
//params : level int
//func (logger *Logger) SetLevel(level int) {
//	logger.level = level
//}

//set logger synchronous false
//params : sync bool
func (logger *Logger) SetAsync(data ...int) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.synchronous = false

	msgChanLen := 100
	if len(data) > 0 {
		msgChanLen = data[0]
	}

	logger.msgChan = make(chan *loggerMessage, msgChanLen)
	logger.signalChan = make(chan string, 1)

	if !logger.synchronous {
		go func() {
			defer func() {
				e := recover()
				if e != nil {
					fmt.Printf("%v", e)
				}
			}()
			logger.startAsyncWrite()
		}()
	}
}

//write log message
//params : level int, msg string
//return : error
func (logger *Logger) Writer(level int, msg string) error {
	funcName := "null"
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "null"
		line = 0
	} else {
		funcName = runtime.FuncForPC(pc).Name()
	}
	_, filename := path.Split(file)

	if levelStringMapping[level] == "" {
		printError("logger: level " + strconv.Itoa(level) + " is illegal!")
	}

	loggerMsg := &loggerMessage{
		Timestamp:         time.Now().Unix(),
		TimestampFormat:   time.Now().Format("2006-01-02 15:04:05"),
		Millisecond:       time.Now().UnixNano() / 1e6,
		MillisecondFormat: time.Now().Format("2006-01-02 15:04:05.999"),
		Level:             level,
		LevelString:       levelStringMapping[level],
		Body:              msg,
		File:              filename,
		Line:              line,
		Function:          funcName,
	}

	if !logger.synchronous {
		logger.wait.Add(1)
		logger.msgChan <- loggerMsg
	} else {
		logger.writeToOutputs(loggerMsg)
	}

	return nil
}

//sync write message to loggerOutputs
//params : loggerMessage
func (logger *Logger) writeToOutputs(loggerMsg *loggerMessage) {
	for _, loggerOutput := range logger.outputs {
		// write level
		if loggerOutput.Level >= loggerMsg.Level {
			err := loggerOutput.Write(loggerMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "logger: unable write loggerMessage to adapter:%v, error: %v\n", loggerOutput.Name, err)
			}
		}
	}
}

//start async write by read logger.msgChan
func (logger *Logger) startAsyncWrite() {
	for {
		select {
		case loggerMsg := <-logger.msgChan:
			logger.writeToOutputs(loggerMsg)
			logger.wait.Done()
		case signal := <-logger.signalChan:
			if signal == "flush" {
				logger.flush()
			}
		}
	}
}

//flush msgChan data
func (logger *Logger) flush() {
	if !logger.synchronous {
		for {
			if len(logger.msgChan) > 0 {
				loggerMsg := <-logger.msgChan
				logger.writeToOutputs(loggerMsg)
				logger.wait.Done()
				continue
			}
			break
		}
		for _, loggerOutput := range logger.outputs {
			loggerOutput.Flush()
		}
	}
}

//if SetAsync() or logger.synchronous is false, must call Flush() to flush msgChan data
func (logger *Logger) Flush() {
	if !logger.synchronous {
		logger.signalChan <- "flush"
		logger.wait.Wait()
		return
	}
	logger.flush()
}

func (logger *Logger) LoggerLevel(levelStr string) int {
	levelStr = strings.ToUpper(levelStr)
	switch levelStr {
	case "EMERGENCY":
		return LOGGER_LEVEL_EMERGENCY
	case "ALERT":
		return LOGGER_LEVEL_ALERT
	case "CRITICAL":
		return LOGGER_LEVEL_CRITICAL
	case "ERROR":
		return LOGGER_LEVEL_ERROR
	case "WARNING":
		return LOGGER_LEVEL_WARNING
	case "NOTICE":
		return LOGGER_LEVEL_NOTICE
	case "INFO":
		return LOGGER_LEVEL_INFO
	case "DEBUG":
		return LOGGER_LEVEL_DEBUG
	default:
		return LOGGER_LEVEL_DEBUG
	}
}

func loggerMessageFormat(format string, loggerMsg *loggerMessage) string {
	message := strings.Replace(format, "%timestamp%", strconv.FormatInt(loggerMsg.Timestamp, 10), 1)
	message = strings.Replace(message, "%timestamp_format%", loggerMsg.TimestampFormat, 1)
	message = strings.Replace(message, "%millisecond%", strconv.FormatInt(loggerMsg.Millisecond, 10), 1)
	message = strings.Replace(message, "%millisecond_format%", loggerMsg.MillisecondFormat, 1)
	message = strings.Replace(message, "%level%", strconv.Itoa(loggerMsg.Level), 1)
	message = strings.Replace(message, "%level_string%", loggerMsg.LevelString, 1)
	message = strings.Replace(message, "%file%", loggerMsg.File, 1)
	message = strings.Replace(message, "%line%", strconv.Itoa(loggerMsg.Line), 1)
	message = strings.Replace(message, "%function%", loggerMsg.Function, 1)
	message = strings.Replace(message, "%body%", loggerMsg.Body, 1)

	return message
}

//log emergency level
func (logger *Logger) Emergency(msg string) {
	logger.Writer(LOGGER_LEVEL_EMERGENCY, msg)
}

//log emergency format
func (logger *Logger) Emergencyf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_EMERGENCY, msg)
}

//log alert level
func (logger *Logger) Alert(msg string) {
	logger.Writer(LOGGER_LEVEL_ALERT, msg)
}

//log alert format
func (logger *Logger) Alertf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_ALERT, msg)
}

//log critical level
func (logger *Logger) Critical(msg string) {
	logger.Writer(LOGGER_LEVEL_CRITICAL, msg)
}

//log critical format
func (logger *Logger) Criticalf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_CRITICAL, msg)
}

//log error level
func (logger *Logger) Error(msg string) {
	logger.Writer(LOGGER_LEVEL_ERROR, msg)
}

//log error format
func (logger *Logger) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_ERROR, msg)
}

//log warning level
func (logger *Logger) Warning(msg string) {
	logger.Writer(LOGGER_LEVEL_WARNING, msg)
}

//log warning format
func (logger *Logger) Warningf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_WARNING, msg)
}

//log notice level
func (logger *Logger) Notice(msg string) {
	logger.Writer(LOGGER_LEVEL_NOTICE, msg)
}

//log notice format
func (logger *Logger) Noticef(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_NOTICE, msg)
}

//log info level
func (logger *Logger) Info(msg string) {
	logger.Writer(LOGGER_LEVEL_INFO, msg)
}

//log info format
func (logger *Logger) Infof(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_INFO, msg)
}

//log debug level
func (logger *Logger) Debug(msg string) {
	logger.Writer(LOGGER_LEVEL_DEBUG, msg)
}

//log debug format
func (logger *Logger) Debugf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Writer(LOGGER_LEVEL_DEBUG, msg)
}

func printError(message string) {
	fmt.Println(message)
	os.Exit(0)
}
