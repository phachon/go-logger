package go_logger

import (
	"sync"
	"io"
	"strconv"
	"github.com/fatih/color"
	"os"
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

type AdapterConsole struct {
	write *ConsoleWriter
	config map[string]interface{}
}

type ConsoleWriter struct {
	lock sync.Mutex
	writer io.Writer
}

func NewAdapterConsole() LoggerAbstract {

	//console default config
	defaultConfig := map[string]interface{}{
		"color": true, //the text need color
	}
	consoleWrite := &ConsoleWriter{
		writer: os.Stdout,
	}
	return &AdapterConsole{
		write: consoleWrite,
		config : defaultConfig,
	}
}

func (adapterConsole *AdapterConsole) Init(config map[string]interface{}) {
	adapterConsole.config = NewMisc().MapIntersect(adapterConsole.config, config)
}

func (adapterConsole *AdapterConsole) Write(loggerMsg *loggerMessage) error {

	//timestamp := loggerMsg.timestamp
	//timestampFormat := loggerMsg.timestampFormat
	//millisecond := loggerMsg.millisecond
	millisecondFormat := loggerMsg.millisecondFormat
	body := loggerMsg.body
	file := loggerMsg.file
	line := loggerMsg.line
	levelPrefix := levelMsgPrefix[loggerMsg.level]
	msg := millisecondFormat +" "+ levelPrefix + " [" + file + ":" + strconv.Itoa(line) + "] " + body

	if adapterConsole.config["color"].(bool) {
		msg = adapterConsole.getColorByLevel(loggerMsg.level, msg)
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