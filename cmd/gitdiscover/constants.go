package main

import "github.com/gotk3/gotk3/glib"

const (
	exitNormal = 0
	exitConfigError = 1
	exitArgumentError = 2
)

const (
	ApplicationId    = "se.softteam.gitdiscover"
	ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)
