package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/arribada/smart-camera/pkg/capturer"
	"github.com/arribada/smart-camera/pkg/uploader"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"
	"github.com/oklog/run"
)

func main() {
	logger := NewLogger()

	_env, err := os.Open(".env")
	if err != nil {
		// When running inside containers the .env file doesn't exist.
		if os.IsNotExist(err) {
			level.Warn(logger).Log("msg", ".env file doesn't exist so using the existing env variables")
		} else {
			ExitOnError(logger, err, "opening the env file")
		}
	}

	env, err := ioutil.ReadAll(_env)
	ExitOnError(logger, err, "reading the env file")

	rr := bytes.NewReader(env)
	envMap, err := godotenv.Parse(rr)
	ExitOnError(logger, err, "parsing the env file")

	// Copied from the godotenv source code.
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] {
			os.Setenv(key, value)
		}
	}

	globalCtx, globalShutdown := context.WithCancel(context.Background())
	defer globalShutdown()

	// We define our run groups here.
	var g run.Group
	// Run groups.
	{
		// Handle interupts.
		g.Add(run.SignalHandler(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM))

		capt, err := capturer.New(
			globalCtx,
			logger,
			capturer.Config{
				DataFolder: "data",
				Interval:   10 * time.Second,
			},
		)
		ExitOnError(logger, err, "creating component:"+capturer.ComponentName)

		g.Add(func() error {
			capt.Start()
			level.Info(logger).Log("msg", "shutdown complete", "component", capturer.ComponentName)
			return nil
		}, func(error) {
			capt.Stop()
		})

		upl, err := uploader.New(
			globalCtx,
			logger,
			uploader.Config{
				DataFolder: "data",
				Interval:   10 * time.Minute,
			},
		)
		ExitOnError(logger, err, "creating component:"+uploader.ComponentName)

		g.Add(func() error {
			upl.Start()
			level.Info(logger).Log("msg", "shutdown complete", "component", uploader.ComponentName)
			return nil
		}, func(error) {
			upl.Stop()
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "main exited with error", "err", err)
		return
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
