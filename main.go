// Encoding: UTF-8
//
// GitHub Notifier
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services. All rights reserved.
//

package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

var configPath string
var debugFlag bool

func init() {
	flag.StringVar(&configPath, "c", "", "Path to junit-gate config file")
	flag.BoolVar(&debugFlag, "debug", false, "Enable verbose log output")
}

func main() {
	// Parse Flags
	flag.Parse()

	if debugFlag {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	if versionFlag {
		showVersion()
		os.Exit(0)
	}

	if configPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		configPath = filepath.Join(pwd, ".config.yml")
	}

	yamlFile, err := RenderConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range config.Repos {
		go func(repo *Repo) {
			for {
				ver, err := repo.Check()
				if err != nil {
					log.Fatal(err)
				}
				if repo.version == "" {
					repo.version = ver
					log.Infoln(repo.version)
				}
				if repo.version == ver {
					log.Debugln(repo.Name, "Versions match or first run... Nothing to do.")
					time.Sleep(repo.interval)
					continue
				}
				repo.version = ver
				log.Infoln(repo.version)

				for _, hook := range repo.Webhooks {
					hook.Run()
				}
			}
		}(repo)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
