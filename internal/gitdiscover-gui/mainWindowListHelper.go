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
var columnColors = []string{"8DB38B", "8DB38B", "D2AB99", "8DB38B", "8DB38B", "4D934B"}
var headerColor = "00002C"

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
	m.discover.AddRepository(dialog.GetFilename(), imagePath, false)
	m.discover.Save()

	m.refreshRepositoryList()
}

func (m *MainWindow) editRepositoryButtonClicked() {
	// Get the selected repo
	repo := m.getSelectedRepo()
	if repo == nil {
		m.infoBar.showInfoWithTimeout("Please select a project to edit.", 5)
		return
	}

	// Create and show the edit repo window
	win := newEditFolderWindow(m)
	win.openWindow(repo)
}

func (m *MainWindow) removeRepositoryButtonClicked() {
	// Get the selected repo
	repo := m.getSelectedRepo()
	if repo == nil {
		m.infoBar.showInfoWithTimeout("Please select a project to remove.", 5)
		return
	}

	// Remove the selected repo
	trimmedPath := strings.Trim(repo.Path(), " ")
	m.discover.RemoveRepository(trimmedPath)

	// Save the config
	m.discover.Save()
	m.refreshRepositoryList()
}

func (m *MainWindow) refreshRepositoryList() {
	// Clear list
	m.clearList()

	// Refresh repository list
	m.discover.Refresh()

	// Sort tracked folders in the order the user have selected
	m.sortRepositories()

	// Fill list
	m.fillRepositoryList()

	// Add separators between git and non-git repositories
	m.addSeparators()

	m.repositoryListBox.ShowAll()
	m.infoBar.hideInfoBar()
}

func (m *MainWindow) addSeparators() {
	// Add the favorites separator
	sepItem := m.createListSeparator("FAVORITES")
	m.repositoryListBox.Insert(sepItem, 0)

	// Add favorites header
	hdrItem := m.createHeaderItem()
	m.repositoryListBox.Insert(hdrItem, 1)

	// Add git repository separator and header
	index := m.getFavoritesBoundaryIndex()
	if index != -1 {
		sepItem = m.createListSeparator("GIT REPOSITORIES")
		m.repositoryListBox.Insert(sepItem, index+2)
		hdrItem = m.createHeaderItem()
		m.repositoryListBox.Insert(hdrItem, index+3)
	}

	// Addd non-git repository separator and header
	index = m.getGitRepositoryBoundaryIndex()
	if index != -1 {
		sepItem = m.createListSeparator("NON-GIT FOLDERS")
		m.repositoryListBox.Insert(sepItem, index+4)
		hdrItem = m.createHeaderItem()
		m.repositoryListBox.Insert(hdrItem, index+5)
	}
}

func (m *MainWindow) fillRepositoryList() {
	// Loop through the list of repos and add them to the list
	sgDate, _ := gtk.SizeGroupNew(gtk.SIZE_GROUP_BOTH)
	for i := range m.discover.Repositories {
		repo := m.discover.GetRepositoryByIndex(i)
		listItem := m.createListItem(i, m.discover.GetDateFormat(), repo, sgDate)
		m.repositoryListBox.Add(listItem)
	}
}

func (m *MainWindow) getFavoritesBoundaryIndex() int {
	// Find where to insert the "Git repositores" separator for non-git repos
	var index = -1
	for i, repo := range m.discover.Repositories {
		if !repo.IsFavorite() {
			index = i
			break
		}
	}
	return index
}

func (m *MainWindow) getGitRepositoryBoundaryIndex() int {
	// Find where to insert the "Non-git repositories" separator for non-git repos
	var index = -1
	for i, repo := range m.discover.Repositories {
		if !repo.IsGit() {
			index = i
			break
		}
	}
	return index
}

func (m *MainWindow) sortRepositories() {
	// Sort repos by [Name|ModifiedDate|Changes] and then [IsGit]
	switch m.sortBy {
	case sortByName:
		sort.Sort(gitdiscover.ByName{Repositories: m.discover.Repositories})
	case sortByModifiedDate:
		sort.Sort(gitdiscover.ByModifiedDate{Repositories: m.discover.Repositories})
	case sortByChanges:
		sort.Sort(gitdiscover.ByChanges{Repositories: m.discover.Repositories})
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
	label.SetMarkup(`<span font="Sans Regular 14" foreground="#8C8C00">` + text + `</span>`)
	box.PackStart(label, true, true, 10)

	return box
}

func (m *MainWindow) createHeaderItem() *gtk.Box {
	// Create main box
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}

	// Icon
	label, err := gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Icon", headerColor))
	label.SetName("hdrIcon")
	label.SetTooltipText("Repository icon")
	box.PackStart(label, false, false, 5)

	// Date
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Date                       ", headerColor))
	label.SetName("hdrDate")
	label.SetTooltipText("Modified date of the git repository folder")
	box.PackStart(label, false, false, 10)

	// Favorite icon
	label, err = gtk.LabelNew("Favorite")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Favorite", headerColor))
	label.SetName("hdrFavorite")
	label.SetTooltipText("Is the repository marked as a user favorite?")
	box.PackStart(label, false, false, 0)

	// Path
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Path", headerColor))
	label.SetName("hdrPath")
	label.SetTooltipText("Repository path")
	box.PackStart(label, false, false, 10)

	// HasRemote
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Remote", headerColor))
	label.SetName("hdrHasRemote")
	label.SetTooltipText("Does the repository have a remote repository?")
	box.PackEnd(label, false, false, 0)

	// GoStatus
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Go status", headerColor))
	label.SetName("hdrGoStatus")
	label.SetTooltipText("The go version set in the go.mod file")
	box.PackEnd(label, false, false, 2)

	// GitStatus
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(m.getMarkup("Git status  ", headerColor))
	label.SetName("hdrGitStatus")
	label.SetTooltipText("The branch and the status of the git branch (modified,added, deleted, etc...)")
	box.PackEnd(label, false, false, 10)

	return box
}

func (m *MainWindow) createListItem(index int, dateFormat string, repo *gitdiscover.Repository,
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
	} else {
		label.SetMarkup(m.getMarkup(repo.HasRemote(), columnColors[5]))
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

	// Favorite icon
	if repo.IsFavorite() {
		iconPath = fw.Resource.GetResourcePath("favorite.png")
		pix, err = gdk.PixbufNewFromFileAtSize(iconPath, 16, 16)
		if err != nil {
			m.logger.Panic(err)
			panic(err)
		}
		image, err = gtk.ImageNewFromPixbuf(pix)
		image.SetTooltipText("Favorite repo")
		if err != nil {
			m.logger.Panic(err)
			panic(err)
		}
	} else {
		iconPath = fw.Resource.GetResourcePath("favorite.png")
		pix, err = gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, 16, 16)
		if err != nil {
			m.logger.Panic(err)
			panic(err)
		}
		pix.Fill(0)
		image, err = gtk.ImageNewFromPixbuf(pix)
		image.SetTooltipText("Not a favorite repo")
		if err != nil {
			m.logger.Panic(err)
			panic(err)
		}
	}
	box.PackStart(image, false, false, 10)

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
