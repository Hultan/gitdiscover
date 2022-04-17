package gitdiscover_gui

import (
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

type popupMenu struct {
	mainWindow *MainWindow
	popupMenu  *gtk.Menu

	popupAddFolder            *gtk.MenuItem
	popupEditFolder           *gtk.MenuItem
	popupRemoveFolder         *gtk.MenuItem
	popupFavorite             *gtk.MenuItem
	popupExternalApplications *gtk.MenuItem
	popupGitStatus            *gtk.MenuItem
	popupGitDiff              *gtk.MenuItem
	popupGitLog               *gtk.MenuItem
	popupGit                  *gtk.MenuItem
}

func newPopupMenu(window *MainWindow) *popupMenu {
	menu := new(popupMenu)
	menu.mainWindow = window
	return menu
}

func (p *popupMenu) setupPopupMenu() {
	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("mainWindow.ui")
	if err != nil {
		panic(err)
	}

	p.popupMenu = builder.GetObject("popupMenu").(*gtk.Menu)

	p.popupAddFolder = builder.GetObject("popupAddFolder").(*gtk.MenuItem)
	p.popupEditFolder = builder.GetObject("popupEditFolder").(*gtk.MenuItem)
	p.popupRemoveFolder = builder.GetObject("popupRemoveFolder").(*gtk.MenuItem)
	p.popupFavorite = builder.GetObject("popupFavorite").(*gtk.MenuItem)
	p.popupExternalApplications = builder.GetObject("popupExternalApplications").(*gtk.MenuItem)
	p.popupGit = builder.GetObject("popupGit").(*gtk.MenuItem)
	p.popupGitStatus = builder.GetObject("popupGitStatus").(*gtk.MenuItem)
	p.popupGitDiff = builder.GetObject("popupGitDiff").(*gtk.MenuItem)
	p.popupGitLog = builder.GetObject("popupGitLog").(*gtk.MenuItem)

	p.setupEvents()
}

func (p *popupMenu) setupEvents() {
	_ = p.mainWindow.window.Connect("button-release-event", func(window *gtk.ApplicationWindow, event *gdk.Event) {
		// If right mouse button is NOT pressed, return
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() != gdk.BUTTON_SECONDARY {
			return
		}

		// Get the currently selected repo
		repo := p.mainWindow.getSelectedRepo()
		if repo == nil {
			p.mainWindow.infoBar.showInfoWithTimeout("Please select a repo...", 5)
			return
		}

		// Disable the git menu for non-git folders
		p.popupGit.SetSensitive(repo.IsGit())

		// Create a sub menu for external applications
		menu, err := gtk.MenuNew()
		if err != nil {
			p.mainWindow.logger.Error(err)
			return
		}

		// Create menu items for external applications
		for i := 0; i < len(p.mainWindow.discover.ExternalApplications); i++ {
			app := p.mainWindow.discover.GetExternalApplicationByIndex(i)
			item, err := gtk.MenuItemNew()
			if err != nil {
				p.mainWindow.logger.Error(err)
				continue
			}
			item.SetLabel(app.Name)
			menu.Add(item)
			item.Connect("activate", func() {
				repo := p.mainWindow.getSelectedRepo()
				p.mainWindow.openInExternalApplication(app.Name, repo)
			})
		}
		p.popupExternalApplications.SetSubmenu(menu)
		p.popupExternalApplications.ShowAll()

		p.popupMenu.PopupAtPointer(event)
	})

	p.popupAddFolder.Connect("activate", func() {
		p.mainWindow.addRepositoryButtonClicked()
	})

	p.popupEditFolder.Connect("activate", func() {
		p.mainWindow.editRepositoryButtonClicked()
	})

	p.popupRemoveFolder.Connect("activate", func() {
		p.mainWindow.removeRepositoryButtonClicked()
	})

	p.popupFavorite.Connect("activate", func() {
		// Get the currently selected repo
		repo := p.mainWindow.getSelectedRepo()

		if repo == nil {
			p.mainWindow.infoBar.showInfoWithTimeout("Please select a repo...", 5)
			return
		}

		repo.SetIsFavorite(!repo.IsFavorite())
		p.mainWindow.discover.Save()
		p.mainWindow.refreshRepositoryList()
	})

	p.popupGitStatus.Connect("activate", func() {
		p.runGitCommand("git status", outputGitStatus)
	})

	p.popupGitDiff.Connect("activate", func() {
		p.runGitCommand("git diff", outputGitDiff)
	})

	p.popupGitLog.Connect("activate", func() {
		p.runGitCommand("git log", outputGitLog)
	})
}

// runGitCommand : Run a GIT command
func (p *popupMenu) runGitCommand(command string, outputType gitCommandType) {
	// Get the currently selected repo
	repo := p.mainWindow.getSelectedRepo()
	if repo == nil {
		p.mainWindow.infoBar.showInfoWithTimeout("Please select a repo...", 5)
		return
	}

	// Create a temp bash file, with the command
	file, err := p.createFile(repo.Path(), command)
	if err != nil {
		p.mainWindow.logger.Error(err)
		p.mainWindow.infoBar.showError(err.Error())
		return
	}

	// Execute the bash file
	result := p.mainWindow.executeCommand("/bin/sh", file)
	output := newOutputWindow(p.mainWindow.builder, p.mainWindow.logger)
	output.openWindow("", result, outputType)

	// Clean up
	err = os.Remove(file)
	if err != nil {
		p.mainWindow.logger.Error(err)
		p.mainWindow.infoBar.showError(err.Error())
	}
}

// createFile : Create a temp bash file
func (p *popupMenu) createFile(path, command string) (string, error) {
	// Create the text for the bash file
	text := "#!/bin/sh\n"
	text += "cd " + path + "\n"
	text += command

	// Create a temp file
	content := []byte(text)
	tmpfile, err := ioutil.TempFile("", "gitdiscover.*.sh")
	if err != nil {
		p.mainWindow.logger.Error(err)
		return "", err
	}

	// Write to the file
	_, err = tmpfile.Write(content)
	if err != nil {
		p.mainWindow.logger.Error(err)
		return "", err
	}

	// Clean up
	err = tmpfile.Close()
	if err != nil {
		p.mainWindow.logger.Error(err)
		return "", err
	}

	return tmpfile.Name(), nil
}
