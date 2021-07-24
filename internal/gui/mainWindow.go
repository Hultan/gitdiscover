package gui

import (
	"bytes"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

type MainWindow struct {
	logger *logrus.Logger
	config *gitConfig.Config

	builder           *SoftBuilder
	window            *gtk.ApplicationWindow
	repositoryListBox *gtk.ListBox
	repositories      []gitdiscover.RepositoryStatus
	terminalOrNemo    *gtk.ToggleToolButton
	aboutDialog       *gtk.AboutDialog
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
	m.builder = SoftBuilderNew("main.glade", m.logger)

	// Get the main window from the glade file
	m.window = m.builder.getObject("mainWindow").(*gtk.ApplicationWindow)

	// Set up main window
	m.window.SetApplication(app)
	m.window.SetTitle(fmt.Sprintf("%s - %s", ApplicationTitle, ApplicationVersion))
	_ = m.window.Connect("destroy", m.closeMainWindow)

	// Toolbar
	m.setupToolBar()

	// MenuBar
	m.setupMenuBar()

	// Status bar
	lblInformation := m.builder.getObject("lblApplicationInfo").(*gtk.Label)
	lblInformation.SetText(fmt.Sprintf("%s %s - %s", ApplicationTitle, ApplicationVersion, ApplicationCopyRight))

	// Repository list box
	m.repositoryListBox = m.builder.getObject("repositoryListBox").(*gtk.ListBox)

	// Refresh repository list
	m.refreshRepositoryList()

	// Show the main window
	m.window.ShowAll()
}

func (m *MainWindow) closeMainWindow() {
	m.logger = nil
	m.window.Close()
	m.repositoryListBox.Destroy()
	m.repositoryListBox = nil
	m.repositories = nil
	m.terminalOrNemo = nil
	m.window = nil
	m.builder.destroy()
	m.builder = nil
}

func (m *MainWindow) setupToolBar() {
	// Quit button
	button := m.builder.getObject("toolbarQuitButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.window.Close)

	// Add button
	button = m.builder.getObject("toolbarAddButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.addButtonClicked)

	// Remove button
	button = m.builder.getObject("toolbarRemoveButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.removeButtonClicked)

	// Refresh button
	button = m.builder.getObject("toolbarRefreshButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.refreshRepositoryList)

	// Open config button
	button = m.builder.getObject("toolbarOpenConfigButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.openConfigButtonClicked)

	// Terminal/Nemo/GoLand buttons
	button = m.builder.getObject("toolbarTerminal").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.openInTerminalButtonClicked)
	button = m.builder.getObject("toolbarNemo").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.openInNemoButtonClicked)
	button = m.builder.getObject("toolbarGoland").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.openInGolandButtonClicked)
}

func (m *MainWindow) setupMenuBar() {
	// File menu
	button := m.builder.getObject("menuFileQuit").(*gtk.MenuItem)
	_ = button.Connect("activate", m.window.Close)

	// Edit menu
	button = m.builder.getObject("menuEditExternalApplications").(*gtk.MenuItem)
	_ = button.Connect("activate", m.openAboutDialog)

	// About menu
	button = m.builder.getObject("menuHelpAbout").(*gtk.MenuItem)
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

	repo := gitConfig.Repository{Path: dialog.GetFilename(), ImagePath: "assets/application.png"}
	m.config.Repositories = append(m.config.Repositories, repo)
	m.config.Save()
	fmt.Println("Added path : ", dialog.GetFilename())
	dialog.Destroy()
	m.refreshRepositoryList()
}

func (m *MainWindow) removeButtonClicked() {
	repo := m.getSelectedRepo()
	if repo == nil {
		return
	}

	trimmedPath := strings.Trim(repo.Path, " ")

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

	git := gitdiscover.NewGit(m.config, m.logger)
	repos, err := git.GetRepositories()
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	m.repositories = repos

	// Fill list
	for i := range m.repositories {
		repo := m.repositories[i]
		listItem := m.createListBoxItem(i, m.config.DateFormat, repo)
		m.repositoryListBox.Add(listItem)
	}
	m.repositoryListBox.ShowAll()
}

func (m *MainWindow) clearList() {
	children := m.repositoryListBox.GetChildren()
	if children == nil {
		return
	}
	var i uint = 0
	for ; i < children.Length(); {
		widget, _ := children.NthData(i).(*gtk.Widget)
		m.repositoryListBox.Remove(widget)
		widget.Destroy()
		i++
	}
}

func (m *MainWindow) createListBoxItem(index int, dateFormat string, repo gitdiscover.RepositoryStatus) *gtk.Box {
	// Create main box
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	box.SetName(fmt.Sprintf("box_%v", index))

	// Icon
	iconPath := path.Join(repo.Path, repo.ImagePath)
	if !fileExists(iconPath) {
		// General icon for project that don't have one
		iconPath, err = getResourcePath("code.png")
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
	var text = ""
	if repo.Date == nil {
		text = `<span font="Sans Regular 10" foreground="#DD4444">No date!                    </span>`
	} else {
		text = `<span font="Sans Regular 10" foreground="#DD4444">` + repo.Date.Format(dateFormat) + `</span>`
	}
	label.SetMarkup(text)
	label.SetName("lblDate")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackStart(label, false, false, 10)

	// Status
	label, err = gtk.LabelNew("")
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	text = `<span font="Sans Regular 10" foreground="#44DD44">` + repo.Status + `</span>`
	label.SetMarkup(text)
	label.SetName("lblStatus")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, false, false, 10)

	// Path
	label, err = gtk.LabelNew(repo.Path)
	if err != nil {
		m.logger.Panic(err)
		panic(err)
	}
	label.SetName("lblPath")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, true, true, 10)

	return box
}

func (m *MainWindow) getSelectedRepo() *gitdiscover.RepositoryStatus {
	row := m.repositoryListBox.GetSelectedRow()
	index := row.GetIndex()
	if index == -1 {
		// TODO : MessageBox "Pleade select a repo!"
		return nil
	}
	repo := m.repositories[index]
	repo.Path = strings.Trim(repo.Path, " ")

	return &repo
}

func (m *MainWindow) openConfigButtonClicked() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		m.logger.Error("Failed to get user home dir", err)
		return
	}

	m.executeCommand("xed", path.Join(homeDir, ".config/softteam/gitdiscover/config.json"))
}

func (m *MainWindow) openInTerminalButtonClicked() {
	repo := m.getSelectedRepo()
	if repo == nil {
		return
	}
	m.openInExternalApplication("Terminal", repo)
}

func (m *MainWindow) openInNemoButtonClicked() {
	repo := m.getSelectedRepo()
	if repo == nil {
		return
	}
	m.openInExternalApplication("Nemo", repo)
}

func (m *MainWindow) openInGolandButtonClicked() {
	repo := m.getSelectedRepo()
	if repo == nil {
		return
	}
	m.openInExternalApplication("GoLand", repo)
}

func (m *MainWindow) openInExternalApplication(name string, repo *gitdiscover.RepositoryStatus) {
	opened := false

	for _, openIn := range m.config.ExternalApplications {
		if openIn.Name == name {
			argument := strings.Replace(openIn.Argument, "%PATH%", repo.Path, -1)
			argument = strings.Replace(argument, "%IMAGEPATH%", repo.ImagePath, -1)
			m.executeCommand(openIn.Command, argument)
			opened = true
			break
		}
	}

	if !opened {
		m.logger.Error("Failed to open external application")
	}
}

func (m *MainWindow) executeCommand(command, path string) {

	cmd := exec.Command(command, path)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	// https://stackoverflow.com/questions/62853835/how-to-use-syscall-sysprocattr-struct-fields-for-windows-when-os-is-set-for-linu
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	// set var to get the output
	var out bytes.Buffer

	// set the output to our variable
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		m.logger.Error(err)
		return
	}

	fmt.Println(out.String())
}

func (m *MainWindow) openAboutDialog() {
	if m.aboutDialog == nil {
		about := m.builder.getObject("aboutDialog").(*gtk.AboutDialog)

		about.SetDestroyWithParent(true)
		about.SetTransientFor(m.window)
		about.SetProgramName(ApplicationTitle)
		about.SetVersion(ApplicationVersion)
		about.SetCopyright(ApplicationCopyRight)
		about.SetComments("Discover your GIT repositories...")
		about.SetModal(true)
		about.SetPosition(gtk.WIN_POS_CENTER)

		_ = about.Connect("response", func(dialog *gtk.AboutDialog, responseId gtk.ResponseType) {
			if responseId == gtk.RESPONSE_CANCEL || responseId == gtk.RESPONSE_DELETE_EVENT {
				about.Hide()
			}
		})

		m.aboutDialog = about
	}

	m.aboutDialog.ShowAll()
}
