package main

import (
	"fmt"

	"github.com/Lin-Jiong-HDU/vow/internal/config"
)

func main() {
	var cfg *config.Config

	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	fmt.Printf("Config loaded successfully: %+v\n", cfg)
}
