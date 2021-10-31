package gui

const (
	ApplicationVersion   = "3.7.1"
	ApplicationTitle     = "GitDiscover"
	ApplicationCopyRight = "©SoftTeam AB, 2021"
)

type SortByColumn int

const (
	SortByName SortByColumn = iota
	SortByModifiedDate
)
