package go_logger

import (
	"errors"
	"github.com/fatih/color"
	"io"
	"os"
	"reflect"
	"sync"
)

const CONSOLE_ADAPTER_NAME = "console"

var levelColors = map[int]color.Attribute{
	LOGGER_LEVEL_EMERGENCY: color.FgWhite,   //white
	LOGGER_LEVEL_ALERT:     color.FgCyan,    //cyan
	LOGGER_LEVEL_CRITICAL:  color.FgMagenta, //magenta
	LOGGER_LEVEL_ERROR:     color.FgRed,     //red
	LOGGER_LEVEL_WARNING:   color.FgYellow,  //yellow
	LOGGER_LEVEL_NOTICE:    color.FgGreen,   //green
	LOGGER_LEVEL_INFO:      color.FgBlue,    //blue
	LOGGER_LEVEL_DEBUG:     color.BgBlue,    //background blue
}

// adapter console
type AdapterConsole struct {
	write  *ConsoleWriter
	config *ConsoleConfig
}

// console writer
type ConsoleWriter struct {
	lock   sync.Mutex
	writer io.Writer
}

// console config
type ConsoleConfig struct {
	// console text is show color
	Color bool

	// is json format
	JsonFormat bool

	// jsonFormat is false, please input format string
	// if format is empty, default format "%millisecond_format% [%level_string%] %body%"
	//
	//  Timestamp "%timestamp%"
	//	TimestampFormat "%timestamp_format%"
	//	Millisecond "%millisecond%"
	//	MillisecondFormat "%millisecond_format%"
	//	Level int "%level%"
	//	LevelString "%level_string%"
	//	Body string "%body%"
	//	File string "%file%"
	//	Line int "%line%"
	//	Function "%function%"
	//
	// example: format = "%millisecond_format% [%level_string%] %body%"
	Format string
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
		write:  consoleWrite,
		config: config,
	}
}

func (adapterConsole *AdapterConsole) Init(consoleConfig Config) error {
	if consoleConfig.Name() != CONSOLE_ADAPTER_NAME {
		return errors.New("logger console adapter init error, config must ConsoleConfig")
	}

	vc := reflect.ValueOf(consoleConfig)
	cc := vc.Interface().(*ConsoleConfig)
	adapterConsole.config = cc

	if cc.JsonFormat == false && cc.Format == "" {
		cc.Format = defaultLoggerMessageFormat
	}

	return nil
}

func (adapterConsole *AdapterConsole) Write(loggerMsg *loggerMessage) error {

	msg := ""
	if adapterConsole.config.JsonFormat == true {
		//jsonByte, _ := json.Marshal(loggerMsg)
		jsonByte, _ := loggerMsg.MarshalJSON()
		msg = string(jsonByte)
	} else {
		msg = loggerMessageFormat(adapterConsole.config.Format, loggerMsg)
	}
	consoleWriter := adapterConsole.write

	if adapterConsole.config.Color {
		colorAttr := adapterConsole.getColorByLevel(loggerMsg.Level, msg)
		consoleWriter.lock.Lock()
		color.New(colorAttr).Println(msg)
		consoleWriter.lock.Unlock()
		return nil
	}

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

func (adapterConsole *AdapterConsole) getColorByLevel(level int, content string) color.Attribute {
	lc, ok := levelColors[level]
	if !ok {
		lc = color.FgWhite
	}
	return lc
}

func init() {
	Register(CONSOLE_ADAPTER_NAME, NewAdapterConsole)
}
