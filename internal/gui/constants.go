package gui

const (
	ApplicationVersion   = "3.7.5"
	ApplicationTitle     = "GitDiscover"
	ApplicationCopyRight = "©SoftTeam AB, 2021"
)

type sortByColumnType int

const (
	SortByName sortByColumnType = iota
	SortByModifiedDate
	SortByChanges
)

type externalApplicationModeType int

const (
	externalApplicationModeNew  externalApplicationModeType = 0
	externalApplicationModeEdit                             = 1
)

type gitCommandType uint

const (
	outputGitStatus gitCommandType = iota
	outputGitDiff
	outputGitLog
)
