package tracker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (f *TrackedFolder) refresh() {
	f.name = path.Base(f.path)
	f.isGit = f.isGitFolder(path.Join(f.path, ".git"))
	f.modifiedDate = f.getModifiedDate(f.path)
	if f.isGit {
		f.hasRemote = f.getHasRemote(f.path)
		f.gitStatus = f.getGitStatus(f.path)
		f.goStatus = f.getGoStatus(f.path)
		f.changes = f.getNoOfChanges(f.gitStatus)
	}
}

// Name returns the name of the repository.
func (f *TrackedFolder) Name() string {
	return f.name
}

// Path returns the repository path.
func (f *TrackedFolder) Path() string {
	return f.path
}

// SetPath lets the user change the path to the repository.
func (f *TrackedFolder) SetPath(newPath string) {
	f.path = newPath
	f.refresh()
}

// ImagePath returns the path to the TrackedFolders image.
func (f *TrackedFolder) ImagePath() string {
	return f.imagePath
}

// setImagePath lets the user change the path to the TrackedFolders image.
func (f *TrackedFolder) setImagePath(newPath string) {
	f.imagePath = newPath
}

// GitStatus returns the git status
func (f *TrackedFolder) GitStatus() string {
	return f.gitStatus
}

// GoStatus returns the go status
func (f *TrackedFolder) GoStatus() string {
	return f.goStatus
}

// HasRemote returns true if the repository has a Git remote repository.
func (f *TrackedFolder) HasRemote() string {
	if !f.IsGit() {
		return "                   "
	}
	if f.hasRemote {
		return "has remote"
	}

	return "                   "
}

// IsGit returns true if the folder points to a Git repository (has a .git folder).
func (f *TrackedFolder) IsGit() bool {
	return f.isGit
}

// ModifiedDate returns the modified date of the folder.
func (f *TrackedFolder) ModifiedDate() time.Time {
	return f.modifiedDate
}

// Changes returns the number of changes to the folder.
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
		return fmt.Sprintf("%15s", "")
	}
	result := strings.Replace(string(out), "\n", "", -1)

	if len(result) == 0 {
		return fmt.Sprintf("%15s", result)
	} else {
		return fmt.Sprintf("%10s", result)
	}
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

func (f *TrackedFolder) getHasRemote(repoPath string) bool {
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
