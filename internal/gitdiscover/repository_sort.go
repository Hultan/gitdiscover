package gitdiscover

// ByName sorts Repositories by their name.
type ByName struct{ Repositories }

// Less is a helper function that sorts by their names.
func (b ByName) Less(i, j int) bool {
	if b.Repositories[i].IsGit() && !b.Repositories[j].IsGit() {
		return true
	}
	if b.Repositories[j].IsGit() && !b.Repositories[i].IsGit() {
		return false
	}
	return b.Repositories[i].name < b.Repositories[j].name
}

// ByModifiedDate sorts Repositories by their modification date.
type ByModifiedDate struct{ Repositories }

// Less is a helper function that sorts by their names.
func (b ByModifiedDate) Less(i, j int) bool {
	if b.Repositories[i].IsGit() && !b.Repositories[j].IsGit() {
		return true
	}
	if b.Repositories[j].IsGit() && !b.Repositories[i].IsGit() {
		return false
	}
	return b.Repositories[i].ModifiedDate().After(b.Repositories[j].ModifiedDate())
}

// ByChanges sorts Repositories by the amount of changed files in the repository.
type ByChanges struct{ Repositories }

// Less is a helper function that sorts by their names.
func (b ByChanges) Less(i, j int) bool {
	if b.Repositories[i].IsGit() && !b.Repositories[j].IsGit() {
		return true
	}
	if b.Repositories[j].IsGit() && !b.Repositories[i].IsGit() {
		return false
	}
	return b.Repositories[i].changes > b.Repositories[j].changes
}
