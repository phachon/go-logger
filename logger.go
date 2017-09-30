package go_logger

import (
	"fmt"
	"sync"
	"runtime"
	"path"
	"time"
	"os"
	"strconv"
)

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
	Init(config map[string]interface{})
	Write(loggerMsg *loggerMessage) error
	Flush()
}

var adapters = make(map[string]adapterLoggerFunc)

var levelMsgPrefix = map[int]string{
	LOGGER_LEVEL_EMERGENCY:   "[Emergency]",
	LOGGER_LEVEL_ALERT:       "[Alert]",
	LOGGER_LEVEL_CRITICAL:    "[Critical]",
	LOGGER_LEVEL_ERROR:       "[Error]",
	LOGGER_LEVEL_WARNING:     "[Warning]",
	LOGGER_LEVEL_NOTICE:      "[Notice]",
	LOGGER_LEVEL_INFO:        "[Info]",
	LOGGER_LEVEL_DEBUG:       "[Debug]",
}

//Register logger adapter
func Register(adapterName string, newLog adapterLoggerFunc)  {
	if adapters[adapterName] != nil {
		panic("logger: logger adapter "+ adapterName +" already registered!")
	}
	if newLog == nil {
		panic("logger: looger adapter "+ adapterName +" is nil!")
	}

	adapters[adapterName] = newLog
}

type outputLogger struct {
	Name string
	LoggerAbstract
}

type Logger struct {
	level       int //日志级别
	lock        sync.Mutex //互斥锁
	outputs     []*outputLogger //输出 loggers
	msgChan     chan *loggerMessage // message channel 通道
	synchronous bool //同步
	wait        sync.WaitGroup //线程阻塞
	signalChan  chan string //信号 channel
}

type loggerMessage struct {
	time int64
	level int
	body string
	file string
	line int
	function string
}

//new logger
//return logger
func NewLogger() *Logger {
	return &Logger{
		level:          LOGGER_LEVEL_EMERGENCY,
		outputs:        []*outputLogger{},
		synchronous:    true,
		wait:           sync.WaitGroup{},
	}
}

//start attach a logger adapter
//param : adapterName console | file | database | ...
//return : error
func (logger *Logger) Attach(adapterName string, config map[string]interface{}) error {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	return logger.attach(adapterName, config)
}

//attach a logger adapter after lock
//param : adapterName console | file | database | ...
//return : error
func (logger *Logger) attach(adapterName string, config map[string]interface{}) (error) {
	for _, output := range logger.outputs {
		if(output.Name == adapterName) {
			return printError("logger: adapter " +adapterName+ "already attached!")
		}
	}
	log, ok := adapters[adapterName]
	if !ok {
		return printError("logger: adapter " +adapterName+ "is nil!")
	}
	adapterLog := log()
	adapterLog.Init(config)

	output := &outputLogger {
		Name:adapterName,
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
		if(output.Name == adapterName) {
			continue
		}
		outputs = append(outputs, output)
	}
	logger.outputs = outputs
	return nil
}

//set logger level
//params : level int
func (logger *Logger) SetLevel(level int) {
	logger.level = level
}

//set logger synchronous
//params : sync bool
func (logger *Logger) SetSync(isSync bool) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.synchronous = isSync

	logger.msgChan = make(chan *loggerMessage, 10)
	logger.signalChan = make(chan string, 1)

	if (!isSync) {
		logger.wait.Add(1)
		go logger.startAsyncWrite()
	}
}

//write log message
//params : level int, msg string
func (logger *Logger) Writer(level int, msg string) error {
	funcName := "null"
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "null"
		line = 0
	}else {
		fun := runtime.FuncForPC(pc)
		funcName = fun.Name()
	}
	_, filename := path.Split(file)

	msgPrefix := levelMsgPrefix[level]
	if(msgPrefix == "") {
		return printError("logger: level " + strconv.Itoa(level) + " is illegal!")
	}

	loggerMsg := &loggerMessage {
		time :time.Now().Unix(),
		level :level,
		body: msg,
		file : filename,
		line : line,
		function: funcName,
	}

	if(!logger.synchronous) {
		logger.msgChan <- loggerMsg
	}else {
		logger.writeToOutputs(loggerMsg)
	}

	return nil
}

//sync write message to loggerOutputs
func (logger *Logger) writeToOutputs(loggerMsg *loggerMessage)  {
	for adapterName, loggerOutput := range logger.outputs {
		err := loggerOutput.Write(loggerMsg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logger: unable write loggerMessage to adapter:%v,error:%v\n", adapterName, err)
		}
	}
}

//start async write by read logger.msgChan
func (logger *Logger) startAsyncWrite()  {
	for {
		select {
		case loggerMsg := <-logger.msgChan:
			logger.writeToOutputs(loggerMsg)
		case signal := <-logger.signalChan:
			if signal == "flush" {
				logger.flush()
			}
			logger.wait.Done()
		}
	}
}

func (logger *Logger) flush() {
	if !logger.synchronous {
		for {
			if len(logger.msgChan) > 0 {
				loggerMsg := <-logger.msgChan
				logger.writeToOutputs(loggerMsg)
				continue
			}
			break
		}
		for _, loggerOutput := range logger.outputs {
			loggerOutput.Flush()
		}
	}
}

func (logger *Logger) Flush()  {
	if !logger.synchronous {
		logger.signalChan <- "flush"
		logger.wait.Wait()
		return
	}
	logger.flush()
}

func (logger *Logger) Emergency(msg string) {
	if logger.level < LOGGER_LEVEL_EMERGENCY {
		return
	}

	logger.Writer(LOGGER_LEVEL_EMERGENCY, msg)
}


func printError(message string) error {
	return fmt.Errorf("%s", message)
}