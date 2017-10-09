package go_logger

import (
	"sync"
	"os"
	"strconv"
)

const FILE_ADAPTER_NAME = "file"

type AdapterFile struct {
	write *FileWrite
	config map[string]interface{}
}

type FileWrite struct {
	lock sync.RWMutex
	writer *os.File
}

func NewAdapterFile() LoggerAbstract {

	//default file config
	defaultConfig := map[string]interface{}{
		"filename": "access.log",  // log file name
		"maxSize": 10, //log file max size (KB)
		"maxLine": 10, //log file max lines
	}
	fileWrite := &FileWrite{

	}
	return &AdapterFile{
		write: fileWrite,
		config: defaultConfig,
	}
}

func (adapterFile *AdapterFile) Init(config map[string]interface{}) {
	adapterFile.config = NewMisc().MapIntersect(adapterFile.config, config)
	filename := adapterFile.config["filename"].(string)
	//check file exits
	_, err := os.Stat(filename)
	if(os.IsNotExist(err)) {
		adapterFile.createFile(filename)
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0766)
	if(err != nil) {
		printError("adapter file open file " + filename + " error "+err.Error())
	}
	adapterFile.write.writer = file
}

func (adapterFile *AdapterFile) Write(loggerMsg *loggerMessage) error {

	//timestamp := loggerMsg.timestamp
	//timestampFormat := loggerMsg.timestampFormat
	//millisecond := loggerMsg.millisecond
	millisecondFormat := loggerMsg.millisecondFormat
	body := loggerMsg.body
	file := loggerMsg.file
	line := loggerMsg.line
	levelPrefix := levelMsgPrefix[loggerMsg.level]
	msg := millisecondFormat +" "+ levelPrefix + " [" + file + ":" + strconv.Itoa(line) + "] " + body + "\n"

	fileWrite := adapterFile.write
	fileWrite.lock.Lock()
	defer fileWrite.lock.Unlock()

	fileWrite.writer.Write([]byte(msg))
	return nil
}

func (adapterFile *AdapterFile) Flush() {
	adapterFile.write.writer.Close()
}

func (adapterFile *AdapterFile) Name() string {
	return FILE_ADAPTER_NAME
}

func (adapterFile *AdapterFile) createFile(filename string) {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	if err != nil {
		printError("create log file " + filename + " error " + err.Error())
	}
}

func init()  {
	Register(FILE_ADAPTER_NAME, NewAdapterFile)
}