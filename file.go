package go_logger

import (
	"sync"
	"os"
	"strconv"
	"path"
	"strings"
	"time"
	"go-logger/utils"
	"errors"
	"encoding/json"
)

const FILE_ADAPTER_NAME = "file"

const (
	FILE_SLICE_DATE_NULL = ""
	FILE_SLICE_DATE_YEAR = "y"
	FILE_SLICE_DATE_MONTH = "m"
	FILE_SLICE_DATE_DAY = "d"
	FILE_SLICE_DATE_HOUR = "h"
)

// adapter file
type AdapterFile struct {
	write *FileWriter
	config *FileConfig
	startLine int64
	startTime int64
}

// file writer
type FileWriter struct {
	lock sync.RWMutex
	writer *os.File
}

// file config
type FileConfig struct {

	// log file name
	Filename string

	// max file size
	MaxSize  int64

	// max file line
	MaxLine  int64

	// file slice by date
	// "y" Log files are cut through year
	// "m" Log files are cut through mouth
	// "d" Log files are cut through day
	// "h" Log files are cut through hour
	DateSlice string

	// is json format
	JsonFormat bool
}

var fileSliceDateMapping = map[string]int{
	FILE_SLICE_DATE_YEAR: 0,
	FILE_SLICE_DATE_MONTH: 1,
	FILE_SLICE_DATE_DAY: 2,
	FILE_SLICE_DATE_HOUR: 3,
}

func NewAdapterFile() LoggerAbstract {
	return &AdapterFile{
		write: &FileWriter{},
		config: &FileConfig{},
		startLine: 0,
		startTime: 0,
	}
}

// init
func (adapterFile *AdapterFile) Init(config *Config) error {

	adapterFile.config = config.File

	if adapterFile.config.Filename == "" {
		return errors.New("config Filename can't be empty!")
	}
	_, ok := fileSliceDateMapping[adapterFile.config.DateSlice]
	if !ok {
		return errors.New("config DateSlice must be one of the 'y', 'd', 'm','h'!")
	}

	err := adapterFile.initFile()
	return err
}

// init file
func (adapterFile *AdapterFile) initFile() error {

	//check file exits, otherwise create a file
	ok, _ := utils.NewFile().PathExists(adapterFile.config.Filename)
	if ok == false {
		err := utils.NewFile().CreateFile(adapterFile.config.Filename)
		if err != nil {
			return err
		}
	}

	// get start time
	adapterFile.startTime = time.Now().Unix()

	// get file start lines
	nowLines, err := utils.NewFile().GetFileLines(adapterFile.config.Filename)
	if err != nil {
		return err
	}
	adapterFile.startLine = nowLines

	//get a file pointer
	file, err := adapterFile.getFileObject(adapterFile.config.Filename)
	if err != nil {
		return err
	}
	adapterFile.write.writer = file

	return nil
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

	fileWrite := adapterFile.write
	fileWrite.lock.Lock()
	defer fileWrite.lock.Unlock()

	if adapterFile.config.DateSlice != "" {
		// file slice by date
		err := adapterFile.sliceByDate()
		if err != nil {
			return err
		}
	}
	if adapterFile.config.MaxLine != 0 {
		// file slice by line
		err := adapterFile.sliceByFileLines()
		if err != nil {
			return err
		}
	}
	if adapterFile.config.MaxSize != 0 {
		// file slice by size
		err := adapterFile.sliceByFileSize()
		if err != nil {
			return err
		}
	}

	msg := ""
	if adapterFile.config.JsonFormat == true  {
		jsonByte, _ := json.Marshal(loggerMsg)
		msg = string(jsonByte) + "\n"
	}else {
		msg = millisecondFormat +" ["+ levelString + "] [" + file + ":" + strconv.Itoa(line) + "] " + body + "\n"
	}

	fileWrite.writer.Write([]byte(msg))
	if adapterFile.config.MaxLine != 0 {
		if adapterFile.config.JsonFormat == true {
			adapterFile.startLine += 1
		}else {
			adapterFile.startLine += int64(strings.Count(msg, "\n"))
		}
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
func (adapterFile *AdapterFile) sliceByDate() error {

	fileWrite := adapterFile.write
	filename := adapterFile.config.Filename
	filenameSuffix := path.Ext(filename)
	startTime := time.Unix(adapterFile.startTime, 0)
	nowTime := time.Now()

	oldFilename := ""
	isHaveSlice := false
	if (adapterFile.config.DateSlice == FILE_SLICE_DATE_YEAR) &&
		(startTime.Year() != nowTime.Year()) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("2006") + filenameSuffix
	}
	if (adapterFile.config.DateSlice == FILE_SLICE_DATE_MONTH) &&
		(startTime.Format("200601") != nowTime.Format("200601")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("200601") + filenameSuffix
	}
	if (adapterFile.config.DateSlice == FILE_SLICE_DATE_DAY) &&
		(startTime.Format("20060102") != nowTime.Format("20060102")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("20060102") + filenameSuffix
	}
	if (adapterFile.config.DateSlice == FILE_SLICE_DATE_HOUR) &&
		(startTime.Format("2006010215") != startTime.Format("2006010215")) {
		isHaveSlice = true
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "_" + startTime.Format("2006010215") + filenameSuffix
	}

	if isHaveSlice == true  {
		//close file handle
		fileWrite.writer.Close()
		err := os.Rename(adapterFile.config.Filename, oldFilename)
		if err != nil {
			return err
		}
		err = adapterFile.initFile()
		if err != nil {
			return err
		}
	}

	return nil
}

//slice file by line, if maxLine < fileLine, rename file is file_line_maxLine_time.log and recreate file
func (adapterFile *AdapterFile) sliceByFileLines() error {
	fileWrite := adapterFile.write
	filename := adapterFile.config.Filename
	filenameSuffix := path.Ext(filename)
	maxLine := adapterFile.config.MaxLine
	startLine := adapterFile.startLine
	randStr := utils.NewMisc().RandString(4)

	if startLine >= maxLine {
		//close file handle
		fileWrite.writer.Close()
		oldFilename := strings.Replace(filename, filenameSuffix, "", 1) + "_line_"+strconv.FormatInt(maxLine, 10)+"_"+randStr+filenameSuffix
		err := os.Rename(filename, oldFilename)
		if err != nil {
			return err
		}
		err = adapterFile.initFile()
		if err != nil {
			return err
		}
	}

	return nil
}

//slice file by size, if maxSize < fileSize, rename file is file_size_maxSize_time.log and recreate file
func (adapterFile *AdapterFile) sliceByFileSize() error {
	fileWrite := adapterFile.write
	filename := adapterFile.config.Filename
	filenameSuffix := path.Ext(filename)
	maxSize := adapterFile.config.MaxSize
	nowSize, _ := adapterFile.getFileSize(filename)
	randStr := utils.NewMisc().RandString(4)

	if nowSize >= maxSize {
		//close file handle
		fileWrite.writer.Close()
		oldFilename := strings.Replace(filename, filenameSuffix, "", 1) + "_size_"+strconv.FormatInt(maxSize, 10)+"_"+randStr+filenameSuffix
		err := os.Rename(filename, oldFilename)
		if err != nil {
			return err
		}
		err = adapterFile.initFile()
		if err != nil {
			return err
		}
	}

	return nil
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