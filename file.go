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
	"syscall"
)

const FILE_ADAPTER_NAME = "file"

const (
	FILE_SLICE_SIZE = "size"
	FILE_SLICE_LINE = "line"
	FILE_SLICE_DATE = "date"
)

const (
	FILE_SLICE_DATE_YEAR = "y"
	FILE_SLICE_DATE_MONTH = "m"
	FILE_SLICE_DATE_DAY = "d"
	FILE_SLICE_DATE_HOUR = "h"
	FILE_SLICE_DATE_MINUTE = "i"
	FILE_SLICE_DATE_SECOND = "s"
)

type AdapterFile struct {
	write *FileWrite
	config map[string]interface{}
	lastWriteTime int64
	sliceType string
	sliceDateType string
	sliceSizeMax int
	sliceLineMax int
}

type FileWrite struct {
	lock sync.RWMutex
	writer *os.File
}

var fileSliceMapping = map[string]interface{}{
	FILE_SLICE_SIZE: 0,   //file max size, default 0
	FILE_SLICE_LINE: 0,   //file max line, default 0
	FILE_SLICE_DATE: "",  //date y (year), m(month), d(day), h(hour), i(minute), s(second)
}

var fileSliceDateMapping = map[string]int{
	FILE_SLICE_DATE_YEAR: 0,
	FILE_SLICE_DATE_MONTH: 1,
	FILE_SLICE_DATE_DAY: 2,
	FILE_SLICE_DATE_HOUR: 3,
	FILE_SLICE_DATE_MINUTE: 4,
	FILE_SLICE_DATE_SECOND: 5,
}

func NewAdapterFile() LoggerAbstract {
	//default file config
	defaultConfig := map[string]interface{}{
		"filename": "access.log",  // log file name
		"slice": map[string]interface{}{},
	}
	fileWrite := &FileWrite{

	}
	return &AdapterFile{
		write: fileWrite,
		config: defaultConfig,
		lastWriteTime: 0,
		sliceType: "",
		sliceDateType: "",
		sliceSizeMax: 0,
		sliceLineMax: 0,
	}
}

func (adapterFile *AdapterFile) Init(config map[string]interface{}) {
	adapterFile.config = NewMisc().MapIntersect(adapterFile.config, config)
	configSlice := adapterFile.config["slice"].(map[string]interface{})
	filename := adapterFile.config["filename"].(string)

	if len(configSlice) > 1 {
		printError("file adapter config error: slice must be one of the 'size','line','date'!")
	}
	configSlice = NewMisc().MapIntersect(fileSliceMapping, configSlice)

	maxSize := configSlice["size"].(int)
	maxLine := configSlice["line"].(int)
	dateType := configSlice["date"].(string)

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

	//slice file by size
	if maxSize > 0 {
		adapterFile.sliceType = FILE_SLICE_SIZE
		adapterFile.sliceSizeMax = maxSize
		go adapterFile.sliceByFileSize()
	}
	//slice file by line
	if maxLine > 0 {
		adapterFile.sliceType = FILE_SLICE_LINE
		adapterFile.sliceLineMax = maxLine
		go adapterFile.sliceByFileLines()
	}
	//slice file by date
	if dateType != "" {
		_, ok := fileSliceDateMapping[dateType]
		if !ok {
			printError("file adapter config slice date error: slice date must be one of the 'y', 'd', 'm','h', 'i', 's'!")
		}
		adapterFile.sliceType = FILE_SLICE_DATE
		adapterFile.sliceDateType = dateType
		adapterFile.lastWriteTime, err = adapterFile.getFileLastTime(filename)
		if(err != nil) {
			printError(err.Error())
		}
	}
}

//slice file by size, if maxSize < fileSize, rename file is file_time.log and recreate file
func (adapterFile *AdapterFile) sliceByFileSize() {
	fileWrite := adapterFile.write
	filename := adapterFile.config["filename"].(string)
	maxSize := adapterFile.sliceSizeMax

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
			//reset file write
			file, err := adapterFile.getFileObject(filename)
			if(err != nil) {
				printError(err.Error())
			}
			fileWrite.writer = file

			fileWrite.lock.Unlock()
		}else {
			fileWrite.lock.RUnlock()
		}
	}
}

//slice file by line, if maxLine < fileLine, rename file is file_time.log and recreate file
func (adapterFile *AdapterFile) sliceByFileLines() {
	fileWrite := adapterFile.write
	filename := adapterFile.config["filename"].(string)
	maxLine := adapterFile.sliceLineMax

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
			//reset file write
			file, err := adapterFile.getFileObject(filename)
			if(err != nil) {
				printError(err.Error())
			}
			fileWrite.writer = file

			fileWrite.lock.Unlock()
		}else {
			fileWrite.lock.RUnlock()
		}
	}
}

//slice file by date (y, m, d, h, i, s), rename file is file_time.log and recreate file
func (adapterFile *AdapterFile) sliceByDate(nowTime int64) {
	fileWrite := adapterFile.write
	filename := adapterFile.config["filename"].(string)
	dateType := adapterFile.sliceDateType

	if adapterFile.lastWriteTime == 0 {
		adapterFile.lastWriteTime = nowTime
	}
	lastTime := adapterFile.lastWriteTime

	lastTimeUnix := time.Unix(lastTime, 0)
	nowTimeUnix := time.Unix(nowTime, 0)

	saveTime := ""
	if (dateType == FILE_SLICE_DATE_YEAR) &&
		(lastTimeUnix.Year() != nowTimeUnix.Year()) {
		saveTime = lastTimeUnix.Format("2006")
	}
	if (dateType == FILE_SLICE_DATE_MONTH) &&
		(lastTimeUnix.Format("200601") != nowTimeUnix.Format("200601")) {
		saveTime = lastTimeUnix.Format("200601")
	}
	if (dateType == FILE_SLICE_DATE_DAY) &&
		(lastTimeUnix.Format("20060102") != nowTimeUnix.Format("20060102")) {
		saveTime = lastTimeUnix.Format("20060102")
	}
	if (dateType == FILE_SLICE_DATE_HOUR) &&
		(lastTimeUnix.Format("2006010215") != nowTimeUnix.Format("2006010215")) {
		saveTime = lastTimeUnix.Format("2006010215")
	}
	if (dateType == FILE_SLICE_DATE_MINUTE) &&
		(lastTimeUnix.Format("200601021504") != nowTimeUnix.Format("200601021504")) {
		saveTime = lastTimeUnix.Format("200601021504")
	}
	if (dateType == FILE_SLICE_DATE_SECOND) &&
		(lastTimeUnix.Format("20060102150405") != nowTimeUnix.Format("20060102150405")) {
		saveTime = lastTimeUnix.Format("20060102150405")
	}

	if(saveTime != "") {
		filenameSuffix := path.Ext(filename)
		realName := strings.Replace(filename, filenameSuffix, "", 1) + "_" + saveTime + filenameSuffix

		//close file handle
		fileWrite.writer.Close()
		err := os.Rename(filename, realName)
		if (err != nil) {
			printError(err.Error())
		}
		//recreate file
		err = adapterFile.createFile(filename)
		if (err != nil) {
			printError(err.Error())
		}
		//reset file write
		file, err := adapterFile.getFileObject(filename)
		if (err != nil) {
			printError(err.Error())
		}
		fileWrite.writer = file
	}

	adapterFile.lastWriteTime = nowTime
}

func (adapterFile *AdapterFile) Write(loggerMsg *loggerMessage) error {

	timestamp := loggerMsg.Timestamp
	//timestampFormat := loggerMsg.TimestampFormat
	//millisecond := loggerMsg.Millisecond
	millisecondFormat := loggerMsg.MillisecondFormat
	body := loggerMsg.Body
	file := loggerMsg.File
	line := loggerMsg.Line
	levelString := loggerMsg.LevelString
	msg := millisecondFormat +" ["+ levelString + "] [" + file + ":" + strconv.Itoa(line) + "] " + body + "\n"

	fileWrite := adapterFile.write
	fileWrite.lock.Lock()
	defer fileWrite.lock.Unlock()

	if adapterFile.sliceType == FILE_SLICE_DATE {
		adapterFile.sliceByDate(timestamp)
	}
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

//get file create time
//params : filename
//return : createTime(int64), error
func (adapterFile *AdapterFile) getFileCreateTime(filename string) (createTime int64, err error) {
	fileInfo, err := os.Lstat(filename)
	if(err != nil) {
		return createTime, err
	}
	fileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	return fileSys.CreationTime.Nanoseconds()/1e9, nil
}

//get file last time
//params : filename
//return : last time(int64), error
func (adapterFile *AdapterFile) getFileLastTime(filename string) (lastWriteTime int64, err error) {
	fileInfo, err := os.Lstat(filename)
	if(err != nil) {
		return lastWriteTime, err
	}
	fileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	return fileSys.LastWriteTime.Nanoseconds()/1e9, nil
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