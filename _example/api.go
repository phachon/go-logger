package main

import (
	"github.com/phachon/go-logger"
)

func main() {

	logger := go_logger.NewLogger()

	apiConfig := &go_logger.ApiConfig{
		Url:        "http://127.0.0.1:8081/index.php",
		Method:     "GET",
		Headers:    map[string]string{},
		IsVerify:   false,
		VerifyCode: 0,
	}
	logger.Attach("api", go_logger.LOGGER_LEVEL_DEBUG, apiConfig)
	logger.SetAsync()

	logger.Emergency("this is a emergency log!")
	logger.Alert("this is a alert log!")

	logger.Flush()
}
