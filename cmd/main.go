package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/biswaone/mongo-oplog-to-sql/config"
	"github.com/biswaone/mongo-oplog-to-sql/database"
	"github.com/biswaone/mongo-oplog-to-sql/internal/app/oplog"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mongo-oplog-to-sql",
	Short: "Transfer data from mongo to sql db",
	Long:  `A Tool to move data from mongo db to sql based database`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		ctx, cancel := context.WithCancel(context.Background())
		// Handle interrupt signal
		handleInterruptSignal(cancel)

		client := database.NewMongoClient(ctx, cfg)
		oplogCh, err := oplog.GetOplog(ctx, client, "local")
		if err != nil {
			log.Fatal("Error getting oplog:", err)
		}

		for entry := range oplogCh {
			if entry.Namespace == "employee.employees" && entry.Operation == "i" {
				queries := oplog.GenerateSqlDocument("employees", entry.Doc, "INSERT")
				fmt.Println(queries)
			}
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func handleInterruptSignal(cancel context.CancelFunc) {
	// Create an interrupt channel to listen for the interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		// Cancel the context to signal the shutdown
		cancel()
	}()
}
