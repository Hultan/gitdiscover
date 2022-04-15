package gitdiscover_gui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/softteam/framework"
)

// Column :                  Path      Date      GitStatus GoStatus  Yes       No
var columnColors = []string{"FFFFFF", "8DB38B", "D2AB99", "ADD3AB", "C290AA", "A26690"}

func (m *MainWindow) addRepositoryButtonClicked() {
	// Create and show the folder chooser dialog
	dialog, err := gtk.FileChooserDialogNewWith2Buttons("Select path...",
		nil,
		gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"OK",
		gtk.RESPONSE_OK,
		"Cancel",
		gtk.RESPONSE_CANCEL)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	defer dialog.Destroy()

	dialog.SetModal(true)
	response := dialog.Run()
	if response == gtk.RESPONSE_CANCEL {
		return
	}

	// Add the new repository, and save the config
	imagePath := filepath.Join(dialog.GetFilename(), "assets/application.png")
	m.config.AddRepository(dialog.GetFilename(), imagePath)
	m.config.Save()

	m.refreshRepositoryList()
}

func (m *MainWindow) editRepositoryButtonClicked() {
	// Get the selected repo
	repo := m.getSelectedRepo()
	if repo == nil {
		m.infoBar.ShowInfoWithTimeout("Please select a project to edit.", 5)
		return
	}

	// Create and show the edit repo window
	win := NewEditFolderWindow(m)
	win.openWindow(repo)
}

func (m *MainWindow) removeRepositoryButtonClicked() {
	// Get the selected repo
	repo := m.getSelectedRepo()
	if repo == nil {
		m.infoBar.ShowInfoWithTimeout("Please select a project to remove.", 5)
		return
	}

	// Remove the selected repo
	trimmedPath := strings.Trim(repo.Path(), " ")
	for i, repository := range m.config.Repositories {
		if repository.Path == trimmedPath {
			m.config.Repositories = append(m.config.Repositories[:i], m.config.Repositories[i+1:]...)
			break
		}
	}

	// Save the config
	m.config.Save()
	m.refreshRepositoryList()
}

func (m *MainWindow) refreshRepositoryList() {
	// Clear list
	m.clearList()

	if m.tracker == nil {
		m.tracker = gitdiscover.NewTracker(m.config)
	} else {
		m.tracker.Refresh()
	}

	// Sort tracked folders in the order the user have selected
	m.sortTrackedFolders()

	// Fill list
	m.fillRepositoryList()

	// Add separators between git and non-git repositories
	m.addSeparators()

	m.repositoryListBox.ShowAll()
	m.infoBar.HideInfoBar()
}

func (m *MainWindow) addSeparators() {
	// Add git repository separator
	sepItem := m.createListSeparator("GIT REPOSITORIES")
	m.repositoryListBox.Insert(sepItem, 0)

	// Addd non-git repository separator
	index := m.getGitRepositoryBoundaryIndex()
	if index != -1 {
		sepItem = m.createListSeparator("NON-GIT FOLDERS")
		m.repositoryListBox.Insert(sepItem, index+1)
	}
}

func (m *MainWindow) fillRepositoryList() {
	// Loop through the list of repos and add them to the list
	sgDate, _ := gtk.SizeGroupNew(gtk.SIZE_GROUP_BOTH)
	for i := range m.tracker.Folders {
		repo := m.tracker.Folders[i]
		listItem := m.createListItem(i, m.config.DateFormat, repo, sgDate)
		m.repositoryListBox.Add(listItem)
	}
}

func (m *MainWindow) getGitRepositoryBoundaryIndex() int {
	// Find where to insert the "Misc folders" separator for non-git repos
	var index = -1
	for i, repo := range m.tracker.Folders {
		if !repo.IsGit() {
			index = i
			break
		}
	}
	return index
}

func (m *MainWindow) sortTrackedFolders() {
	// Sort repos by [Name|ModifiedDate|Changes] and then [IsGit]
	switch m.sortBy {
	case sortByName:
		sort.Sort(gitdiscover.ByName{TrackedFolders: m.tracker.Folders})
	case sortByModifiedDate:
		sort.Sort(gitdiscover.ByModifiedDate{TrackedFolders: m.tracker.Folders})
	case sortByChanges:
		sort.Sort(gitdiscover.ByChanges{TrackedFolders: m.tracker.Folders})
	}
}

func (m *MainWindow) clearList() {
	// Get the list of repos
	children := m.repositoryListBox.GetChildren()
	if children == nil {
		return
	}

	// Remove all elements from the list
	var i uint = 0
	for i < children.Length() {
		widget, _ := children.NthData(i).(*gtk.Widget)
		m.repositoryListBox.Remove(widget)
		widget.Destroy()
		i++
	}
}

func (m *MainWindow) createListSeparator(text string) *gtk.Box {
	// Create main box
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	box.SetName("sep")

	// Date
	label, err := gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(`<span font="Sans Regular 14" foreground="#DDDD00">` + text + `</span>`)
	box.PackStart(label, true, true, 10)

	return box
}

func (m *MainWindow) createListItem(index int, dateFormat string, repo *gitdiscover.TrackedFolder,
	sgDate *gtk.SizeGroup) *gtk.Box {

	// Create main box
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	box.SetName(fmt.Sprintf("box_%v", index))

	// Icon
	fw := framework.NewFramework()
	iconPath := repo.ImagePath()
	if !fw.IO.FileExists(iconPath) {
		// General icon for project that don't have one
		if repo.IsGit() {
			iconPath = fw.Resource.GetResourcePath("gitFolder.png")
		} else {
			iconPath = fw.Resource.GetResourcePath("folder.png")
		}
	}
	pix, err := gdk.PixbufNewFromFileAtSize(iconPath, 16, 16)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	image, err := gtk.ImageNewFromPixbuf(pix)
	image.SetTooltipText("Application.png from assets folder in repository")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	box.PackStart(image, false, false, 10)

	// Date
	label, err := gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup(repo.ModifiedDate().Format(dateFormat), columnColors[1]))
	label.SetName("lblDate")
	label.SetTooltipText("Modifed date of repository folder")
	sgDate.AddWidget(label)
	label.SetXAlign(0.0)
	box.PackStart(label, false, false, 10)

	// HasRemote
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	if strings.Trim(repo.HasRemote(), " ") == "yes" {
		label.SetMarkup(m.getMarkup(repo.HasRemote(), columnColors[4]))
		// label.SetMarkup(`<span font="Sans Regular 10" foreground="#44DD44">` + repo.HasRemote() + `</span>`)
	} else {
		label.SetMarkup(m.getMarkup(repo.HasRemote(), columnColors[5]))
		// label.SetMarkup(`<span font="Sans Regular 10" foreground="#DD4444">` + repo.HasRemote() + `</span>`)
	}

	label.SetName("lblHasRemote")
	label.SetTooltipText("Has a Git remote repository")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, false, false, 10)

	// GoStatus
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	// label.SetMarkup(`<span font="Sans Regular 10" foreground="#6666DD">` + repo.GoStatus() + `</span>`)
	label.SetMarkup(m.getMarkup(repo.GoStatus(), columnColors[3]))
	label.SetName("lblGoStatus")
	label.SetTooltipText("Go version from the go.mod file")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, false, false, 10)

	// GitStatus
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	// label.SetMarkup(`<span font="Sans Regular 10" foreground="#22BB88">` + repo.GitStatus() + `</span>`)
	label.SetMarkup(m.getMarkup(repo.GitStatus(), columnColors[2]))
	label.SetName("lblStatus")
	label.SetTooltipText("Result from the Git status command")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, false, false, 10)

	// Path
	label, err = gtk.LabelNew(repo.Path())
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup(repo.Path(), columnColors[0]))
	label.SetName("lblPath")
	label.SetTooltipText("Repository path")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, true, true, 10)

	return box
}

func (m *MainWindow) toggleSortBy(radio *gtk.RadioMenuItem) {
	// Only sort by the selected radio button
	if !radio.GetActive() {
		return
	}

	// Sort by the selected radio button
	name := radio.GetLabel()
	switch name {
	case "Name":
		m.sortBy = sortByName
	case "Modified date":
		m.sortBy = sortByModifiedDate
	case "Changes":
		m.sortBy = sortByChanges
	}

	m.refreshRepositoryList()
}
