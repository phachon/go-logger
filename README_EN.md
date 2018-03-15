# go-logger
A simple golang log Toolkit

[中文文档](/README.md)

# Feature
- Support at the same time to console, file, URL
- console output fonts can be colored with
- File output supports three types of segmentation based on the size of the file, the number of file lines, and the date.
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
- console  //write console
- file     //write file
- api      // http request url
- ...


# Use

- ### example

```
import (
    "github.com/phachon/go-logger"
)
func main()  {
    logger := go_logger.NewLogger()

    // The default has been added to the output of console, and the default does not display the color. If you need to modify it, delete the console first
    logger.Detach("console")

    // config console
    console := &go_logger.ConsoleConfig{
        Color: true, // text show color
        JsonFormat: true, // json format
        ShowFileLine: true,  // when JsonFormat is false, console show file line, default false 
    }
    // attach console to outputs
    logger.Attach("console", go_logger.NewConfigConsole(console))

    // config file
    fileConfig := &go_logger.FileConfig{
        Filename : "./test.log", // filename
        MaxSize : 1024 * 1024,  // max file size(kb)
        MaxLine : 100000, // max file line
        DateSlice : "d", // slice file by date, support "y", "m", "d", "h", default "" not slice
        JsonFormat: true, // json format
    }
    logger.Attach("file", go_logger.NewConfigFile(fileConfig))

    // set logger level
    logger.SetLevel(go_logger.LOGGER_LEVEL_DEBUG)
    // Set to asynchronous, default is synchronous output
    logger.SetAsync()

    logger.Emergency("this is a emergency log!")
    logger.Alert("this is a alert log!")
    logger.Critical("this is a critical log!")
    logger.Error("this is a error log!")
    logger.Warning("this is a warning log!")
    logger.Notice("this is a notice log!")
    logger.Info("this is a info log!")
    logger.Debug("this is a debug log!")
    
    logger.Emergencyf("this is a emergency %d log!", 10)
    logger.Alertf("this is a alert %s log!", "format")
    logger.Criticalf("this is a critical %s log!", "format")
    logger.Errorf("this is a error %s log!", "format")
    logger.Warningf("this is a warning %s log!", "format")
    logger.Noticef("this is a notice %s log!", "format")
    logger.Infof("this is a info %s log!", "format")
    logger.Debugf("this is a debug %s log!", "format")
    
    // If set to asynchronous, the flush method must finally be invoked to ensure that all the logs are out
    logger.Flush()
}
```
- ### console adapter
```
// config console
console := &go_logger.ConsoleConfig{
    Color: true, // text show color
    JsonFormat: true, // json format
    ShowFileLine: true,  // when JsonFormat is false, console show file line, default false 
}
// attach
logger.Attach("console", go_logger.NewConfigConsole(console))
```
#### console color preview
![image](https://github.com/phachon/go-logger/blob/master/_example/images/console.png)

- ### file adapter

```
fileConfig := &go_logger.FileConfig{
    Filename : "./test.log", // filename
    MaxSize : 1024 * 1024,  // max file size(kb)
    MaxLine : 100000, // max file line
    DateSlice : "d", // slice file by date, support "y", "m", "d", "h", default "" not slice
    JsonFormat: true, // json format
}
logger.Attach("file", go_logger.NewConfigFile(fileConfig))
```

- ### api adapter

```
apiConfig := &go_logger.ApiConfig{
    Url: "http://127.0.0.1:8081/index.php", //request url address, not empty
    Method: "GET", //request method GET or POST
    Headers: map[string]string{},  //request headers, default empty
    IsVerify: false, //response is verify code, default false
    VerifyCode: 0, //verify code value, if isVerify is true, verifyCode is not be 0
}
logger.Attach("api", go_logger.NewConfigApi(apiConfig))
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
