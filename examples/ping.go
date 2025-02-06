package examples

import (
	"context"
	"fmt"
	"github.com/jordanstewart/up-client/up"
	"log"
)

func ping() {
	client := up.NewClient("your-token-here", nil)
	pingResp, _, err := client.Utility.Ping(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API Status: %s\n", pingResp.Meta.StatusEmoji)
}
