# go-logger
A simple but powerful golang log Toolkit  

[![Sourcegraph](https://sourcegraph.com/github.com/phachon/go-logger/-/badge.svg)](https://sourcegraph.com/github.com/phachon/go-logger?badge)
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/phachon/go-logger) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phachon/go-logger/master/LICENSE)

[中文文档](/README_CN.md)

# Feature
- Support at the same time to console, file, URL
- console output fonts can be colored with
- File output supports three types of segmentation based on the size of the file, the number of file lines, and the date.
- file output support is saved to different files at the log level.
- Two ways of writing to support asynchronous and synchronous
- Support json format output
- The code is designed to be extensible, and you can design your own adapter as needed

# Install

```
go get github.com/phachon/go-logger
go get ./...
```

# Requirement
go 1.8

# Support outputs
- console  // write console
- file     // write file
- api      // http request url
- ...


# Quick Used

- sync

```
import (
    "github.com/phachon/go-logger"
)
func main()  {
    logger := go_logger.NewLogger()

    logger.Info("this is a info log!")
    logger.Errorf("this is a error %s log!", "format")
}
```

- async

```
import (
    "github.com/phachon/go-logger"
)
func main()  {
    logger := go_logger.NewLogger()
    logger.SetAsync()

    logger.Info("this is a info log!")
    logger.Errorf("this is a error %s log!", "format")

    // Flush must be called before the end of the process
    logger.Flush()
}
```

- Multiple output

```
import (
    "github.com/phachon/go-logger"
)
func main()  {
    logger := go_logger.NewLogger()

    logger.Detach("console")

    // console adapter config
    consoleConfig := &go_logger.ConsoleConfig{
        Color: true, // Does the text display the color
        JsonFormat: true, // Whether or not formatted into a JSON string
        Format: "", // JsonFormat is false, logger message output to console format string
    }
    // add output to the console
    logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)

    // file adapter config
    fileConfig := &go_logger.FileConfig {
        Filename : "./test.log", // The file name of the logger output, does not exist automatically
        // If you want to separate separate logs into files, configure LevelFileName parameters.
        LevelFileName : map[int]string {
            logger.LoggerLevel("error"): "./error.log",    // The error level log is written to the error.log file.
            logger.LoggerLevel("info"): "./info.log",      // The info level log is written to the info.log file.
            logger.LoggerLevel("debug"): "./debug.log",    // The debug level log is written to the debug.log file.
        },
        MaxSize : 1024 * 1024,  // File maximum (KB), default 0 is not limited
        MaxLine : 100000, // The maximum number of lines in the file, the default 0 is not limited
        DateSlice : "d",  // Cut the document by date, support "Y" (year), "m" (month), "d" (day), "H" (hour), default "no".
        JsonFormat: true, // Whether the file data is written to JSON formatting
        Format: "", // JsonFormat is false, logger message written to file format string
    }
    // add output to the file
    logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)


    logger.Info("this is a info log!")
    logger.Errorf("this is a error %s log!", "format")
}
```

## Console text with color effect
![image](https://github.com/phachon/go-logger/blob/master/_example/images/console.png)

## Customize Format output

### Logger Message

| Field | Alias |Type  | Comment | Example |
|-------|-------|------|---------|----------|
| Timestamp | timestamp | int64 | unix timestamp| 1521791201 |
| TimestampFormat | timestamp_format| string | timestamp format | 2018-3-23 15:46:41|
| Millisecond | millisecond | int64 | millisecond | 1524472688352 |
| MillisecondFormat | millisecond_format| string | millisecond_format | 2018-3-23 15:46:41.970 |
| Level | level| int | logger level |  1  |
| LevelString | level_string | string | logger level string | Error |
| Body | body | string | logger message body | this is a info log |
| File | file | string | Call the file of the logger | main.go |
| Line | line | int | The number of specific lines to call logger |64|
| Function | function| string | The function name to call logger  | main.main |

>> If you want to customize the format of the log output ?

**config format**:
```
consoleConfig := &go_logger.ConsoleConfig{
    Format: "%millisecond_format% [%level_string%] %body%",
}
fileConfig := &go_logger.FileConfig{
    Format: "%millisecond_format% [%level_string%] %body%",
}
```
**output**:
```
2018-03-23 14:55:07.003 [Critical] this is a critical log!
```

>> You can customize the format, Only needs to be satisfied Format: "%Logger Message Alias%"

## More adapter examples
- [console](./_example/console.go)
- [file](./_example/file.go)
- [api](./_example/api.go)


## Benchmark

system: Linux Mint 18.2 Sonya  
cpu(s): 4  
model name: Intel(R) Core(TM) i5-3210M CPU @ 2.50GHz  
memery: 4G  

```
BenchmarkLoggerConsoleText          500000             11375 ns/op             672 B/op         15 allocs/op
BenchmarkLoggerConsoleText-2        500000             11345 ns/op             672 B/op         15 allocs/op
BenchmarkLoggerConsoleText-4        500000              9897 ns/op             672 B/op         15 allocs/op
```

```
BenchmarkLoggerConsoleAsyncText     500000              9323 ns/op             672 B/op         15 allocs/op
BenchmarkLoggerConsoleAsyncText-2   500000              9087 ns/op             672 B/op         15 allocs/op
BenchmarkLoggerConsoleAsyncText-4   500000             10685 ns/op             672 B/op         15 allocs/op
```

```
BenchmarkLoggerConsoleJson          200000             30918 ns/op            2048 B/op         10 allocs/op
BenchmarkLoggerConsoleJson-2        200000             33153 ns/op            2048 B/op         10 allocs/op
BenchmarkLoggerConsoleJson-4        200000             30918 ns/op            2048 B/op         10 allocs/op
```

```
BenchmarkLoggerFileText             300000             14083 ns/op             912 B/op         21 allocs/op
BenchmarkLoggerFileText-2           200000             21159 ns/op             912 B/op         21 allocs/op
BenchmarkLoggerFileText-4           200000             23776 ns/op             912 B/op         21 allocs/op
```

```
BenchmarkLoggerFileAsyncText        300000             13956 ns/op             912 B/op         21 allocs/op
BenchmarkLoggerFileAsyncText-2      300000             16124 ns/op             912 B/op         21 allocs/op
BenchmarkLoggerFileAsyncText-4      300000             18641 ns/op             912 B/op         21 allocs/op
```

```
BenchmarkLoggerFileJson             200000             15472 ns/op            1968 B/op         15 allocs/op
BenchmarkLoggerFileJson-2           200000             22523 ns/op            1968 B/op         15 allocs/op
BenchmarkLoggerFileJson-4           200000             25596 ns/op            1968 B/op         15 allocs/op
```

## Reference
beego/logs : github.com/astaxie/beego/logs

## Feedback

Welcome to submit comments and code, contact information phachon@163.com

## License

MIT

Thanks
---------
Create By phachon@163.com
