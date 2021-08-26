// Encoding: UTF-8
//
// GitHub Release Client
//
// Copyright © 2021 Brian Dwyer - Broadridge Financial Solutions. All rights reserved.
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func (r *Repo) GetLatestGithubRelease() string {
	// https://docs.github.com/en/rest/reference/repos#get-the-latest-release
	client := http.Client{Timeout: 600 * time.Millisecond}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%v/releases/latest", r.Name), nil)
	if err != nil {
		log.Fatal(err)
	}

	if token, present := os.LookupEnv("GITHUB_API_KEY"); present {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		if e, ok := err.(*url.Error); ok {
			if e.Timeout() {
				log.Debugf("Timeout - %#v", e)
				return ""
			}
		}
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalln(resp.Status, req.URL.String())
	}

	release := struct {
		TagName     string    `json:"tag_name"`
		Name        string    `json:"name"`
		Draft       bool      `json:"draft"`
		Prerelease  bool      `json:"prerelease"`
		PublishedAt time.Time `json:"published_at"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Fatal(err)
	}

	return release.TagName
}
