# go-logger
一个简单的 golang 日志工具包  
[English document](/README_EN.md)  

# 功能
- 支持同时输出到 console, file, url 
- 命令行输出字体可带颜色
- 文件输出支持根据 文件大小，文件行数，日期三种方式切分
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


# 使用例子

- ### example

```
import (
    "go-logger"
)
func main()  {
    logger := go_logger.NewLogger()
    
    // 默认已经添加了 console 的输出，默认不显示颜色，如果需要修改，先删除掉 console
    logger.Detach("console")
    
    // 配置 console
    console := &go_logger.ConsoleConfig{
        Color: true, // 文字是否显示颜色 
        JsonFormat: true, // 是否格式化成 json 字符串
    }
    // 添加输出到命令行
    logger.Attach("console", go_logger.NewConfigConsole(console))
    
    // 配置 file
    fileConfig := &go_logger.FileConfig{
        Filename : "./test.log", // 文件名
        MaxSize : 1024 * 1024,  // 文件最大 ，默认 0 不限制
        MaxLine : 100000, // 文件最多多少行，默认 0 不限制
        DateSlice : "d", // 按日期切分文件，支持 "y"(年), "m"(月), "d"(日), "h"(小时), 默认 "" 不限制
        JsonFormat: true, // 写入文件数据是否 json 格式化
    }
    logger.Attach("file", go_logger.NewConfigFile(fileConfig))
    
    // 设置日志级别
    logger.SetLevel(go_logger.LOGGER_LEVEL_DEBUG)
    // 设置为异步，默认是同步方式输出
    logger.SetAsync()

    logger.Emergency("this is a emergency log!")
    logger.Alert("this is a alert log!")
    logger.Critical("this is a critical log!")
    logger.Error("this is a error log!")
    logger.Warning("this is a warning log!")
    logger.Notice("this is a notice log!")
    logger.Info("this is a info log!")
    logger.Debug("this is a debug log!")

    // 如果设置为异步，最后必须调用 flush 方法确保所有的日志都输出完
    logger.Flush()
}
```
- ### console adapter
```
// 配置 console
console := &go_logger.ConsoleConfig{
    Color: true, // 文字是否显示颜色 
    JsonFormat: true, // 是否格式化成 json 字符串
}
// 添加
logger.Attach("console", go_logger.NewConfigConsole(console))
```
#### console 文字带颜色效果
![image](https://github.com/phachon/go-logger/blob/master/_example/images/console.png)

- ### file adapter

```
fileConfig := &go_logger.FileConfig{
    Filename : "./test.log", // 文件名
    MaxSize : 1024 * 1024,  // 文件最大 ，默认 0 不限制
    MaxLine : 100000, // 文件最多多少行，默认 0 不限制
    DateSlice : "d", // 按日期切分文件，支持 "y"(年), "m"(月), "d"(日), "h"(小时), 默认 "" 不限制
    JsonFormat: true, // 写入文件数据是否 json 格式化
}
logger.Attach("file", go_logger.NewConfigFile(fileConfig))
```

- ### api adapter

```
apiConfig := &go_logger.ApiConfig{
    Url: "http://127.0.0.1:8081/index.php", // 请求的 url  地址,不能为空
    Method: "GET", // 请求方式 GET, POST
    Headers: map[string]string{}, // request header
    IsVerify: false, // 是否验证 url 请求返回 http code
    VerifyCode: 0, // 如果 IsVerify 为 true, 需要验证的成功的 http code 码, 不能为 0
}
logger.Attach("api", go_logger.NewConfigApi(apiConfig))
```

## 参考
beego/logs : github.com/astaxie/beego/logs

## 反馈

欢迎提交意见和代码，联系信息 phachon@163.com

## License

MIT

谢谢
---
Create By phachon@163.com
