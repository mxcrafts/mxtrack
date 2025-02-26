package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"runtime/debug"

	"github.com/mxcrafts/ltrack/cmd/options"
	"github.com/mxcrafts/ltrack/internal/config"
	"github.com/mxcrafts/ltrack/internal/monitor/file"
	"github.com/mxcrafts/ltrack/internal/monitor/network"
	exec "github.com/mxcrafts/ltrack/internal/monitor/syscall"
	"github.com/mxcrafts/ltrack/pkg/logger"
	"github.com/mxcrafts/ltrack/pkg/version"
)

func main() {

	// Defer panic handler
	defer func() {
		if r := recover(); r != nil {
			logger.Global.Error("Program encountered a critical error",
				"error", r)
			os.Exit(1)
		}
	}()

	// Parse command line options
	opts, err := options.Parse()
	if err != nil {
		fmt.Printf("Error parsing options: %v\n", err)
		os.Exit(1)
	}

	// Validate options
	if err := opts.Validate(); err != nil {
		fmt.Printf("Invalid options: %v\n", err)
		os.Exit(1)
	}

	// Defer panic handler
	defer func() {
		if r := recover(); r != nil {
			logger.Global.Error("Program encountered a critical error",
				"error", r,
				"stack", string(debug.Stack()))
			os.Exit(1)
		}
	}()

	// Load configuration
	config, err := config.Load(opts.ConfigPath)
	if err != nil {
		logger.Global.Error("Load Configuration Failed",
			"path", opts.ConfigPath,
			"error", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(&config.Log); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// Log startup message
	logger.Global.Info("ltrack Started",
		"version", version.Version)

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create wait group
	var wg sync.WaitGroup

	// Create file monitor
	if config.FileMonitor.Enabled {
		monitor, err := file.NewMonitor()
		if err != nil {
			logger.Global.Error("Create File Monitor Failed",
				"error", err)
			os.Exit(1)
		}

		// Type assertion
		if fileMonitor, ok := monitor.(*file.Monitor); ok {
			// Add monitored directories
			for _, dir := range config.FileMonitor.Directories {
				if err := fileMonitor.AddMonitoredDir(dir); err != nil {
					logger.Global.Error("Add Monitored Directory Failed",
						"dir", dir,
						"error", err)
					continue
				}
			}
		}

		// Start monitoring
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := monitor.Start(ctx); err != nil {
				logger.Global.Error("Start File Monitor Failed",
					"error", err)
				return
			}
			<-ctx.Done()
			monitor.Stop(ctx)
		}()
	} else {
		logger.Global.Info("File Monitor Disabled")
	}

	// Create exec monitor
	if config.ExecMonitor.Enabled {
		logger.Global.Info("Initializing Exec Monitor...",
			"watch_commands", config.ExecMonitor.WatchCommands)

		monitor, err := exec.NewMonitor(config)
		if err != nil {
			logger.Global.Error("Create Exec Monitor Failed",
				"error", err)
			os.Exit(1)
		}
		logger.Global.Info("Exec Monitor Created Successfully!")

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := monitor.Start(ctx); err != nil {
				logger.Global.Error("Start Exec Monitor Failed",
					"error", err)
				return
			}

			logger.Global.Info("Exec Monitor Running...",
				"watch_commands", config.ExecMonitor.WatchCommands)

			<-ctx.Done()
			logger.Global.Info("Stopping exec monitor...")
			monitor.Stop(ctx)
			logger.Global.Info("Exec Monitor Stopped Successfully!")
		}()
	} else {
		logger.Global.Info("Exec Monitor Disabled")
	}

	// Create network monitor
	if config.NetworkMonitor.Enabled {
		logger.Global.Info("Initializing network monitor...",
			"ports", config.NetworkMonitor.Ports,
			"protocols", config.NetworkMonitor.Protocols)

		monitor, err := network.NewMonitor(config)
		if err != nil {
			logger.Global.Error("Create Network Monitor Failed",
				"error", err)
			os.Exit(1)
		}
		logger.Global.Info("Network Monitor Created Successfully!")

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := monitor.Start(ctx); err != nil {
				logger.Global.Error("Start Network Monitor Failed",
					"error", err, "stack", debug.Stack())
				return
			}

			logger.Global.Info("Network Monitor Running...",
				"monitored_ports", config.NetworkMonitor.Ports,
				"monitored_protocols", config.NetworkMonitor.Protocols)

			<-ctx.Done()
			logger.Global.Info("Stopping Network Monitor...")
			monitor.Stop(ctx)
			logger.Global.Info("Network Monitor Stopped Successfully!")
		}()
	} else {
		logger.Global.Info("Network Monitor Disabled")
	}

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-sigChan
	logger.Global.Info("Received Signal, Preparing Exit",
		"signal", sig.String())

	// Cancel context
	cancel()

	// Wait for all goroutines to complete with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Global.Info("Program Exited Normally!")
	case <-time.After(5 * time.Second):
		logger.Global.Warn("Program Exit Timed Out!")
	}
}
