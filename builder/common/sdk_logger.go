package common

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Logger struct{}

func newLogger() *Logger {
	return &Logger{}
}

func (l *Logger) RequestHttp(_ context.Context, req *http.Request) {
	body := ""

	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err == nil {
			bodyBytes, _ := io.ReadAll(bodyReader)
			if len(bodyBytes) > 0 {
				body = prettyJSON(bodyBytes)
			}
		}
	}

	log.Print(l.format("[DEBUG] SDK HTTP request", req.Method, req.URL.String(), body, "", ""))
}

func (l *Logger) ResponseHttp(_ context.Context, resp *http.Response, d time.Duration) {
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	body := ""

	if len(bodyBytes) > 0 {
		body = prettyJSON(bodyBytes)
	}

	log.Print(l.format("[DEBUG] SDK HTTP response", resp.Request.Method, resp.Request.URL.String(), body, resp.Status, d.String()))
}

func (l *Logger) format(title, method, url, body, status, duration string) string {
	var msg strings.Builder
	msg.WriteString(title + ":\n")
	msg.WriteString("  | " + method + " " + url + "\n")
	if status != "" {
		msg.WriteString("  | Status: " + status + " (" + duration + ")\n")
	}

	if body != "" {
		msg.WriteString("  |\n")
		for line := range strings.SplitSeq(body, "\n") {
			msg.WriteString("  | " + line + "\n")
		}
	} else {
		msg.WriteString("\n")
	}
	msg.WriteString("\n")

	return msg.String()
}

func (l *Logger) Request(_ context.Context, req any) {
	// PACKER_LOGS cannot filter debugging level, so we only log the http req/resp and error
}

func (l *Logger) Response(_ context.Context, resp any) {
	// PACKER_LOGS cannot filter debugging level, so we only log the http req/resp and error
}

func (l *Logger) Error(_ context.Context, err error) {
	log.Printf("[ERROR] SDK error: %s", err)
}

func prettyJSON(body []byte) string {
	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}

	formatted, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return string(body)
	}

	return string(formatted)
}
