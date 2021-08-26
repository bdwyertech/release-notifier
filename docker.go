// Encoding: UTF-8
//
// Docker Registry Client
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services. All rights reserved.
//

package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1"
)

func (r *Repo) GetLatestDockerRelease() string {
	log.Infoln("Checking:", r.Name)
	opts := []crane.Option{}

	if r.AuthConfig != nil {
		opts = append(opts, crane.WithAuth(authn.FromConfig(*r.AuthConfig)))
	}

	resp, err := crane.Manifest(r.Name, opts...)
	if err != nil {
		log.Fatal(err)
	}

	var manifest v1.Manifest
	if err = json.Unmarshal(resp, &manifest); err != nil {
		log.Fatal(err)
	}

	return manifest.Config.Digest.String()
}
