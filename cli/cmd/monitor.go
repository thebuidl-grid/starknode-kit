package cmd

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/spf13/cobra"
    "starknode-kit/pkg/monitoring"
)

var monitorCmd = &cobra.Command{
    Use:   "monitor",
    Short: "Launch real-time monitoring dashboard",
    Long:  `Start the terminal-based monitoring dashboard for your Ethereum clients`,
    Run: func(cmd *cobra.Command, args []string) {
        runMonitor()
    },
}

func init() {
    rootCmd.AddCommand(monitorCmd)
}

func runMonitor() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        cancel()
    }()
    
    monitor := monitoring.NewMonitorApp()
    
    fmt.Println("Starting StarkNode Monitor Dashboard...")
    fmt.Println("Press 'q' or ESC to quit")
    
    if err := monitor.Start(ctx); err != nil {
        fmt.Printf("Error running monitor: %v\n", err)
        os.Exit(1)
    }
}