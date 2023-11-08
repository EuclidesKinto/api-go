package main

import (
	"api-go/configs"
	"fmt"
)

func main() {
	config, _ := configs.LoadConfig(".")
	fmt.Println(config.DBDriver)
}
