package main

import (
	"os"
	"fmt"
)

func usage() {
	fmt.Println(os.Args[0], " configfile")
	os.Exit(-1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	err := config.Init(os.Args[1])
	if err != nil {
		fmt.Println("config init failed, err:", err)
		os.Exit(-1)
	}

	fmt.Println(config)
	logger.Info(config)
}
