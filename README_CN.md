# go-logger
一个简单而强大的 golang 日志工具包

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/phachon/go-logger) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phachon/go-logger/master/LICENSE)

[English document](/README.md)  

# 功能
- 支持同时输出到 console, file, url 
- 命令行输出字体可带颜色
- 文件输出支持根据 文件大小，文件行数，日期三种方式切分
- 文件输出支持根据日志级别分别保存到不同的文件
- 支持异步和同步两种方式写入
- 支持 json 格式化输出
- 代码设计易扩展，可根据需要设计自己的 adapter

# 安装使用

```
go get github.com/phachon/go-logger
go get ./...
```
# 环境需要
go 1.8

# 支持输出
- console  // 输出到命令行
- file     // 文件
- api      // http url 接口
- ...

# 快速使用

- 同步方式

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

- 异步方式

```
import (
    "github.com/phachon/go-logger"
)
func main()  {
    logger := go_logger.NewLogger()
    logger.SetAsync()

    logger.Info("this is a info log!")
    logger.Errorf("this is a error %s log!", "format")

    // 程序结束前必须调用 Flush
    logger.Flush()
}
```

- 多个输出

```
import (
    "github.com/phachon/go-logger"
)
func main()  {
    logger := go_logger.NewLogger()

    logger.Detach("console")

    // 命令行输出配置
    consoleConfig := &go_logger.ConsoleConfig{
        Color: true, // 命令行输出字符串是否显示颜色
        JsonFormat: true, // 命令行输出字符串是否格式化
        Format: "" // 如果输出的不是 json 字符串，JsonFormat: false, 自定义输出的格式
    }
    // 添加 console 为 logger 的一个输出
    logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)

    // 文件输出配置
    fileConfig := &go_logger.FileConfig {
        Filename : "./test.log", // 日志输出文件名，不自动存在
        // 如果要将单独的日志分离为文件，请配置LealFrimeNem参数。
        LevelFileName : map[int]string {
            logger.LoggerLevel("error"): "./error.log",    // Error 级别日志被写入 error .log 文件
            logger.LoggerLevel("info"): "./info.log",      // Info 级别日志被写入到 info.log 文件中
            logger.LoggerLevel("debug"): "./debug.log",    // Debug 级别日志被写入到 debug.log 文件中
        },
        MaxSize : 1024 * 1024,  // 文件最大值（KB），默认值0不限
        MaxLine : 100000, // 文件最大行数，默认 0 不限制
        DateSlice : "d",  // 文件根据日期切分， 支持 "Y" (年), "m" (月), "d" (日), "H" (时), 默认 "no"， 不切分
        JsonFormat: true, // 写入文件的数据是否 json 格式化
        Format: "" // 如果写入文件的数据不 json 格式化，自定义日志格式
    }
    // 添加 file 为 logger 的一个输出
    logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)


    logger.Info("this is a info log!")
    logger.Errorf("this is a error %s log!", "format")
}
```

## 命令行下的文本带颜色效果
![image](https://github.com/phachon/go-logger/blob/master/_example/images/console.png)

## 自定义格式化输出

Logger Message

| 字段 | 别名 |类型  | 说明 | 例子 |
|-------|-------|------|---------|----------|
| Timestamp | timestamp | int64 | Unix时间戳 | 1521791201 |
| TimestampFormat | timestamp_format| string | 时间戳格式化字符串 | 2018-3-23 15:46:41|
| Millisecond | millisecond | int64 | 毫秒时间戳 | 1524472688352 |
| MillisecondFormat | millisecond_format| string | 毫秒时间戳格式化字符串 | 2018-3-23 15:46:41.970 |
| Level | level| int | 日志级别 |  1  |
| LevelString | level_string | string | 日志级别字符串 | Error |
| Body | body | string | 日志内容 | this is a info log |
| File | file | string | 调用本次日志输出的文件名 | main.go |
| Line | line | int | 调用本次日志输出的方法 |64|
| Function | function| string | 调用本次日志输出的方法名  | main.main |

>> 你想要自定义日志输出格式 ?

**配置 Format 参数**:
```
consoleConfig := &go_logger.ConsoleConfig{
    Format: "%millisecond_format% [%level_string%] %body%"
}
fileConfig := &go_logger.FileConfig{
    Format: "%millisecond_format% [%level_string%] %body%"
}
```
**输出结果**:
```
2018-03-23 14:55:07.003 [Critical] this is a critical log!
```

>> 你只需要配置参数 Format: "% Logger Message 别名%" 来自定义输出字符串格式

## 更多的 adapter 例子
- [console](./_example/console.go)
- [file](./_example/file.go)
- [api](./_example/api.go)


## 性能测试结果

## 参考
beego/logs : github.com/astaxie/beego/logs

## 反馈

欢迎提交意见和代码，联系信息 phachon@163.com

## License

MIT

谢谢
---
Create By phachon@163.com
