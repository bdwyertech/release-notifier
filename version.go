// Encoding: UTF-8
//
// GitHub Notifier - Version
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services. All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"runtime"
)

var versionFlag bool

func init() {
	flag.BoolVar(&versionFlag, "version", false, "Display version")
}

var GitCommit, ReleaseVer, ReleaseDate string

func showVersion() {
	if GitCommit == "" {
		GitCommit = "DEVELOPMENT"
	}
	if ReleaseVer == "" {
		ReleaseVer = "DEVELOPMENT"
	}
	fmt.Println("version:", ReleaseVer)
	fmt.Println("commit:", GitCommit)
	fmt.Println("date:", ReleaseDate)
	fmt.Println("runtime:", runtime.Version())
}
