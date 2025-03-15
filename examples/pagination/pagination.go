package main

import (
	"context"
	"github.com/jordanst3wart/up-client/up"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {
	value, exists := os.LookupEnv("UP_TOKEN")
	if !exists {
		log.Fatal("UP_TOKEN environment variable not set")
	}
	client := up.NewClient(value, nil)
	opts := &up.ListAccountsOptions{
		AccountType: up.AccountTypeTransactional,
	}
	accounts, _, err := client.Accounts.List(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	if len(accounts.Data) != 1 {
		log.Fatal("Expected 1 transaction account")
	}
	previousWeek := time.Now().AddDate(0, 0, -7)
	now := time.Now()
	// iterates over paginated results
	list, _, err := client.Transactions.ListByAccount(context.TODO(), accounts.Data[0].ID, &up.ListTransactionsOptions{
		Since: &previousWeek,
		Until: &now,
	})
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("All transactions for transactions account", "transactions", list)
}
