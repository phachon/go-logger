package go_logger

import (
	"sync"
	"os"
	"strconv"
	"path"
	"strings"
	"time"
	"bufio"
	"io"
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
		"maxSize": 0,  // max file size, default 0 KB, unlimited size
		"maxLine": 0,  // max file line, default 0 line, Unlimited number
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
	maxSize := adapterFile.config["maxSize"].(int)
	maxLine := adapterFile.config["maxLine"].(int)

	//check file exits, otherwise create a file
	_, err := os.Stat(filename)
	if(os.IsNotExist(err)) {
		err = adapterFile.createFile(filename)
		if err != nil {
			printError(err.Error())
		}
	}
	//get a file pointer
	file, err := adapterFile.getFileObject(filename)
	if(err != nil) {
		printError(err.Error())
	}
	adapterFile.write.writer = file

	//monitor file size
	if maxSize > 0 {
		go adapterFile.monitorFileSize()
	}
	//monitor file line
	if maxLine > 0 {
		go adapterFile.monitorFileLines()
	}
}

//monitor file size, if maxSize < fileSize, rename file is file_time.log and recreate file
func (adapterFile *AdapterFile) monitorFileSize() {
	fileWrite := adapterFile.write
	maxSize := adapterFile.config["maxSize"].(int)
	filename := adapterFile.config["filename"].(string)
	for  {
		fileWrite.lock.RLock()
		fileSize, err := adapterFile.getFileSize(filename)
		if(err != nil) {
			fileWrite.lock.RUnlock()
			printError(err.Error())
		}
		fileSizeKb := float32(fileSize) / 1024  //kb
		if (float32(maxSize) < fileSizeKb) {
			saveTime := time.Now().Format("200601021504")
			filenameSuffix := path.Ext(filename)
			realName := strings.Replace(filename, filenameSuffix, "", 1) + "_"+ saveTime + filenameSuffix
			fileWrite.lock.RUnlock()
			fileWrite.lock.Lock()

			//close file handle
			fileWrite.writer.Close()
			err := os.Rename(filename, realName)
			if(err != nil) {
				fileWrite.lock.Unlock()
				continue
			}
			//recreate file
			err = adapterFile.createFile(filename)
			if(err != nil) {
				fileWrite.lock.Unlock()
				printError(err.Error())
			}
			fileWrite.lock.Unlock()
			//reset file write
			file, err := adapterFile.getFileObject(filename)
			if(err != nil) {
				printError(err.Error())
			}
			fileWrite.writer = file
		}else {
			fileWrite.lock.RUnlock()
		}
	}
}

//monitor file size, if maxLine < fileLine, rename file is file_time.log and recreate file
func (adapterFile *AdapterFile) monitorFileLines() {
	fileWrite := adapterFile.write
	maxLine := adapterFile.config["maxLine"].(int)
	filename := adapterFile.config["filename"].(string)
	for  {
		fileWrite.lock.RLock()
		fileLine, err := adapterFile.getFileLines(filename)
		if(err != nil) {
			fileWrite.lock.RUnlock()
			printError(err.Error())
		}
		if maxLine <= fileLine {
			saveTime := time.Now().Format("200601021504")
			filenameSuffix := path.Ext(filename)
			realName := strings.Replace(filename, filenameSuffix, "", 1) + "_"+ saveTime + filenameSuffix
			fileWrite.lock.RUnlock()
			fileWrite.lock.Lock()

			//close file handle
			fileWrite.writer.Close()
			err := os.Rename(filename, realName)
			if(err != nil) {
				fileWrite.lock.Unlock()
				continue
			}
			//recreate file
			err = adapterFile.createFile(filename)
			if(err != nil) {
				fileWrite.lock.Unlock()
				printError(err.Error())
			}
			fileWrite.lock.Unlock()
			//reset file write
			file, err := adapterFile.getFileObject(filename)
			if(err != nil) {
				printError(err.Error())
			}
			fileWrite.writer = file
		}else {
			fileWrite.lock.RUnlock()
		}
	}
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

//create file
//params : filename
//return error
func (adapterFile *AdapterFile) createFile(filename string) error {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	return err
}

//get file object
//params : filename
//return : *os.file, error
func (adapterFile *AdapterFile) getFileObject(filename string) (file *os.File, err error) {
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0766)
	return file, err
}

//get file size
//params : filename
//return : fileSize(byte int64), error
func (adapterFile *AdapterFile) getFileSize(filename string) (fileSize int64, err error) {
	fileInfo, err := os.Stat(filename)
	if(err != nil) {
		return fileSize, err
	}

	return fileInfo.Size(), nil
}

//get file lines
//params : filename
//return : fileLine, error
func (adapterFile *AdapterFile) getFileLines(filename string) (fileLine int, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0766)
	if(err != nil) {
		return fileLine, err
	}
	defer file.Close()

	fileLine = 1
	r := bufio.NewReader(file)
	for {
		_, err := r.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		fileLine += 1
	}
	return fileLine, nil
}

func init()  {
	Register(FILE_ADAPTER_NAME, NewAdapterFile)
}