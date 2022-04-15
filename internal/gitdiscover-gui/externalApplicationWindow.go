package gitdiscover_gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/softteam/framework"

	"github.com/sirupsen/logrus"
)

type externalApplicationDialog struct {
	window  *gtk.Window
	builder *framework.GtkBuilder
	config  *gitdiscover.Config
	logger  *logrus.Logger

	nameEntry     *gtk.Entry
	commandEntry  *gtk.Entry
	argumentEntry *gtk.Entry

	externalApplication gitdiscover.ExternalApplication
	originalName        string
	mode                externalApplicationModeType

	saveCallback func() bool
}

func newExternalApplicationDialog(logger *logrus.Logger, config *gitdiscover.Config) *externalApplicationDialog {
	// TODO : Send in parent object instead (e) of builder, logger and config
	dialog := new(externalApplicationDialog)
	dialog.config = config
	dialog.logger = logger
	return dialog
}

func (e *externalApplicationDialog) openDialog(parent *gtk.Window, saveCallback func() bool) {
	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("externalApplicationWindow.ui")
	if err != nil {
		panic(err)
	}
	e.builder = builder

	window := e.builder.GetObject("externalApplicationWindow").(*gtk.Window)
	window.Connect("destroy", window.Destroy)
	if e.mode == externalApplicationModeNew {
		window.SetTitle("New external application")
	} else {
		window.SetTitle(fmt.Sprintf("External Application '%s'", e.externalApplication.Name))
	}
	window.SetTransientFor(parent)
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := e.builder.GetObject("saveButton").(*gtk.Button)
	button.Connect("clicked", e.save)
	button = e.builder.GetObject("cancelButton").(*gtk.Button)
	button.Connect("clicked", e.cancel)

	e.nameEntry = e.builder.GetObject("nameEntry").(*gtk.Entry)
	e.commandEntry = e.builder.GetObject("commandEntry").(*gtk.Entry)
	e.argumentEntry = e.builder.GetObject("argumentEntry").(*gtk.Entry)
	if e.mode == externalApplicationModeEdit {
		e.originalName = e.externalApplication.Name

		e.nameEntry.SetText(e.externalApplication.Name)
		e.commandEntry.SetText(e.externalApplication.Command)
		e.argumentEntry.SetText(e.externalApplication.Argument)
	} else {
		e.nameEntry.SetText("")
		e.commandEntry.SetText("")
		e.argumentEntry.SetText("")
	}

	e.saveCallback = saveCallback
	e.window = window
	window.ShowAll()
}

func (e *externalApplicationDialog) save() {
	// TODO : Make sure name is not empty
	text, err := e.nameEntry.GetText()
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	e.externalApplication.Name = text

	text, err = e.commandEntry.GetText()
	// TODO : Make sure command is not empty and exists etc...
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	e.externalApplication.Command = text

	text, err = e.argumentEntry.GetText()
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	e.externalApplication.Argument = text

	if e.saveCallback() {
		e.window.Hide()
	}
}

func (e *externalApplicationDialog) cancel() {
	e.window.Hide()
}
