package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/arribada/smart-camera/pkg/capturer"
	"github.com/arribada/smart-camera/pkg/edgeimpulse"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
)

var cli struct {
	Capture     capturer.Config    `cmd:"" help:"Captures images to directory"`
	EdgeImpulse edgeimpulse.Config `cmd:"" help:"Edge Impulse commands"`
}

func main() {
	kongCtx := kong.Parse(&cli, kong.UsageOnError())

	logger := NewLogger()
	globalCtx, globalShutdown := context.WithCancel(context.Background())
	defer globalShutdown()

	switch kongCtx.Command() {
	case "capture":
		runCapture(globalCtx, logger, cli.Capture)
	case "edge-impulse remove-all":
		runEdgeImpRemoveAll(globalCtx, logger, cli.EdgeImpulse)
	case "edge-impulse upload <file-or-dir>":
		runEdgeImpUpload(globalCtx, logger, cli.EdgeImpulse)
	case "edge-impulse upload-many <file-or-dir>":
		runEdgeImpUploadMany(globalCtx, globalShutdown, logger, cli.EdgeImpulse)
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
	capt, err := capturer.New(ctx, logger, cfg)
	ExitOnError(logger, err, "creating component:"+capturer.ComponentName)
	runInGroup(logger, func() error {
		capt.Start()
		level.Info(logger).Log("msg", "shutdown complete", "component", capturer.ComponentName)
		return nil
	}, func(error) { capt.Stop() })
}

func runInGroup(logger log.Logger, startF func() error, stopF func(error)) {
	var g run.Group
	// Run groups.
	// Handle interupts.
	g.Add(run.SignalHandler(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM))

	g.Add(startF, stopF)

	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "main exited with error", "err", err)
		return
	}
}

func runEdgeImpRemoveAll(ctx context.Context, logger log.Logger, cfg edgeimpulse.Config) {
	if !confirm(logger, "WARNING! Remove all training data for project? Type (Remove)/(n)o: ", "Remove") {
		level.Info(logger).Log("msg", "canceled")
		return
	}
	ei, err := edgeimpulse.New(ctx, logger, cfg)
	ExitOnError(logger, err, "creating component:"+edgeimpulse.ComponentName)
	err = ei.RemoveAll()
	ExitOnError(logger, err, "removeAll failed")
	level.Info(logger).Log("msg", "done")
}

func runEdgeImpUpload(ctx context.Context, logger log.Logger, cfg edgeimpulse.Config) {
	ei, err := edgeimpulse.New(ctx, logger, cfg)
	ExitOnError(logger, err, "creating component:"+edgeimpulse.ComponentName)
	err = ei.Upload(cfg.Upload.FileOrDir, cfg.Upload.Label, cfg.Upload.HMAC)
	ExitOnError(logger, err, "upload failed")
	level.Info(logger).Log("msg", "done")
}

func runEdgeImpUploadMany(ctx context.Context, cancelF func(), logger log.Logger, cfg edgeimpulse.Config) {

	ei, err := edgeimpulse.New(ctx, logger, cfg)
	ExitOnError(logger, err, "creating component:"+edgeimpulse.ComponentName)
	runInGroup(logger, func() error {
		return ei.UploadDir(cfg.UploadMany.FileOrDir, cfg.UploadMany.Label, cfg.UploadMany.HMAC)
	}, func(err error) {
		cancelF()
	})

}

func confirm(logger log.Logger, msg, check string) bool {
	var response string
	fmt.Fprint(os.Stdout, msg)
	_, err := fmt.Scanln(&response)
	ExitOnError(logger, err, "confirm failed")
	return response == check
}
