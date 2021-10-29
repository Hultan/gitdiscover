package gui

import (
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

type PopupMenu struct {
	mainWindow *MainWindow
	popupMenu  *gtk.Menu

	popupAddFolder            *gtk.MenuItem
	popupEditFolder			  *gtk.MenuItem
	popupRemoveFolder         *gtk.MenuItem
	popupExternalApplications *gtk.MenuItem
	popupGitStatus            *gtk.MenuItem
	popupGitDiff              *gtk.MenuItem
	popupGitLog               *gtk.MenuItem
}

func NewPopupMenu(window *MainWindow) *PopupMenu {
	menu := new(PopupMenu)
	menu.mainWindow = window
	return menu
}

func (p *PopupMenu) Setup() {
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
	p.popupExternalApplications = builder.GetObject("popupExternalApplications").(*gtk.MenuItem)
	p.popupGitStatus = builder.GetObject("popupGitStatus").(*gtk.MenuItem)
	p.popupGitDiff = builder.GetObject("popupGitDiff").(*gtk.MenuItem)
	p.popupGitLog = builder.GetObject("popupGitLog").(*gtk.MenuItem)

	p.setupEvents()
}

func (p *PopupMenu) setupEvents() {
	_ = p.mainWindow.window.Connect("button-release-event", func(window *gtk.ApplicationWindow, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() == gdk.BUTTON_SECONDARY {
			//list := p.popupExternalApplications.GetChildren()

			menu, err := gtk.MenuNew()
			if err != nil {
				p.mainWindow.logger.Error(err)
			} else {
				for i := 0; i < len(p.mainWindow.config.ExternalApplications); i++ {
					app := p.mainWindow.config.ExternalApplications[i]
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
			}
			p.popupMenu.PopupAtPointer(event)
		}
	})

	p.popupAddFolder.Connect("activate", func() {
		p.mainWindow.addButtonClicked()
	})

	p.popupEditFolder.Connect("activate", func() {
		p.mainWindow.editButtonClicked()
	})

	p.popupRemoveFolder.Connect("activate", func() {
		p.mainWindow.removeButtonClicked()
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

func (p *PopupMenu) createFile(path, command string) (string, error) {
	text := "#!/bin/sh\n"
	text += "cd " + path + "\n"
	text += command

	content := []byte(text)
	tmpfile, err := ioutil.TempFile("", "gitdiscover.*.sh")
	if err != nil {
		p.mainWindow.logger.Error(err)
		return "", err
	}

	// clean up
	if _, err := tmpfile.Write(content); err != nil {
		p.mainWindow.logger.Error(err)
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		p.mainWindow.logger.Error(err)
		return "", err
	}

	return tmpfile.Name(), nil
}

func (p *PopupMenu) runGitCommand(command string, outputType OutputType) {
	repo := p.mainWindow.getSelectedRepo()
	if repo == nil {
		p.mainWindow.infoBar.ShowInfo("Please select a repo...")
		return
	}
	p.mainWindow.infoBar.hideInfoBar()

	file, err := p.createFile(repo.Path, command)
	if err != nil {
		p.mainWindow.logger.Error(err)
		p.mainWindow.infoBar.ShowError(err.Error())
		return
	}

	result := p.mainWindow.executeCommand("/bin/sh", file)
	output := NewOutputWindow(p.mainWindow.builder, p.mainWindow.logger)
	output.openWindow("", result, outputType)

	err = os.Remove(file)
	if err != nil {
		p.mainWindow.logger.Error(err)
		p.mainWindow.infoBar.ShowError(err.Error())
	}
}
