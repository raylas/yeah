package main

import (
	"log"

	"rymnd.net/yeah/internal/encode"
)

func main() {
	if err := encode.Vendors(); err != nil {
		log.Fatalf("failed to encode vendors: %v", err)
	}
}
