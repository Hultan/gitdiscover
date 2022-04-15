package gitdiscover_gui

import (
	"github.com/gotk3/gotk3/gtk"
)

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
	_ = button.Connect("clicked", m.addRepositoryButtonClicked)

	// Edit button
	button = m.builder.GetObject("toolbarEditButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.editRepositoryButtonClicked)

	// Remove button
	button = m.builder.GetObject("toolbarRemoveButton").(*gtk.ToolButton)
	_ = button.Connect("clicked", m.removeRepositoryButtonClicked)

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
