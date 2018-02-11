package go_logger

import (
	"sync"
	"os"
	"strconv"
	"path"
	"strings"
	"time"
	"go-logger/utils"
)

const FILE_ADAPTER_NAME = "file"

const (
	FILE_SLICE_DATE_NULL = ""
	FILE_SLICE_DATE_YEAR = "y"
	FILE_SLICE_DATE_MONTH = "m"
	FILE_SLICE_DATE_DAY = "d"
	FILE_SLICE_DATE_HOUR = "h"
)

type AdapterFile struct {
	write *FileWrite
	filename string
	maxSize int64
	maxLine int64
	startLine int64
	startTime int64
	dateSlice string
}

type FileWrite struct {
	lock sync.RWMutex
	writer *os.File
}

var fileSliceDateMapping = map[string]int{
	FILE_SLICE_DATE_YEAR: 0,
	FILE_SLICE_DATE_MONTH: 1,
	FILE_SLICE_DATE_DAY: 2,
	FILE_SLICE_DATE_HOUR: 3,
}

//default file config
var defaultConfig = map[string]interface{}{
	"filename": "access.log",  // default log file name
	"maxSize":  1024 * 1024,  // default max size 1G size
	"maxLine":  100000,         // default max line 100000
	"dateSlice": FILE_SLICE_DATE_NULL, // default date slice is day
}

func NewAdapterFile() LoggerAbstract {
	return &AdapterFile{
		write: &FileWrite{},
		filename: "",
		maxSize: 0,
		maxLine: 0,
		startLine: 0,
		startTime: 0,
		dateSlice: "",
	}
}

// init
func (adapterFile *AdapterFile) Init(config map[string]interface{}) {
	config = utils.NewMisc().MapIntersect(defaultConfig, config)

	adapterFile.filename = config["filename"].(string)
	adapterFile.maxSize = config["maxSize"].(int64)
	adapterFile.maxLine = config["maxLine"].(int64)
	adapterFile.dateSlice = config["dateSlice"].(string)

	_, ok := fileSliceDateMapping[adapterFile.dateSlice]
	if !ok {
		printError("file adapter config slice date error: slice date must be one of the 'y', 'd', 'm','h'!")
	}

	adapterFile.initFile()
}

// init file
func (adapterFile *AdapterFile) initFile()  {

	//check file exits, otherwise create a file
	ok, _ := utils.NewFile().PathExists(adapterFile.filename)
	if ok == false {
		err := utils.NewFile().CreateFile(adapterFile.filename)
		if err != nil {
			printError(err.Error())
		}
	}

	// get start time
	adapterFile.startTime = time.Now().Unix()

	// get file start lines
	nowLines, err := utils.NewFile().GetFileLines(adapterFile.filename)
	if err != nil {
		printError(err.Error())
	}
	adapterFile.startLine = nowLines

	//get a file pointer
	file, err := adapterFile.getFileObject(adapterFile.filename)
	if err != nil {
		printError(err.Error())
	}
	adapterFile.write.writer = file
}

func (adapterFile *AdapterFile) Write(loggerMsg *loggerMessage) error {

	//timestamp := loggerMsg.Timestamp
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

	if adapterFile.dateSlice != "" {
		// file slice by date
		adapterFile.sliceByDate()
	}
	if adapterFile.maxLine != 0 {
		// file slice by line
		adapterFile.sliceByFileLines()
	}
	if adapterFile.maxSize != 0 {
		// file slice by line
		adapterFile.sliceByFileSize()
	}
	fileWrite.writer.Write([]byte(msg))
	if adapterFile.maxLine != 0 {
		adapterFile.startLine += int64(strings.Count(msg, "\n"))
	}
	return nil
}

func (adapterFile *AdapterFile) Flush() {
	adapterFile.write.writer.Close()
}

func (adapterFile *AdapterFile) Name() string {
	return FILE_ADAPTER_NAME
}

//slice file by date (y, m, d, h, i, s), rename file is file_time.log and recreate file
func (adapterFile *AdapterFile) sliceByDate() {

	fileWrite := adapterFile.write
	filename := adapterFile.filename
	filenameSuffix := path.Ext(filename)
	startTime := time.Unix(adapterFile.startTime, 0)
	nowTime := time.Now()

	oldFilename := ""
	isHaveSlice := false
	if (adapterFile.dateSlice == FILE_SLICE_DATE_YEAR) &&
		(startTime.Year() != nowTime.Year()) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("2006") + filenameSuffix
	}
	if (adapterFile.dateSlice == FILE_SLICE_DATE_MONTH) &&
		(startTime.Format("200601") != nowTime.Format("200601")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("200601") + filenameSuffix
	}
	if (adapterFile.dateSlice == FILE_SLICE_DATE_DAY) &&
		(startTime.Format("20060102") != nowTime.Format("20060102")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("20060102") + filenameSuffix
	}
	if (adapterFile.dateSlice == FILE_SLICE_DATE_HOUR) &&
		(startTime.Format("2006010215") != startTime.Format("2006010215")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("2006010215") + filenameSuffix
	}

	if isHaveSlice == true  {
		//close file handle
		fileWrite.writer.Close()
		err := os.Rename(adapterFile.filename, oldFilename)
		if err != nil {
			printError(err.Error())
		}
		adapterFile.initFile()
	}
}

//slice file by line, if maxLine < fileLine, rename file is file_line_maxLine_time.log and recreate file
func (adapterFile *AdapterFile) sliceByFileLines() {
	fileWrite := adapterFile.write
	filename := adapterFile.filename
	filenameSuffix := path.Ext(filename)
	maxLine := adapterFile.maxLine
	startLine := adapterFile.startLine
	randStr := utils.NewMisc().RandString(4)

	if startLine >= maxLine {
		//close file handle
		fileWrite.writer.Close()
		oldFilename := strings.Replace(filename, filenameSuffix, "", 1) + "_line_"+strconv.FormatInt(maxLine, 10)+"_"+randStr+filenameSuffix
		err := os.Rename(filename, oldFilename)
		if err != nil {
			printError(err.Error())
		}
		adapterFile.initFile()
	}
}

//slice file by size, if maxSize < fileSize, rename file is file_size_maxSize_time.log and recreate file
func (adapterFile *AdapterFile) sliceByFileSize() {
	fileWrite := adapterFile.write
	filename := adapterFile.filename
	filenameSuffix := path.Ext(filename)
	maxSize := adapterFile.maxSize
	nowSize, _ := adapterFile.getFileSize(filename)
	randStr := utils.NewMisc().RandString(4)

	if nowSize >= maxSize {
		//close file handle
		fileWrite.writer.Close()
		oldFilename := strings.Replace(filename, filenameSuffix, "", 1) + "_size_"+strconv.FormatInt(maxSize, 10)+"_"+randStr+filenameSuffix
		err := os.Rename(filename, oldFilename)
		if err != nil {
			printError(err.Error())
		}
		adapterFile.initFile()
	}
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

func init()  {
	Register(FILE_ADAPTER_NAME, NewAdapterFile)
}