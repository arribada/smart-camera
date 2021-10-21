package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	filename := "data/output.png"
	cmd := exec.Command("pylepton_capture", filename)
	stdout, err := cmd.Output()

	ExitOnError(logger, err, "capture run")

	level.Info(logger).Log("msg", "capture command run", "output", string(stdout))

	// Initialize minio client object.
	clt, err := minio.New("storage.googleapis.com", &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("STORAGE_KEY_ID"), os.Getenv("STORAGE_KEY_SECRET"), ""),
		Secure: true,
	})
	ExitOnError(logger, err, "initalize the storage upload client")

	objectName := filename
	filePath := filename
	contentType := "image/png"

	n, err := clt.FPutObject(context.Background(), "arribada-smart", objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	ExitOnError(logger, err, "uploading file")

	level.Info(logger).Log("msg", "upload complete", "file", objectName, "metadata", fmt.Sprintf("%+v", n))
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
