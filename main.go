package main

import (
	"fmt"

	"github.com/Ajasf444/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	s := state{cfg: &cfg}
}
