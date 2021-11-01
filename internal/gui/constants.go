package gui

const (
	ApplicationVersion   = "3.7.3"
	ApplicationTitle     = "GitDiscover"
	ApplicationCopyRight = "©SoftTeam AB, 2021"
)

type SortByColumn int

const (
	SortByName SortByColumn = iota
	SortByModifiedDate
	SortByChanges
)
