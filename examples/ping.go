package examples

import (
	"fmt"
	"log"
)

func ping() {
	client := upclient.NewClient("your-token-here")
	pingResp, err := client.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API Status: %s\n", pingResp.Meta.StatusEmoji)
}
