package gui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/tracker"
	"github.com/hultan/softteam/framework"
)

type MainWindow struct {
	ApplicationLogPath string

	logger *logrus.Logger
	config *gitConfig.Config

	builder           *framework.GtkBuilder
	window            *gtk.ApplicationWindow
	repositoryListBox *gtk.ListBox
	tracker           *tracker.Tracker
	infoBar           *InfoBar
	toolBar           *gtk.Toolbar

	sortBy             SortByColumn
	sortByName         *gtk.RadioMenuItem
	sortByModifiedDate *gtk.RadioMenuItem
	sortByChanges      *gtk.RadioMenuItem
}

// NewMainWindow : Creates a new MainWindow object
func NewMainWindow(logger *logrus.Logger, config *gitConfig.Config) *MainWindow {
	mainWindow := new(MainWindow)
	mainWindow.logger = logger
	mainWindow.config = config
	return mainWindow
}

// OpenMainWindow : Opens the MainWindow window
func (m *MainWindow) OpenMainWindow(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("mainWindow.ui")
	if err != nil {
		panic(err)
	}
	m.builder = builder

	// Get the main window from the glade file
	m.window = m.builder.GetObject("mainWindow").(*gtk.ApplicationWindow)

	// Set up main window
	m.window.SetApplication(app)
	m.window.SetTitle(fmt.Sprintf("%s - %s", ApplicationTitle, ApplicationVersion))
	_ = m.window.Connect("destroy", m.closeMainWindow)

	// Toolbar
	m.toolBar = m.builder.GetObject("toolbar").(*gtk.Toolbar)
	m.setupToolBar()

	// MenuBar
	m.setupMenuBar()

	// Status bar
	lblInformation := m.builder.GetObject("lblApplicationInfo").(*gtk.Label)
	lblInformation.SetText(fmt.Sprintf("%s %s - %s", ApplicationTitle, ApplicationVersion, ApplicationCopyRight))

	// Info bar
	infoBar := m.builder.GetObject("infoBar").(*gtk.InfoBar)
	labelInfoBar := m.builder.GetObject("labelInfoBar").(*gtk.Label)
	m.infoBar = NewInfoBar(infoBar, labelInfoBar)

	// Repository list box
	m.repositoryListBox = m.builder.GetObject("repositoryListBox").(*gtk.ListBox)

	// Refresh repository list
	m.refreshRepositoryList()

	// Popup menu
	popup := NewPopupMenu(m)
	popup.Setup()

	// Show the main window
	m.window.ShowAll()
	m.infoBar.hideInfoBar()
}

func (m *MainWindow) closeMainWindow() {
	m.logger = nil
	m.window.Close()
	m.repositoryListBox.Destroy()
	m.repositoryListBox = nil
	m.tracker = nil
	m.window.Destroy()
	m.window = nil
	m.builder = nil
}

func (m *MainWindow) setupToolBar() {
	// Quit button
	button := m.builder.GetObject("toolbarQuitButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.window.Close)

	// Add button
	button = m.builder.GetObject("toolbarAddButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.addButtonClicked)

	// Edit button
	button = m.builder.GetObject("toolbarEditButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.editButtonClicked)

	// Remove button
	button = m.builder.GetObject("toolbarRemoveButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.removeButtonClicked)

	// Refresh button
	button = m.builder.GetObject("toolbarRefreshButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.refreshRepositoryList)

	m.refreshExternalApplications(m.toolBar)
}

func (m *MainWindow) setupMenuBar() {
	// File menu
	button := m.builder.GetObject("menuFileQuit").(*gtk.MenuItem)
	_ = button.Connect("activate", m.window.Close)

	// View menu
	m.sortByName = m.builder.GetObject("mnuSortByName").(*gtk.RadioMenuItem)
	m.sortByName.SetActive(true)
	_ = m.sortByName.Connect("activate", m.toggleSortBy)
	m.sortByModifiedDate = m.builder.GetObject("mnuSortByModifiedDate").(*gtk.RadioMenuItem)
	m.sortByModifiedDate.JoinGroup(m.sortByName)
	_ = m.sortByModifiedDate.Connect("activate", m.toggleSortBy)
	m.sortByChanges = m.builder.GetObject("mnuSortByChanges").(*gtk.RadioMenuItem)
	m.sortByChanges.JoinGroup(m.sortByName)
	_ = m.sortByChanges.Connect("activate", m.toggleSortBy)

	// Edit menu
	button = m.builder.GetObject("menuEditExternalApplications").(*gtk.MenuItem)
	_ = button.Connect("activate", m.openExternalToolsDialog)
	button = m.builder.GetObject("menuEditConfig").(*gtk.MenuItem)
	_ = button.Connect("activate", m.openConfig)
	button = m.builder.GetObject("menuEditLog").(*gtk.MenuItem)
	_ = button.Connect("activate", m.openLog)

	// About menu
	button = m.builder.GetObject("menuHelpAbout").(*gtk.MenuItem)
	_ = button.Connect("activate", m.openAboutDialog)
}

func (m *MainWindow) addButtonClicked() {
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
	dialog.SetModal(true)
	response := dialog.Run()
	if response == gtk.RESPONSE_CANCEL {
		dialog.Destroy()
		return
	}

	imagePath := filepath.Join(dialog.GetFilename(), "assets/application.png")
	repo := gitConfig.Repository{Path: dialog.GetFilename(), ImagePath: imagePath}
	m.config.Repositories = append(m.config.Repositories, repo)
	m.config.Save()
	fmt.Println("Added path : ", dialog.GetFilename())
	dialog.Destroy()
	m.refreshRepositoryList()
}

func (m *MainWindow) editButtonClicked() {
	folder := m.getSelectedRepo()
	if folder == nil {
		// TODO : Handle this error, must select folder
		return
	}
	win := NewEditFolderWindow(m)
	win.openWindow(folder)
}

func (m *MainWindow) removeButtonClicked() {
	repo := m.getSelectedRepo()
	if repo == nil {
		return
	}

	trimmedPath := strings.Trim(repo.Path(), " ")

	for i, repository := range m.config.Repositories {
		if repository.Path == trimmedPath {
			m.config.Repositories = append(m.config.Repositories[:i], m.config.Repositories[i+1:]...)
			break
		}
	}
	m.config.Save()
	fmt.Println("Removed path : ", trimmedPath)
	m.refreshRepositoryList()
}

func (m *MainWindow) refreshRepositoryList() {
	// Clear list
	m.clearList()

	m.tracker = tracker.NewTracker(m.config)

	switch m.sortBy {
	case SortByName:
		sort.Sort(tracker.ByName{TrackedFolders: m.tracker.Folders})
	case SortByModifiedDate:
		sort.Sort(tracker.ByModifiedDate{TrackedFolders: m.tracker.Folders})
	case SortByChanges:
		sort.Sort(tracker.ByChanges{TrackedFolders: m.tracker.Folders})
	}

	sgDate, _ := gtk.SizeGroupNew(gtk.SIZE_GROUP_BOTH)

	// Fill list
	for i := range m.tracker.Folders {
		repo := m.tracker.Folders[i]
		listItem := m.createListItem(i, m.config.DateFormat, repo, sgDate)
		m.repositoryListBox.Add(listItem)
	}
	sepItem := m.createListSeparator("GIT REPOSITORIES")
	m.repositoryListBox.Insert(sepItem, 0)

	// Find where to insert the "Misc folders" separator
	var index = -1
	for i, repo := range m.tracker.Folders {
		if !repo.IsGit() {
			index = i
			break
		}
	}

	if index != -1 {
		sepItem = m.createListSeparator("NON-GIT FOLDERS")
		m.repositoryListBox.Insert(sepItem, index+1)
	}

	m.repositoryListBox.ShowAll()
	m.infoBar.hideInfoBar()
}

func (m *MainWindow) clearList() {
	children := m.repositoryListBox.GetChildren()
	if children == nil {
		return
	}
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

func (m *MainWindow) createListItem(index int, dateFormat string, repo *tracker.TrackedFolder,
	sgDate *gtk.SizeGroup) *gtk.Box {

	// Create main box
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	box.SetName(fmt.Sprintf("box_%v", index))

	// Icon
	iconPath := repo.ImagePath()
	if !fileExists(iconPath) {
		// General icon for project that don't have one
		if repo.IsGit() {
			iconPath, err = getResourcePath("gitFolder.png")
		} else {
			iconPath, err = getResourcePath("folder.png")
		}
		if err != nil {
			m.logger.Panic(err)
			panic(err)
		}
	}
	pix, err := gdk.PixbufNewFromFileAtSize(iconPath, 16, 16)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	image, err := gtk.ImageNewFromPixbuf(pix)
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
	label.SetMarkup(`<span font="Sans Regular 10" foreground="#44DD44">` + repo.ModifiedDate().Format(dateFormat) + `</span>`)
	label.SetName("lblDate")
	sgDate.AddWidget(label)
	label.SetXAlign(0.0)
	box.PackStart(label, false, false, 10)

	// GitStatus
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetMarkup(`<span font="Sans Regular 10" foreground="#44DD44">` + repo.GitStatus() + `</span>`)
	label.SetName("lblStatus")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, false, false, 10)

	// Path
	label, err = gtk.LabelNew(repo.Path())
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetName("lblPath")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, true, true, 10)

	return box
}

func (m *MainWindow) getSelectedRepo() *tracker.TrackedFolder {
	row := m.repositoryListBox.GetSelectedRow()
	if row == nil {
		return nil
	}
	boxObj, err := row.GetChild()
	if err != nil {
		m.infoBar.ShowError(err.Error())
		return nil
	}
	box, ok := boxObj.(*gtk.Box)
	if !ok {
		m.infoBar.ShowError("Failed to convert to *gtk.Widget")
		return nil
	}
	name, err := box.GetName()
	if err != nil {
		m.infoBar.ShowError(err.Error())
		return nil
	}
	if name == "sep" {
		return nil
	}
	indexString := name[4:]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		m.infoBar.ShowError(err.Error())
		return nil
	}
	repo := m.tracker.Folders[index]

	return repo
}

func (m *MainWindow) openConfig() {
	go func() {
		m.executeCommand("xed", m.config.GetConfigPath())
	}()
}

func (m *MainWindow) openLog() {
	go func() {
		m.executeCommand("xed", m.ApplicationLogPath)
	}()
}

func (m *MainWindow) openInExternalApplication(name string, repo *tracker.TrackedFolder) {
	// Find application
	app := m.config.GetExternalApplicationByName(name)
	if app == nil {
		// Failed to find application, show info bar.
		// This should not happen, but if there is an issue
		// with when the external application buttons are
		// refreshed, then it can happen.
		text := fmt.Sprintf("Failed to find an application with name : %s", name)
		m.infoBar.ShowError(text)
		m.logger.Error(text)
		return
	}
	m.infoBar.hideInfoBar()

	// Open external application
	var argument = ""

	if repo != nil {
		argument = strings.Replace(app.Argument, "%PATH%", repo.Path(), -1)
		argument = strings.Replace(argument, "%IMAGEPATH%", repo.ImagePath(), -1)
	}

	m.logger.Info("Trying to open external application: ", app.Command, " ", argument)
	go func() {
		m.executeCommand(app.Command, argument)
	}()
}

func (m *MainWindow) executeCommand(command, arguments string) string {
	cmd := exec.Command(command, arguments)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	// https://stackoverflow.com/questions/62853835/how-to-use-syscall-sysprocattr-struct-fields-for-windows-when-os-is-set-for-linu
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	//	Setpgid: true,
	//	Pgid:    0,
	// }

	// set the output to our variable
	out, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.Error("Failed to open external application: ", command, " ", arguments)
		m.logger.Error(err)
		m.infoBar.ShowError(err.Error())
		return ""
	}

	return string(out)
}

func (m *MainWindow) openAboutDialog() {
	about := NewAboutDialog(m.logger, m.window)
	about.openAboutDialog()
}

func (m *MainWindow) openExternalToolsDialog() {
	window := NewExternalApplicationsWindow(m.logger, m.config)
	window.openWindow(func() {
		m.refreshExternalApplications(m.toolBar)
	})
}

func (m *MainWindow) refreshExternalApplications(toolbar *gtk.Toolbar) {
	children := toolbar.GetChildren()
	if children.Length() > 0 {
		var i uint
		for i = 0; i < children.Length(); i++ {
			child := children.NthData(i)
			toolButton, ok := child.(*gtk.Widget)
			if ok {
				name, err := toolButton.GetName()
				if err != nil {
					m.logger.Error(err)
					return
				}
				if name[:3] == "ea_" {
					toolbar.Remove(toolButton)
				}
			}
		}
	}

	for i := 0; i < len(m.config.ExternalApplications); i++ {
		app := m.config.ExternalApplications[i]
		toolButton, err := gtk.ToolButtonNew(nil, app.Name)
		if err != nil {
			m.logger.Error(err)
			panic(err)
		}
		toolButton.SetName("ea_" + app.Name)
		toolButton.Connect("clicked", func(button *gtk.ToolButton) {
			name, err := button.GetName()
			if err != nil {
				m.logger.Error(err)
				return
			}
			appName := name[3:]
			repo := m.getSelectedRepo()
			if repo == nil {
				m.logger.Error("repo not found when clicking application '", appName, "'")
			}
			m.openInExternalApplication(appName, repo)
		})
		toolbar.Add(toolButton)
	}
	toolbar.ShowAll()
}

func (m *MainWindow) toggleSortBy(radio *gtk.RadioMenuItem) {
	if !radio.GetActive() {
		return
	}

	name := radio.GetLabel()
	switch name {
	case "Name":
		m.sortBy = SortByName
		// m.sortByModifiedDate.SetActive(false)
		// m.sortByChanges.SetActive(false)
	case "Modified date":
		m.sortBy = SortByModifiedDate
		// m.sortByName.SetActive(false)
		// m.sortByChanges.SetActive(false)
	case "Changes":
		m.sortBy = SortByChanges
		// m.sortByName.SetActive(false)
		// m.sortByModifiedDate.SetActive(false)
	}
	m.refreshRepositoryList()
}
