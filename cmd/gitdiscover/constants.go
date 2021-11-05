package main

import "github.com/gotk3/gotk3/glib"

const (
	exitNormal      = 0
	exitConfigError = 1
	// exitArgumentError = 2 // Not used anymore
	exitUnknown = 3
)

const (
	applicationId      = "se.softteam.gitdiscover"
	applicationFlags   = glib.APPLICATION_FLAGS_NONE
	applicationLogPath = "/tmp/gitdiscover.log"
)
