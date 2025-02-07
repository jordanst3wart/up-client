package main

import (
	"context"
	"fmt"
	"github.com/jordanst3wart/up-client/up"
	"log"
	"os"
)

func main() {
	value, exists := os.LookupEnv("UP_TOKEN")
	if !exists {
		log.Fatal("UP_TOKEN environment variable not set")
	}
	client := up.NewClient(value, nil)
	pingResp, _, err := client.Utility.Ping(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API Status: %s\n", pingResp.Meta.StatusEmoji)
}
