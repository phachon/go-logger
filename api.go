package go_logger

import (
	"errors"
	"fmt"
	"github.com/phachon/go-logger/utils"
	"reflect"
	"strconv"
)

const API_ADAPTER_NAME = "api"

// adapter api
type AdapterApi struct {
	config *ApiConfig
}

// api config
type ApiConfig struct {

	// request url adddress
	Url string

	// request method
	// GET, POST
	Method string

	// request headers
	Headers map[string]string

	// is verify response code
	IsVerify bool

	// verify response http code
	VerifyCode int
}

func (ac *ApiConfig) Name() string {
	return API_ADAPTER_NAME
}

func NewAdapterApi() LoggerAbstract {
	return &AdapterApi{}
}

func (adapterApi *AdapterApi) Init(apiConfig Config) error {

	if apiConfig.Name() != API_ADAPTER_NAME {
		return errors.New("logger api adapter init error, config must ApiConfig")
	}

	vc := reflect.ValueOf(apiConfig)
	ac := vc.Interface().(*ApiConfig)
	adapterApi.config = ac

	if adapterApi.config.Url == "" {
		return errors.New("config Url cannot be empty!")
	}
	if adapterApi.config.Method != "GET" && adapterApi.config.Method != "POST" {
		return errors.New("config Method must one of the 'GET', 'POST'!")
	}
	if adapterApi.config.IsVerify && (adapterApi.config.VerifyCode == 0) {
		return errors.New("config if IsVerify is true, VerifyCode cannot be 0!")
	}
	return nil
}

func (adapterApi *AdapterApi) Write(loggerMsg *loggerMessage) error {

	url := adapterApi.config.Url
	method := adapterApi.config.Method
	isVerify := adapterApi.config.IsVerify
	verifyCode := adapterApi.config.VerifyCode
	headers := adapterApi.config.Headers

	loggerMap := map[string]string{
		"timestamp":          strconv.FormatInt(loggerMsg.Timestamp, 10),
		"timestamp_format":   loggerMsg.TimestampFormat,
		"millisecond":        strconv.FormatInt(loggerMsg.Millisecond, 10),
		"millisecond_format": loggerMsg.MillisecondFormat,
		"level":              strconv.Itoa(loggerMsg.Level),
		"level_string":       loggerMsg.LevelString,
		"body":               loggerMsg.Body,
		"file":               loggerMsg.File,
		"line":               strconv.Itoa(loggerMsg.Line),
		"function":           loggerMsg.Function,
	}

	var err error
	var code int
	if method == "GET" {
		_, code, err = utils.NewMisc().HttpGet(url, loggerMap, headers, 0)
	} else {
		_, code, err = utils.NewMisc().HttpPost(url, loggerMap, headers, 0)
	}
	if err != nil {
		return err
	}
	if isVerify && (code != verifyCode) {
		return fmt.Errorf("%s", "request "+url+" faild, code="+strconv.Itoa(code))
	}

	return nil
}

func (adapterApi *AdapterApi) Flush() {

}

func (adapterApi *AdapterApi) Name() string {
	return API_ADAPTER_NAME
}

func init() {
	Register(API_ADAPTER_NAME, NewAdapterApi)
}
