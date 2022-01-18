// Copyright (c) The Arribada initiative.
// Licensed under the MIT License.

package capturer

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kennygrant/sanitize"
	"github.com/pkg/errors"
)

const ComponentName = "Camera Capture"

type Config struct {
	DataFolder string
	Interval   time.Duration
	LogLevel   string
}

func New(ctx context.Context, logger log.Logger, cfg Config) (*Capturer, error) {
	if cfg.DataFolder == "" {
		return nil, errors.New("upload dir data folder can't be empty")
	}

	if err := os.MkdirAll(cfg.DataFolder, os.ModePerm); err != nil {
		return nil, errors.Wrapf(err, "create destination dir:%+v", cfg.DataFolder)
	}

	ctx, close := context.WithCancel(ctx)
	d := &Capturer{
		logger: log.With(logger, "component", ComponentName),
		config: cfg,
		ctx:    ctx,
		close:  close,
		ticker: time.NewTicker(cfg.Interval),
	}

	return d, nil
}

type Capturer struct {
	logger log.Logger
	config Config
	ctx    context.Context
	close  context.CancelFunc
	ticker *time.Ticker
}

func (self *Capturer) Start() error {
	level.Info(self.logger).Log("msg", "starting")

	for {
		select {
		case <-self.ctx.Done():
			return nil
		case <-self.ticker.C:
			if err := self.capture(); err != nil {
				level.Error(self.logger).Log("msg", "capturing", "err", err)
			}
		}
	}
}

func (self *Capturer) capture() error {
	fName := sanitize.Accents(time.Now().UTC().Format("02 Jan 15:04:05.99") + ".png")
	fName = strings.Replace(fName, " ", "-", -1)
	fPath := filepath.Join("data/", fName)

	cmd := exec.Command("pylepton_capture", fPath)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "running capture script:%v", string(stdout))
	}

	level.Info(self.logger).Log("msg", "capture command run", "path", fPath, "output", string(stdout))
	return nil
}

func (self *Capturer) Stop() {
	self.close()
	level.Info(self.logger).Log("msg", "stopped")
}
