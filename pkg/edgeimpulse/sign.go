package edgeimpulse

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const jsonPayload = `{
"protected": {"ver": "v1", "alg": "HS256"}, 
"signature": %s, 
"payload": {"device_name": "test-dev",
	"device_type": "phone",
	"interval_ms": 0,
	"sensors": [{"name": "image", "units": "rgba"}],
	"values": ["Ref-BINARY-image/jpeg (%d bytes) %s"]}}` 

func sign(key string, data []byte) (string, error) {
	signer := hmac.New(sha256.New, []byte(key))
    _, err := signer.Write(data)
	if (err != nil) {
		return "", errors.Wrapf(err, "can't sign file")
	}
    fileSig := signer.Sum(nil)
	emptyS := strings.Repeat("0", 64)
	json := fmt.Sprintf(jsonPayload, emptyS, len(data), hex.EncodeToString(fileSig))
	signer.Reset()
	_, err = signer.Write([]byte(json))
	if (err != nil) {
		return "", errors.Wrapf(err, "can't sign json")
	}
	jsonSig := signer.Sum(nil)
	return strings.Replace(json, emptyS, hex.EncodeToString(jsonSig), 1), nil
}