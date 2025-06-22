package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sangnt1552314/digimontex/internal/app"
)

func main() {
	// Setup logging
	if err := os.MkdirAll("storage/logs", 0755); err != nil {
		panic(fmt.Errorf("failed to create logs directory: %w", err))
	}

	logFile, err := os.OpenFile("storage/logs/develop.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Initialize the application
	application := app.NewApp()

	if err := application.Run(); err != nil {
		panic(err)
	}
}
