package edgeimpulse

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

func validateResp(resp *http.Response) error {
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("resp code: %d%s", resp.StatusCode, getBodyStr(resp.Body))
	}
	return nil
}

func getBodyStr(rd io.Reader) string {
	b, err := ioutil.ReadAll(rd)
	if err != nil {
		return "\ncan't read body"
	}
	if len(b) == 0 {
		return ""
	}
	return "\n" + string(b)
}

func createFormFileHeader(fieldName, name, fileType string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldName), escapeQuotes(name)))
	h.Set("Content-Type", fileType)
	return h
}

func addFile(writer *multipart.Writer, name, fType string, r io.Reader) error {
	part, err := writer.CreatePart(createFormFileHeader("attachments[]", name, fType))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, r)
	return err
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
