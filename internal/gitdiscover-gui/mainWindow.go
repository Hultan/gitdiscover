package gitdiscover_gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/softteam/framework"
)

// MainWindow is the main window
type MainWindow struct {
	ApplicationLogPath string

	logger *logrus.Logger
	config *gitdiscover.Config

	builder           *framework.GtkBuilder
	window            *gtk.ApplicationWindow
	repositoryListBox *gtk.ListBox
	tracker           *gitdiscover.Tracker
	infoBar           *infoBar
	toolBar           *gtk.Toolbar

	sortBy             sortByColumnType
	sortByName         *gtk.RadioMenuItem
	sortByModifiedDate *gtk.RadioMenuItem
	sortByChanges      *gtk.RadioMenuItem
}

// NewMainWindow creates a new MainWindow object
func NewMainWindow(logger *logrus.Logger, config *gitdiscover.Config) *MainWindow {
	mainWindow := new(MainWindow)
	mainWindow.logger = logger
	mainWindow.config = config
	return mainWindow
}

// OpenMainWindow opens the MainWindow window
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
	m.window.SetTitle(fmt.Sprintf("%s - %s", applicationTitle, applicationVersion))
	_ = m.window.Connect("destroy", m.closeMainWindow)

	// Toolbar
	m.toolBar = m.builder.GetObject("toolbar").(*gtk.Toolbar)
	m.setupToolBar()

	// MenuBar
	m.setupMenuBar()

	// Status bar
	lblInformation := m.builder.GetObject("lblApplicationInfo").(*gtk.Label)
	lblInformation.SetText(fmt.Sprintf("%s %s - %s", applicationTitle, applicationVersion, applicationCopyRight))

	// Info bar
	infoBar := m.builder.GetObject("infoBar").(*gtk.InfoBar)
	labelInfoBar := m.builder.GetObject("labelInfoBar").(*gtk.Label)
	m.infoBar = newInfoBar(infoBar, labelInfoBar)

	// Repository list box
	m.repositoryListBox = m.builder.GetObject("repositoryListBox").(*gtk.ListBox)

	// Refresh repository list
	m.refreshRepositoryList()

	// Popup menu
	popup := newPopupMenu(m)
	popup.setupPopupMenu()

	// Show the main window
	m.window.ShowAll()
	m.infoBar.hideInfoBar()
}

func (m *MainWindow) openAboutDialog() {
	about := newAboutDialog(m.logger, m.window)
	about.openAboutDialog()
}
