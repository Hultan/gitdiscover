package gitdiscover

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	gitStatus "github.com/hultan/gitstatus"
	gitStatusPrompt "github.com/hultan/gitstatusprompt"
	goMod "github.com/hultan/gomod"
)

// TrackedFolders is a slice of tracked folders.
type TrackedFolders []*TrackedFolder

// TrackedFolder represents a tracked git repositry or a standard folder.
type TrackedFolder struct {
	name         string
	path         string
	isGit        bool
	modifiedDate time.Time
	imagePath    string
	gitStatus    string
	goStatus     string
	changes      int
	hasRemote    bool
}

func newFolder(folder string) *TrackedFolder {
	f := TrackedFolder{path: strings.Trim(folder, " ")}
	f.refresh()
	return &f
}

// Len makes sure that TrackedFolders implements the Interface interface
func (f TrackedFolders) Len() int { return len(f) }

// Swap makes sure that TrackedFolders implements the Interface interface
func (f TrackedFolders) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func (t *TrackedFolder) refresh() {
	t.name = path.Base(t.path)
	t.isGit = t.isGitFolder(path.Join(t.path, ".git"))
	t.modifiedDate = t.getModifiedDate(t.path)
	if t.isGit {
		t.hasRemote = t.getHasRemote(t.path)
		t.gitStatus = t.getGitStatus(t.path)
		t.goStatus = t.getGoStatus(t.path)
		t.changes = t.getNoOfChanges(t.path)
	}
}

// Name returns the name of the repository.
func (t *TrackedFolder) Name() string {
	return t.name
}

// Path returns the repository path.
func (t *TrackedFolder) Path() string {
	return t.path
}

// SetPath lets the user change the path to the repository.
func (t *TrackedFolder) SetPath(newPath string) {
	t.path = newPath
	t.refresh()
}

// ImagePath returns the path to the TrackedFolders image.
func (t *TrackedFolder) ImagePath() string {
	return t.imagePath
}

// setImagePath lets the user change the path to the TrackedFolders image.
func (t *TrackedFolder) setImagePath(newPath string) {
	t.imagePath = newPath
}

// GitStatus returns the git status
func (t *TrackedFolder) GitStatus() string {
	return t.gitStatus
}

// GoStatus returns the go status
func (t *TrackedFolder) GoStatus() string {
	return t.goStatus
}

// HasRemote returns true if the repository has a Git remote repository.
func (t *TrackedFolder) HasRemote() string {
	if !t.IsGit() {
		return "                   "
	}
	if t.hasRemote {
		return "has remote"
	}

	return "                   "
}

// IsGit returns true if the folder points to a Git repository (has a .git folder).
func (t *TrackedFolder) IsGit() bool {
	return t.isGit
}

// ModifiedDate returns the modified date of the folder.
func (t *TrackedFolder) ModifiedDate() time.Time {
	return t.modifiedDate
}

// Changes returns the number of changes to the folder.
func (t *TrackedFolder) Changes() int {
	return t.changes
}

func (t *TrackedFolder) isGitFolder(gitFolder string) bool {
	_, err := os.Stat(gitFolder)
	return !os.IsNotExist(err)
}

// Get the modified date of a file
func (t *TrackedFolder) getModifiedDate(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// Get the git status
func (t *TrackedFolder) getGitStatus(path string) string {
	p := gitStatusPrompt.GitStatusPrompt{}
	status := p.GetPrompt(path)
	return status
}

// Get the go status
func (t *TrackedFolder) getGoStatus(path string) string {
	m := goMod.GoMod{}
	info := m.GetInfo(path)
	if info == nil {
		return fmt.Sprintf("%15s", "")
	}
	result := fmt.Sprintf("Go %s", info.Version())
	if len(result) == 0 {
		return fmt.Sprintf("%15s", result)
	} else {
		return fmt.Sprintf("%10s", result)
	}
}

func (t *TrackedFolder) getNoOfChanges(path string) int {
	gs := gitStatus.GitStatus{}
	status, err := gs.GetStatus(path)
	if err != nil {
		return 0
	}
	return status.Untracked() + status.Modified() + status.Deleted() + status.Unmerged()
}

func (t *TrackedFolder) getHasRemote(repoPath string) bool {
	configPath := path.Join(repoPath, ".git", "config")
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return false
	}
	expr := `\[remote`
	r := regexp.MustCompile(expr)
	result := r.FindString(string(buf))
	if result == "" {
		return false
	}
	return true
}
