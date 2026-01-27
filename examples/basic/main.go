package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arslanbekov/statusgator-go-client/statusgator"
)

func main() {
	token := os.Getenv("STATUSGATOR_API_TOKEN")
	if token == "" {
		log.Fatal("STATUSGATOR_API_TOKEN environment variable is required")
	}

	client, err := statusgator.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Verify API connectivity
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("API ping failed: %v", err)
	}
	fmt.Println("API connection successful!")

	// List all boards
	fmt.Println("\n=== Boards ===")
	boards, err := client.Boards.ListAll(ctx)
	if err != nil {
		log.Fatalf("Failed to list boards: %v", err)
	}

	for _, board := range boards {
		fmt.Printf("Board: %s (ID: %s)\n", board.Name, board.ID)

		// List monitors for each board
		monitors, err := client.Monitors.ListAll(ctx, board.ID)
		if err != nil {
			log.Printf("  Failed to list monitors: %v", err)
			continue
		}

		for _, monitor := range monitors {
			statusIcon := getStatusIcon(monitor.Status)
			fmt.Printf("  %s %s [%s] - %s\n", statusIcon, monitor.Name, monitor.Type, monitor.Status)
		}
	}

	// List available monitoring regions
	fmt.Println("\n=== Monitoring Regions ===")
	regions, err := client.Regions.List(ctx)
	if err != nil {
		log.Printf("Failed to list regions: %v", err)
	} else {
		for _, region := range regions {
			fmt.Printf("  %s (%s)\n", region.Name, region.Code)
		}
	}

	// List organization users
	fmt.Println("\n=== Users ===")
	users, err := client.Users.List(ctx)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
	} else {
		for _, user := range users {
			fmt.Printf("  %s <%s> - %s\n", user.Name, user.Email, user.Role)
		}
	}
}

func getStatusIcon(status statusgator.MonitorStatus) string {
	switch status {
	case statusgator.MonitorStatusUp:
		return "✓"
	case statusgator.MonitorStatusDown:
		return "✗"
	case statusgator.MonitorStatusWarn:
		return "⚠"
	case statusgator.MonitorStatusMaintenance:
		return "🔧"
	default:
		return "?"
	}
}
