package main

import (
	"fmt"

	"github.com/fummbly/gatorcli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
	}

	cfg.SetUser("tom")

	newCfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading changed config: %v", err)
	}

	fmt.Println(newCfg)

}
