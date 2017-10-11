package main

import (
	"go-logger"
)

func main()  {

	logger := go_logger.NewLogger()

	logger.Attach("api", map[string]interface{}{
		"url": "http://127.0.0.1:8081/test.php",
		"method": "POST",//GET,POST
		"headers": map[string]string{},
		"isVerify": true,
		"verifyCode": 200,
	})
	logger.SetLevel(go_logger.LOGGER_LEVEL_DEBUG)
	logger.SetAsync()

	logger.Emergency("this is a emergency log!")
	logger.Alert("this is a alert log!")

	logger.Flush()
}