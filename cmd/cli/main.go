package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/arribada/smart-camera/pkg/capturer"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
)

var cli struct {
	Capture capturer.Config `cmd:"" help:"Captures images to directory"`
}

func main() {
	kongCtx := kong.Parse(&cli, kong.UsageOnError())

	logger := NewLogger()
	globalCtx, globalShutdown := context.WithCancel(context.Background())
	defer globalShutdown()

	switch kongCtx.Command() {
	case "capture":
		runCapture(globalCtx, logger, cli.Capture)
	default:
		ExitOnError(logger, errors.Errorf("unknown cmd: %s", kongCtx.Command()), "initialization")
	}
	
	level.Info(logger).Log("msg", "main shutdown complete")
}

const DefaultTimeFormat = "02 Jan 15:04:05.00"

// NewLogger create a new logger.
func NewLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	return log.With(logger, "ts", log.TimestampFormat(func() time.Time { return time.Now().UTC() }, DefaultTimeFormat), "caller", log.Caller(5))
}

func ExitOnError(logger log.Logger, err error, msg string) {
	if err != nil {
		level.Error(logger).Log("msg", msg, "err", err)
		os.Exit(1)
	}
}

func runCapture(ctx context.Context, logger log.Logger, cfg capturer.Config) {
	var g run.Group
	// Run groups.
	{
		// Handle interupts.
		g.Add(run.SignalHandler(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM))

		capt, err := capturer.New(ctx, logger, cfg)
		ExitOnError(logger, err, "creating component:"+capturer.ComponentName)

		g.Add(func() error {
			capt.Start()
			level.Info(logger).Log("msg", "shutdown complete", "component", capturer.ComponentName)
			return nil
		}, func(error) {
			capt.Stop()
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "main exited with error", "err", err)
		return
	}
}
