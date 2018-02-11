package go_logger

import (
	"fmt"
	"strconv"
	"go-logger/utils"
)

const API_ADAPTER_NAME = "api"

type AdapterApi struct {
	config map[string]interface{}
}

func NewAdapterApi() LoggerAbstract {
	//api default config
	defaultConfig := map[string]interface{}{
		"url": "",
		"method": "GET",
		"headers": map[string]string{},
		"isVerify": false,
		"verifyCode": 0,
	}
	return &AdapterApi{
		config: defaultConfig,
	}
}

func (adapterApi *AdapterApi) Init(config map[string]interface{}) {
	adapterApi.config = utils.NewMisc().MapIntersect(adapterApi.config, config)
	if adapterApi.config["url"].(string) == "" {
		printError("logger: api adapter config url cannot be empty!")
	}
	if adapterApi.config["method"].(string) != "GET" && adapterApi.config["method"].(string) != "POST" {
		printError("logger: api adapter config method must one of the 'GET', 'POST'!")
	}
	if adapterApi.config["isVerify"].(bool) && (adapterApi.config["verifyCode"] == 0) {
		printError("logger: api adapter config if isVerify is true, verifyCode cannot be 0!")
	}
}

func (adapterApi *AdapterApi) Write(loggerMsg *loggerMessage) error {

	url :=  adapterApi.config["url"].(string)
	method :=  adapterApi.config["method"].(string)
	isVerify :=  adapterApi.config["isVerify"].(bool)
	verifyCode :=  adapterApi.config["verifyCode"].(int)
	headers :=  adapterApi.config["headers"].(map[string]string)

	loggerMap := map[string]string {
		"timestamp": strconv.FormatInt(loggerMsg.Timestamp, 10),
		"timestamp_format": loggerMsg.TimestampFormat,
		"millisecond": strconv.FormatInt(loggerMsg.Millisecond, 10),
		"millisecond_format": loggerMsg.MillisecondFormat,
		"level": strconv.Itoa(loggerMsg.Level),
		"level_string": loggerMsg.LevelString,
		"body": loggerMsg.Body,
		"file": loggerMsg.File,
		"line": strconv.Itoa(loggerMsg.Line),
		"function": loggerMsg.Function,
	}

	var err error
	var code int
	if method == "GET" {
		_, code, err = utils.NewMisc().HttpGet(url, loggerMap, headers, 0)
	}else {
		_, code, err = utils.NewMisc().HttpPost(url, loggerMap, headers, 0)
	}
	if err != nil {
		return err
	}
	if(isVerify && (code != verifyCode)) {
		return fmt.Errorf("%s", "request "+ url +" faild, code=" + strconv.Itoa(code))
	}

	return nil
}

func (adapterApi *AdapterApi) Flush() {

}

func (adapterApi *AdapterApi)Name() string {
	return API_ADAPTER_NAME
}

func init()  {
	Register(API_ADAPTER_NAME, NewAdapterApi)
}