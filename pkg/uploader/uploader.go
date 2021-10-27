// Copyright (c) The Arribada initiative.
// Licensed under the MIT License.

package uploader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

const ComponentName = "File Uploader"

type Config struct {
	DataFolder string
	Interval   time.Duration
	LogLevel   string
}

func New(ctx context.Context, logger log.Logger, cfg Config) (*Uploader, error) {
	mClt, err := minio.New("storage.googleapis.com", &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("STORAGE_KEY_ID"), os.Getenv("STORAGE_KEY_SECRET"), ""),
		Secure: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "initalize the storage upload client")
	}

	ctx, close := context.WithCancel(ctx)
	d := &Uploader{
		logger:       log.With(logger, "component", ComponentName),
		config:       cfg,
		ctx:          ctx,
		close:        close,
		uploadClient: mClt,
		ticker:       time.NewTicker(cfg.Interval),
	}

	return d, nil
}

type Uploader struct {
	logger       log.Logger
	config       Config
	uploadClient *minio.Client
	ctx          context.Context
	close        context.CancelFunc
	ticker       *time.Ticker
}

func (self *Uploader) Start() error {
	level.Info(self.logger).Log("msg", "starting")

	if self.config.DataFolder == "" {
		return errors.New("upload dir data folder can't be empty")
	}

	for {
		select {
		case <-self.ctx.Done():
			return nil
		case <-self.ticker.C:
			if err := self.upload(); err != nil {
				level.Error(self.logger).Log("msg", "uploading", "err", err)
			}
		}
	}
}

func (self *Uploader) upload() error {
	return filepath.Walk(self.config.DataFolder, func(path string, info os.FileInfo, err error) error {
		select {
		case <-self.ctx.Done():
			return nil
		default:
		}
		if err != nil {
			return err
		}
		if !info.IsDir() {
			contentType := "image/png"

			n, err := self.uploadClient.FPutObject(context.Background(), "arribada-smart", path, path, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				return errors.Wrap(err, "put object")
			}
			level.Info(self.logger).Log("msg", "upload complete", "file", path, "metadata", fmt.Sprintf("%+v", n))

			level.Info(self.logger).Log("msg", "deleting transferred file", "path", path)
			return os.Remove(filepath.Join(path))
		}
		return nil
	})
}

func (self *Uploader) Stop() {
	self.close()
	level.Info(self.logger).Log("msg", "stopped")
}
