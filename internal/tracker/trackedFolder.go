package tracker

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type ByName struct { TrackedFolders }

func (b ByName) Less(i,j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() { return true }
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() { return false }
	return b.TrackedFolders[i].name < b.TrackedFolders[j].name
}

type ByModifiedDate struct { TrackedFolders }

func (b ByModifiedDate) Less(i,j int) bool {
	if b.TrackedFolders[i].IsGit() && !b.TrackedFolders[j].IsGit() { return true }
	if b.TrackedFolders[j].IsGit() && !b.TrackedFolders[i].IsGit() { return false }
	return b.TrackedFolders[i].ModifiedDate().After(b.TrackedFolders[j].ModifiedDate())
}

type TrackedFolders []*TrackedFolder

type TrackedFolder struct {
	name         string
	path         string
	isGit        bool
	modifiedDate time.Time
	imagePath    string
	gitStatus    string
}

func NewFolder(folder string) *TrackedFolder {
	f := TrackedFolder{path: strings.Trim(folder, " ")}
	f.Refresh()
	return &f
}

func (f TrackedFolders) Len() int { return len(f) }
func (f TrackedFolders) Swap(i, j int) { f[i], f[j] = f[j],f[i] }

func (f *TrackedFolder) Refresh() {
	f.name = path.Base(f.path)
	f.isGit = f.isGitFolder(path.Join(f.path, ".git"))
	f.modifiedDate = f.getModifiedDate(f.path)
	if f.isGit {
		f.gitStatus = f.getGitStatus(f.path)
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

func (f *TrackedFolder) IsGit() bool {
	return f.isGit
}

func (f *TrackedFolder) ModifiedDate() time.Time {
	return f.modifiedDate
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
	cmd := exec.Command("/home/per/bin/gitprompt-go")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.Replace(string(out), "\n", "", -1)
}
