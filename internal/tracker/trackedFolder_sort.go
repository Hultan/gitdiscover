package tracker

type ByName struct{ TrackedFolders }

func (b ByName) Less(i, j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() {
		return true
	}
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() {
		return false
	}
	return b.TrackedFolders[i].name < b.TrackedFolders[j].name
}

type ByModifiedDate struct{ TrackedFolders }

func (b ByModifiedDate) Less(i, j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() {
		return true
	}
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() {
		return false
	}
	return b.TrackedFolders[i].ModifiedDate().After(b.TrackedFolders[j].ModifiedDate())
}

type ByChanges struct{ TrackedFolders }

func (b ByChanges) Less(i, j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() {
		return true
	}
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() {
		return false
	}
	return b.TrackedFolders[i].changes > b.TrackedFolders[j].changes
}
