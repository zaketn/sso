package main

import (
	"fmt"
	"github.com/zaketn/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
