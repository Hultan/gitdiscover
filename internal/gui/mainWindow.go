package gui

import (
	"bytes"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

type MainWindow struct {
	builder           *SoftBuilder
	window            *gtk.ApplicationWindow
	repositoryListBox *gtk.ListBox
	repositories      []gitdiscover.RepositoryStatus
	terminalOrNemo    *gtk.ToggleToolButton
}

// NewMainWindow : Creates a new MainWindow object
func NewMainWindow() *MainWindow {
	mainWindow := new(MainWindow)
	return mainWindow
}

// OpenMainWindow : Opens the MainWindow window
func (m *MainWindow) OpenMainWindow(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	m.builder = SoftBuilderNew("main.glade")

	// Get the main window from the glade file
	m.window = m.builder.getObject("mainWindow").(*gtk.ApplicationWindow)

	// Set up main window
	m.window.SetApplication(app)
	m.window.SetTitle(fmt.Sprintf("%s - %s", applicationTitle, applicationVersion))

	// Hook up the destroy event
	_ = m.window.Connect("destroy", m.closeMainWindow)

	// Toolbar
	m.setupToolBar()

	//// Menu
	//m.setupMenu(m.window)

	// Status bar
	lblInformation := m.builder.getObject("lblApplicationInfo").(*gtk.Label)
	lblInformation.SetText(fmt.Sprintf("%s %s - %s", applicationTitle, applicationVersion, applicationCopyRight))

	// Repository list box
	m.repositoryListBox = m.builder.getObject("repositoryListBox").(*gtk.ListBox)

	// Refresh repository list
	m.refreshList()

	// Show the main window
	m.window.ShowAll()
}

func (m *MainWindow) closeMainWindow() {
	m.window.Close()
	m.repositoryListBox.Destroy()
	m.repositoryListBox = nil
	m.repositories = nil
	m.terminalOrNemo = nil
	m.window = nil
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
	_ = button.Connect("clicked", m.refreshButtonClicked)

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

func (m *MainWindow) addButtonClicked() {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons("Select path...",
		nil,
		gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"OK",
		gtk.RESPONSE_OK,
		"Cancel",
		gtk.RESPONSE_CANCEL)
	if err != nil {
		panic(err)
	}
	dialog.SetModal(true)
	response := dialog.Run()
	if response == gtk.RESPONSE_CANCEL {
		dialog.Destroy()
		return
	}

	config := gitConfig.NewConfig()
	config.Load()
	config.Paths = append(config.Paths, dialog.GetFilename())
	config.Save()
	fmt.Println("Added path : ", dialog.GetFilename())
	dialog.Destroy()
	m.refreshList()
}

func (m *MainWindow) removeButtonClicked() {
	repo := m.getSelectedRepo()
	path := strings.Trim(repo.Path, " ")

	config := gitConfig.NewConfig()
	config.Load()
	for i, v := range config.Paths {
		if v == path {
			config.Paths = append(config.Paths[:i], config.Paths[i+1:]...)
			break
		}
	}
	config.Save()
	fmt.Println("Removed path : ", path)
	m.refreshList()
}

func (m *MainWindow) refreshButtonClicked() {
	m.refreshList()
}

func (m *MainWindow) refreshList() {
	// Clear list
	m.clearList()

	// Get repositories
	config := gitConfig.NewConfig()
	err := config.Load()
	if err != nil {
		panic(err)
	}

	git := gitdiscover.GitNew(config)
	m.repositories, err = git.GetRepositories()
	if err != nil {
		panic(err)
	}

	// Fill list
	for i := range m.repositories {
		repo := m.repositories[i]
		listItem, err := m.createListBoxItem(i, config.DateFormat, repo)
		if err != nil {
			panic(err)
		}
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
		i++
	}
}

func (m *MainWindow) createListBoxItem(index int, dateFormat string, repo gitdiscover.RepositoryStatus) (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	box.SetName(fmt.Sprintf("box_%v", index))
	if err != nil {
		return nil, err
	}

	// Icon
	assetsPath := path.Join(repo.Path, "assets")
	iconPath := path.Join(assetsPath, "application.png")
	if !fileExists(assetsPath) || !fileExists(iconPath) {
		// General icon for project that don't have one
		iconPath, err = getResourcePath("code.png")
		if err != nil {
			return nil, err
		}
	}
	pix, err := gdk.PixbufNewFromFileAtSize(iconPath,16,16)
	if err != nil {
		return nil,err
	}
	image, err := gtk.ImageNewFromPixbuf(pix)
	if err != nil {
		return nil,err
	}
	box.PackStart(image, false, false, 10)


	// Date
	label, err := gtk.LabelNew("")
	if err != nil {
		return nil, err
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
		return nil, err
	}
	text = `<span font="Sans Regular 10" foreground="#44DD44">` + repo.Status + `</span>`
	label.SetMarkup(text)
	label.SetName("lblStatus")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, false, false, 10)

	// Path
	label, err = gtk.LabelNew(repo.Path)
	if err != nil {
		return nil, err
	}
	label.SetName("lblPath")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackEnd(label, true, true, 10)

	return box, nil
}

func (m *MainWindow) openInTerminal(path string) {
	m.executeCommand("gnome-terminal", "--working-directory="+path)
}

func (m *MainWindow) openInNemo(path string) {
	m.executeCommand("nemo", path)
}

func (m *MainWindow) openInGoland(path string) {
	m.executeCommand("goland", path)
}

func (m *MainWindow) executeCommand(command, path string) {

	cmd := exec.Command(command, path)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	
	// set var to get the output
	var out bytes.Buffer

	// set the output to our variable
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(out.String())
}

func (m *MainWindow) openConfigButtonClicked() {
	m.executeCommand("xed","/home/per/.config/softteam/gitdiscover/config.json")
}

func (m *MainWindow) openInTerminalButtonClicked() {
	repo := m.getSelectedRepo()
	m.openInTerminal(repo.Path)
}

func (m *MainWindow) openInNemoButtonClicked() {
	repo := m.getSelectedRepo()
	m.openInNemo(repo.Path)
}

func (m *MainWindow) openInGolandButtonClicked() {
	repo := m.getSelectedRepo()
	m.openInGoland(repo.Path)
}

func (m *MainWindow) getSelectedRepo() gitdiscover.RepositoryStatus {
	row := m.repositoryListBox.GetSelectedRow()
	index := row.GetIndex()
	repo := m.repositories[index]
	repo.Path = strings.Trim(repo.Path, " ")

	return repo
}