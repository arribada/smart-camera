package edgeimpulse

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
)

const ComponentName = "EdgeImpulse Wrapper"

type Config struct {
	Project   string        `env:"PROJECT" short:"p" help:"Project ID on EdgeImpulse"`
	ApiKey    string        `env:"API_KEY" name:"api-key" help:"API key for request auth on EdgeImpulse"`
	RemoveAll *removeAllCfg `cmd:"" help:"Removes all training data from project"`
	Upload    *uploadCfg    `cmd:"" help:"Uploads one file to edge impulse"`
}

type removeAllCfg struct {
}

type uploadCfg struct {
	Label string `default:"Unknown" short:"l" help:"File label to assign on EdgeImpulse"`
	HMAC  string `env:"HMAC" help:"HMAC string for signing a file"`
	File  string `arg:"" help:"File to upload"`
}

type Wrapper struct {
	logger log.Logger
	config Config
	ctx    context.Context
}

func New(ctx context.Context, logger log.Logger, cfg Config) (*Wrapper, error) {
	if cfg.ApiKey == "" {
		return nil, errors.New("api-key can't be empty")
	}
	res := &Wrapper{
		logger: log.With(logger, "component", ComponentName),
		config: cfg,
		ctx:    ctx,
	}

	return res, nil
}

//RemoveAll drops all training files from EdgeImpulse
func (w *Wrapper) RemoveAll() error {
	if w.config.Project == "" {
		return errors.New("project can't be empty")
	}
	level.Info(w.logger).Log("msg", "removeAll", "project", w.config.Project)

	url := fmt.Sprintf("https://studio.edgeimpulse.com/v1/api/%s/raw-data/delete-all", w.config.Project)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", w.config.ApiKey)
	ctx, cf := context.WithTimeout(w.ctx, time.Second*10)
	defer cf()
	req = req.WithContext(ctx)

	return w.invoke(req)
}

//Upload uploads file to EdgeImpulse
func (w *Wrapper) Upload(file, label, hmac string) error {
	if file == "" {
		return errors.New("no file")
	}
	if hmac == "" {
		return errors.New("no hmac")
	}
	level.Info(w.logger).Log("msg", "upload", "file", w.config.Upload.File)

	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "can't read %s", file)
	}
	signData, err := sign(hmac, fileData)
	if err != nil {
		return errors.Wrapf(err, "can't sign file data")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err = addFile(writer, "metadata.json", "application/json", strings.NewReader(signData))
	if err != nil {
		return errors.Wrapf(err, "can't add json header file")
	}
	err = addFile(writer, filepath.Base(file), "image/jpeg", bytes.NewReader(fileData))
	if err != nil {
		return errors.Wrapf(err, "can't add image")
	}
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, "https://ingestion.edgeimpulse.com/api/training/data", body)
	if err != nil {
		return errors.Wrapf(err, "can't prepare request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", w.config.ApiKey)
	req.Header.Add("x-file-name", filepath.Base(file))
	req.Header.Add("x-disallow-duplicates", "true")
	if label != "" {
		req.Header.Add("x-label", label)
	}
	ctx, cf := context.WithTimeout(w.ctx, time.Second*10)
	defer cf()
	req = req.WithContext(ctx)
	return w.invoke(req)
}

func (w *Wrapper) invoke(req *http.Request) error {
	level.Info(w.logger).Log("msg", "invoke", "url", req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "can't invoke %s", req.URL.String())
	}
	defer resp.Body.Close()
	return validateResp(resp)
}
