package main

import (
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
)

func main() {
	logger, err := common.NewLogger("")
	if err != nil {
		panic(err)
	}

	defer logger.Close()

	logger.LogI(msg.I2000)
}
