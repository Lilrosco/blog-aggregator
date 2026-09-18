package main

import (
	"fmt"

	"github.com/Lilrosco/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		fmt.Println(err)
	}

	username := "lilrosco"
	config.SetUser(&username, &cfg)
	cfg, err = config.Read()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", cfg)
}
