package main

import (
	"fmt"

	"github.com/noredis/subscriptions/internal/common/config"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg
	fmt.Println(cfg.App.Env)

	// logger

	// db

	// http

	fmt.Println("Hello, World!")
}
