package main

import (
	"go-logger"
)

func main()  {

	logger := go_logger.NewLogger()

	logger.Attach("console", map[string]interface{}{
		"color": false,
	})
	logger.Attach("file", map[string]interface{}{
		"filename": "test.log",
	})

	//logger.SetSync(false)

	logger.Emergency("test1, test, test")
	logger.Emergency("test2, test, test")
	logger.Emergency("test3, test, test")
	logger.Emergency("test4, test, test")
	logger.Emergency("test5, test, test")
	logger.Emergency("test6, test, test")

	//logger.Flush()
}