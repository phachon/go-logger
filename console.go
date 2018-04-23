package go_logger

import (
	"sync"
	"io"
	"strconv"
	"github.com/fatih/color"
	"os"
	"encoding/json"
	"reflect"
	"errors"
)

const CONSOLE_ADAPTER_NAME  = "console"

var levelColors = map[int] color.Attribute {
	LOGGER_LEVEL_EMERGENCY: color.FgWhite,  //white
	LOGGER_LEVEL_ALERT:     color.FgCyan,   //cyan
	LOGGER_LEVEL_CRITICAL:  color.FgMagenta,//magenta
	LOGGER_LEVEL_ERROR:     color.FgRed,    //red
	LOGGER_LEVEL_WARNING:   color.FgYellow, //yellow
	LOGGER_LEVEL_NOTICE:    color.FgGreen,  //green
	LOGGER_LEVEL_INFO:      color.FgBlue,   //blue
	LOGGER_LEVEL_DEBUG:     color.BgBlue,   //background blue
}

// adapter console
type AdapterConsole struct {
	write *ConsoleWriter
	config *ConsoleConfig
}

// console writer
type ConsoleWriter struct {
	lock sync.Mutex
	writer io.Writer
}

// console config
type ConsoleConfig struct {
	// console text is show color
	Color bool

	// is json format
	JsonFormat bool

	// is show file
	ShowFileLine bool
}

func (cc *ConsoleConfig) Name() string {
	return CONSOLE_ADAPTER_NAME
}

func NewAdapterConsole() LoggerAbstract {
	consoleWrite := &ConsoleWriter{
		writer: os.Stdout,
	}
	config := &ConsoleConfig{}
	return &AdapterConsole{
		write: consoleWrite,
		config : config,
	}
}

func (adapterConsole *AdapterConsole) Init(consoleConfig Config) error {
	if consoleConfig.Name() != CONSOLE_ADAPTER_NAME {
		return errors.New("logger console adapter init error, config must ConsoleConfig")
	}

	vc := reflect.ValueOf(consoleConfig)
	cc := vc.Interface().(*ConsoleConfig)
	adapterConsole.config = cc
	return nil
}

func (adapterConsole *AdapterConsole) Write(loggerMsg *loggerMessage) error {

	//timestamp := loggerMsg.Timestamp
	//timestampFormat := loggerMsg.TimestampFormat
	//millisecond := loggerMsg.Millisecond
	millisecondFormat := loggerMsg.MillisecondFormat
	body := loggerMsg.Body
	file := loggerMsg.File
	line := loggerMsg.Line
	levelString := loggerMsg.LevelString

	msg := ""
	if adapterConsole.config.JsonFormat == true  {
		jsonByte, _ := json.Marshal(loggerMsg)
		msg = string(jsonByte)
	}else {
		msg = millisecondFormat +" ["+ levelString + "] "
		if adapterConsole.config.ShowFileLine {
			msg += "[" + file + ":" + strconv.Itoa(line) + "] "
		}
		msg += body
	}

	if adapterConsole.config.Color {
		msg = adapterConsole.getColorByLevel(loggerMsg.Level, msg)
	}

	consoleWriter := adapterConsole.write
	consoleWriter.lock.Lock()
	consoleWriter.writer.Write([]byte(msg + "\n"))
	consoleWriter.lock.Unlock()
	return nil
}

func (adapterConsole *AdapterConsole) Name() string {
	return CONSOLE_ADAPTER_NAME
}

func (adapterConsole *AdapterConsole) Flush() {

}

func (adapterConsole *AdapterConsole) getColorByLevel(level int, content string) string {
	lc, ok := levelColors[level]
	if !ok {
		lc = color.FgWhite
	}
	colorFunc := color.New(lc).SprintFunc()
	return colorFunc(content)
}

func init()  {
	Register(CONSOLE_ADAPTER_NAME, NewAdapterConsole)
}