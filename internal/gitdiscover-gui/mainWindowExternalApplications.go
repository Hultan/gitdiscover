package gitdiscover_gui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
)

func (m *MainWindow) openInExternalApplication(name string, repo *gitdiscover.TrackedFolder) {
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

func (m *MainWindow) openExternalToolsDialog() {
	window := newExternalApplicationsWindow(m.logger, m.config)
	window.openWindow(func() {
		m.refreshExternalApplications(m.toolBar)
	})
}

func (m *MainWindow) refreshExternalApplications(toolbar *gtk.Toolbar) {
	// Remove the old external applications from the toolbar
	m.removeToolbarApplications(toolbar)

	// Add the new external applications to the toolbar
	m.addToolbarApplications(toolbar)
}

func (m *MainWindow) addToolbarApplications(toolbar *gtk.Toolbar) {
	for i := 0; i < len(m.config.ExternalApplications); i++ {
		extApp := m.config.ExternalApplications[i]

		// Create a new toolbar button, and panic on error.
		// If we can't create a new button, we have bigger problems.
		toolButton, err := gtk.ToolButtonNew(nil, extApp.Name)
		if err != nil {
			m.logger.Error(err)
			panic(err)
		}
		toolButton.SetName("ea_" + extApp.Name)

		// Create a clicked signal handler for the new button
		toolButton.Connect("clicked", func(button *gtk.ToolButton) {
			name, err := button.GetName()
			if err != nil {
				m.logger.Error(err)
				panic(err)
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

func (m *MainWindow) removeToolbarApplications(toolbar *gtk.Toolbar) {
	// Get toolbar children, and return if there are none
	children := toolbar.GetChildren()
	if children.Length() == 0 {
		return
	}

	// Loop through toolbar children, and remove external applications
	var i uint
	for i = 0; i < children.Length(); i++ {
		// Get toolbar button
		child := children.NthData(i)
		toolButton, ok := child.(*gtk.Widget)
		if !ok {
			continue
		}

		// Get toolbar button name
		name, err := toolButton.GetName()
		if err != nil {
			m.logger.Error(err)
			panic(err)
		}

		// Remove the button if it is an externa application
		if name[:3] == "ea_" {
			toolbar.Remove(toolButton)
		}
	}
	return
}
