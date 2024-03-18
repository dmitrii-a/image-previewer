package main

import (
	"context"
	"os"
	"time"

	"github.com/dmitrii-a/image-previewer/internal/common"
	"github.com/dmitrii-a/image-previewer/internal/presentation/http/fiber"
)

func main() {
	common.Config.SetConfigFileSettings(common.GetConfigPathFromArg())

	server := fiber.NewServer()

	ctx, cancel := common.GetNotifyCancelCtx()
	defer cancel()

	go func() {
		<-ctx.Done()

		stopCtx, stopCancel := context.WithTimeout(
			context.Background(), time.Duration(common.Config.Server.ShutdownTimeout)*time.Second,
		)
		defer stopCancel()

		if err := server.Stop(stopCtx); common.IsErr(err) { //nolint: contextcheck
			common.Logger.Error().Msg("failed to stop http server: " + err.Error())
		}
	}()

	common.Logger.Info().Msg("image-previewer service is starting...")

	go func() {
		if err := server.Start(ctx); common.IsErr(err) {
			common.Logger.Error().Msg("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	<-ctx.Done()
}
