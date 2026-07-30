package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	account "github.com/curtrube/aws-account-service/internal/aws"
	"github.com/curtrube/aws-account-service/internal/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	dbPool, err := postgres.New(ctx, "")
	if err != nil {
		log.Fatalf("creating postgres pool: %v", err)
	}
	defer dbPool.Close()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	orgClient := organizations.NewFromConfig(cfg)

	ticker := time.NewTicker(10 * time.Second)

	conn, err := dbPool.Acquire(ctx)
	defer conn.Release()
	store := postgres.NewAccountStore(conn)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("aws account poller shutting down")
				return

			case <-ticker.C:
				fmt.Println("polling aws account metadata")
				accounts, err := account.ListAccounts(ctx, orgClient)
				fmt.Println(accounts)
				if err != nil {
					log.Fatalf("listing aws accounts: %v", err)
				}
				for _, account := range accounts {
					if _, err := store.Create(ctx, account); err != nil {
						fmt.Printf("err creating account: %v", err)
					}
				}
			}
		}
	}()

	<-ctx.Done()
	ticker.Stop()
	wg.Wait()
	fmt.Println("graceful shutdown")
}
