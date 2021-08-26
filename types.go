// Encoding: UTF-8
//
// GitHub Notifier - Types
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services. All rights reserved.
//

package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-containerregistry/pkg/authn"
)

type Config struct {
	Repos []*Repo `yaml:"repos"`
}

type Repo struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`
	AuthConfig *authn.AuthConfig `yaml:"authconfig"`
	Webhooks   []*Webhook        `yaml:"webhooks"`
	Interval   string            `yaml:"interval"`
	interval   time.Duration
	version    string
}

func (r *Repo) Check() (version string, err error) {
	if r.interval == 0 {
		if r.Interval != "" {
			var err error
			if r.interval, err = time.ParseDuration(r.Interval); err != nil {
				log.Fatal(err)
			}
		} else {
			r.interval = 10 * time.Minute
		}
	}

	if r.Type == "docker" {
		return r.GetLatestDockerRelease(), nil
	} else {
		return r.GetLatestGithubRelease(), nil
	}
}
