package common

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
)

// IsErr helper function to return bool if err != nil.
func IsErr(err error) bool {
	return err != nil
}

// GetServerAddr returns server address in format host:port.
func GetServerAddr(host string, port int) string {
	return fmt.Sprintf("%v:%v", host, port)
}

// GetConfigPathFromArg returns path to configuration file from command line argument.
func GetConfigPathFromArg() string {
	var path string

	flag.StringVar(&path, "config", "", "Path to configuration file")
	flag.Parse()

	return path
}

// GetNotifyCancelCtx returns context with cancel function for graceful shutdown.
func GetNotifyCancelCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)

	return ctx, cancel
}
