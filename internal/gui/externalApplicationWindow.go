package gui

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/gitdiscover/internal/config"
	"github.com/sirupsen/logrus"
)

type ModeType int

const (
	modeNew  ModeType = 0
	modeEdit          = 1
)

type ExternalApplicationDialog struct {
	window  *gtk.Window
	builder *GtkBuilder
	config  *config.Config
	logger  *logrus.Logger

	nameEntry     *gtk.Entry
	commandEntry  *gtk.Entry
	argumentEntry *gtk.Entry

	externalApplication config.ExternalApplication
	originalName string
	mode                ModeType

	saveCallback func() bool
}

func NewExternalApplicationDialog(logger *logrus.Logger, config *config.Config) *ExternalApplicationDialog {
	// TODO : Send in parent object instead (e) of builder, logger and config
	dialog := new(ExternalApplicationDialog)
	dialog.config = config
	dialog.logger = logger
	return dialog
}

func (e *ExternalApplicationDialog) openDialog(parent *gtk.Window, saveCallback func() bool) {
	// Create a new softBuilder
	e.builder = NewGtkBuilder("externalApplicationWindow.ui", e.logger)

	window := e.builder.getObject("externalApplicationWindow").(*gtk.Window)
	window.Connect("destroy", window.Destroy)
	if e.mode == modeNew {
		window.SetTitle("New external application")
	} else {
		window.SetTitle(fmt.Sprintf("External Application '%s'", e.externalApplication.Name))
	}
	window.SetTransientFor(parent)
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := e.builder.getObject("saveButton").(*gtk.Button)
	button.Connect("clicked", e.save)
	button = e.builder.getObject("cancelButton").(*gtk.Button)
	button.Connect("clicked", e.cancel)

	e.nameEntry = e.builder.getObject("nameEntry").(*gtk.Entry)
	e.commandEntry = e.builder.getObject("commandEntry").(*gtk.Entry)
	e.argumentEntry = e.builder.getObject("argumentEntry").(*gtk.Entry)
	if e.mode == modeEdit {
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

func (e *ExternalApplicationDialog) save() {
	// TODO : Make sure name is not empty
	text, err := e.nameEntry.GetText()
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	e.externalApplication.Name = text

	text, err = e.commandEntry.GetText()
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

func (e *ExternalApplicationDialog) cancel() {
	e.window.Hide()
}
