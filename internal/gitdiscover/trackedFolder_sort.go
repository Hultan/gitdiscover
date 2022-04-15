package gitdiscover

// ByName sorts TrackedFolders by their name.
type ByName struct{ TrackedFolders }

// Less is a helper function that sorts by their names.
func (b ByName) Less(i, j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() {
		return true
	}
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() {
		return false
	}
	return b.TrackedFolders[i].name < b.TrackedFolders[j].name
}

// ByModifiedDate sorts TrackedFolders by their modification date.
type ByModifiedDate struct{ TrackedFolders }

// Less is a helper function that sorts by their names.
func (b ByModifiedDate) Less(i, j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() {
		return true
	}
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() {
		return false
	}
	return b.TrackedFolders[i].ModifiedDate().After(b.TrackedFolders[j].ModifiedDate())
}

// ByChanges sorts TrackedFolders by the amount of changed files in the repository.
type ByChanges struct{ TrackedFolders }

// Less is a helper function that sorts by their names.
func (b ByChanges) Less(i, j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() {
		return true
	}
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() {
		return false
	}
	return b.TrackedFolders[i].changes > b.TrackedFolders[j].changes
}
