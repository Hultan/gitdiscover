package gitdiscover_gui

const (
	applicationVersion   = "3.9.2"
	applicationTitle     = "GitDiscover"
	applicationCopyRight = "©SoftTeam AB, 2021"
)

type sortByColumnType int

const (
	sortByName sortByColumnType = iota
	sortByModifiedDate
	sortByChanges
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
