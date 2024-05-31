package main

import (
	"fmt"
	"os"

	"github.com/uandersonricardo/masterclass-go/internal"
)

func main() {
	fmt.Println("Starting server...")

	server := internal.NewGrpcServer(":8080")
	err := server.Start()

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
