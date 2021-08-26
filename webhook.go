// Encoding: UTF-8
//
// GitHub Notifier - Webhooks
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services. All rights reserved.
//

package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Webhook struct {
	Url      string            `yaml:"url"`
	Method   string            `yaml:"method,omitempty"`
	Body     string            `yaml:"body"`
	Headers  map[string]string `yaml:"headers"`
	Timeout  string            `yaml:"timeout"`
	Insecure bool              `yaml:"insecure"`
}

func (hook *Webhook) Run() {
	if hook.Method == "" {
		hook.Method = "POST"
	}
	timeout := 5 * time.Second
	if hook.Timeout != "" {
		var err error
		if timeout, err = time.ParseDuration(hook.Timeout); err != nil {
			log.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, hook.Method, hook.Url, strings.NewReader(strings.TrimSpace(hook.Body)))
	if err != nil {
		log.Fatal(err)
	}

	r.Header.Set("User-Agent", "github-notifier/"+ReleaseVer)
	for k, v := range hook.Headers {
		r.Header.Set(k, v)
	}

	// Copy of http.DefaultTransport with Flippable TLS Verification
	// https://golang.org/pkg/net/http/#Client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: hook.Insecure},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	log.Debugln("Executing Request:", r.URL.String())
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	response := make(map[string]interface{})

	contentType := resp.Header.Get("Content-Type")
	switch strings.Split(contentType, ";")[0] {
	case "application/json":
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Fatal(err)
		}
	default:
		if body, err := io.ReadAll(resp.Body); err != nil {
			response["BodyDecodeFailure"] = err
		} else {
			if bodyString := string(body); bodyString != "" {
				response["Body"] = bodyString
			} else {
				return
			}
		}
	}

	out, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof(hook.Url, string(out))
}
