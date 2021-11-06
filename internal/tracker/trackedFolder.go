package tracker

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// TrackedFolders is a slice of tracked folders.
type TrackedFolders []*TrackedFolder

type TrackedFolder struct {
	name         string
	path         string
	isGit        bool
	modifiedDate time.Time
	imagePath    string
	gitStatus    string
	goStatus     string
	changes      int
}

func NewFolder(folder string) *TrackedFolder {
	f := TrackedFolder{path: strings.Trim(folder, " ")}
	f.Refresh()
	return &f
}

func (f TrackedFolders) Len() int      { return len(f) }
func (f TrackedFolders) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func (f *TrackedFolder) Refresh() {
	f.name = path.Base(f.path)
	f.isGit = f.isGitFolder(path.Join(f.path, ".git"))
	f.modifiedDate = f.getModifiedDate(f.path)
	if f.isGit {
		f.gitStatus = f.getGitStatus(f.path)
		f.goStatus = f.getGoStatus(f.path)
		f.changes = f.getNoOfChanges(f.gitStatus)
	}
}

func (f *TrackedFolder) Name() string {
	return f.name
}

func (f *TrackedFolder) Path() string {
	return f.path
}

func (f *TrackedFolder) SetPath(newPath string) {
	f.path = newPath
	f.Refresh()
}

func (f *TrackedFolder) ImagePath() string {
	return f.imagePath
}

func (f *TrackedFolder) SetImagePath(newPath string) {
	f.imagePath = newPath
}

func (f *TrackedFolder) GitStatus() string {
	return f.gitStatus
}

func (f *TrackedFolder) GoStatus() string {
	return f.goStatus
}
func (f *TrackedFolder) IsGit() bool {
	return f.isGit
}

func (f *TrackedFolder) ModifiedDate() time.Time {
	return f.modifiedDate
}

func (f *TrackedFolder) Changes() int {
	return f.changes
}

func (f *TrackedFolder) isGitFolder(gitFolder string) bool {
	_, err := os.Stat(gitFolder)
	return !os.IsNotExist(err)
}

// Get the modified date of a file
func (f *TrackedFolder) getModifiedDate(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// Get the git status
func (f *TrackedFolder) getGitStatus(path string) string {
	const gitPromptCommand = "/home/per/bin/gitprompt-go"
	const gitPromptCommandFormat = "$(BRANCH)$(AHEAD)$(BEHIND)$(SEPARATOR)$(UNTRACKED)$(MODIFIED)$(DELETED)$(UNMERGED)$(STAGED)"
	cmd := exec.Command(gitPromptCommand, "-f", gitPromptCommandFormat)
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.Replace(string(out), "\n", "", -1)
}

// Get the go status
func (f *TrackedFolder) getGoStatus(path string) string {
	const gitPromptCommand = "/home/per/bin/gitprompt-go"
	const gitPromptCommandFormat = "$(GOVERSION)"
	cmd := exec.Command(gitPromptCommand, "-f", gitPromptCommandFormat)
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.Replace(string(out), "\n", "", -1)
}

func (f *TrackedFolder) getNoOfChanges(status string) int {
	fields := strings.FieldsFunc(status, func(r rune) bool {
		return !strings.ContainsRune("0123456789", r)
	})
	changes := 0
	for _, field := range fields {
		c, err := strconv.Atoi(field)
		if err != nil {
			c = 0
		}
		changes += c
	}
	return changes
}
