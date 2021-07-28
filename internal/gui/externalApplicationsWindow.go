package gui

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/gitdiscover/internal/config"
	"github.com/sirupsen/logrus"
)

type ExternalApplicationsWindow struct {
	window  *gtk.Window
	builder *GtkBuilder
	config  *config.Config
	logger  *logrus.Logger
	listBox *gtk.ListBox
}

func NewExternalApplicationsWindow(logger *logrus.Logger, config *config.Config) *ExternalApplicationsWindow {
	window := new(ExternalApplicationsWindow)
	window.config = config
	window.logger = logger
	return window
}

func (e *ExternalApplicationsWindow) openWindow() {
	// Create a new softBuilder
	e.builder = NewGtkBuilder("externalApplicationsWindow.ui", e.logger)

	window := e.builder.getObject("externalApplicationsWindow").(*gtk.Window)
	window.Connect("destroy", window.Destroy)
	window.SetTitle("External Applications...")
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := e.builder.getObject("closeButton").(*gtk.Button)
	button.Connect("clicked", window.Hide)

	// Toolbar
	tool := e.builder.getObject("toolbarAddApplication").(*gtk.ToolButton)
	tool.Connect("clicked", e.addExternalApplication)
	tool = e.builder.getObject("toolbarRemoveApplication").(*gtk.ToolButton)
	tool.Connect("clicked", e.removeExternalApplication)
	tool = e.builder.getObject("toolbarEditApplication").(*gtk.ToolButton)
	tool.Connect("clicked", e.editExternalApplication)

	e.listBox = e.builder.getObject("externalApplicationsList").(*gtk.ListBox)
	e.fillExternalApplicationsList()

	e.window = window
	window.ShowAll()
}

func (e *ExternalApplicationsWindow) fillExternalApplicationsList() {
	e.clearListBox()

	for _, application := range e.config.ExternalApplications {
		item := e.createListItem(application)
		e.listBox.Add(item)
	}

	e.listBox.ShowAll()
}

func (e *ExternalApplicationsWindow) addExternalApplication() {
	dialog := NewExternalApplicationDialog(e.logger, e.config)
	dialog.mode = modeNew
	dialog.openDialog(e.window, func() bool {
		app := config.ExternalApplication{
			Name:     dialog.externalApplication.Name,
			Command:  dialog.externalApplication.Command,
			Argument: dialog.externalApplication.Argument,
		}
		e.config.ExternalApplications = append(e.config.ExternalApplications, app)
		// TODO : Config.Save needs error handling?
		e.config.Save()
		e.fillExternalApplicationsList()
		return true
	})
}

func (e *ExternalApplicationsWindow) removeExternalApplication() {
	app, index := e.getSelectedApplication()
	if app==nil {
		// TODO Please select an application
		return
	}
	// Remove external application from config
	e.config.ExternalApplications = append(e.config.ExternalApplications[:index], e.config.ExternalApplications[index+1:]...)
	e.config.Save()
	e.fillExternalApplicationsList()
}

func (e *ExternalApplicationsWindow) editExternalApplication() {
	dialog := NewExternalApplicationDialog(e.logger, e.config)
	_, index := e.getSelectedApplication()
	if index == -1 {
		// TODO Please select an application
		return
	}
	dialog.externalApplication = e.config.ExternalApplications[index]
	dialog.mode = modeEdit
	dialog.openDialog(e.window, func() bool {
		e.config.ExternalApplications[index].Name = dialog.externalApplication.Name
		e.config.ExternalApplications[index].Command = dialog.externalApplication.Command
		e.config.ExternalApplications[index].Argument = dialog.externalApplication.Argument
		// TODO : Config.Save needs error handling?
		e.config.Save()
		e.fillExternalApplicationsList()
		return true
	})
}

func (e *ExternalApplicationsWindow) clearListBox() {
	children := e.listBox.GetChildren()
	if children == nil {
		return
	}
	var i uint = 0
	for ; i < children.Length(); {
		widget, _ := children.NthData(i).(*gtk.Widget)
		e.listBox.Remove(widget)
		i++
	}
}

func (e *ExternalApplicationsWindow) createListItem(application config.ExternalApplication) *gtk.Box {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}

	labelName, err := gtk.LabelNew("")
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	labelName.SetHAlign(gtk.ALIGN_START)
	applicationName := `<span font="Sans Regular 10" foreground="#44DD44">` + fmt.Sprintf("%-40s", application.Name) + `</span>`
	labelName.SetMarkup(applicationName)
	box.PackStart(labelName, false, false, 10)

	labelCommand, err := gtk.LabelNew("")
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	labelCommand.SetHAlign(gtk.ALIGN_START)
	applicationCommand := `<span font="Sans Regular 10" foreground="#FFFFFF">` + application.Command + `</span>`
	labelCommand.SetMarkup(applicationCommand)
	box.PackStart(labelCommand, false, false, 10)

	labelArgument, err := gtk.LabelNew("")
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	labelArgument.SetHAlign(gtk.ALIGN_START)
	applicationArgument := `<span font="Sans Regular 10" foreground="#FFFFFF">` + application.Argument + `</span>`
	labelArgument.SetMarkup(applicationArgument)
	box.PackEnd(labelArgument, false, false, 10)
	return box
}

func (e *ExternalApplicationsWindow) getSelectedApplication() (*config.ExternalApplication, int) {
	row := e.listBox.GetSelectedRow()
	if row == nil {
		// TODO : MessageBox "Pleade select an application!"
		return nil, -1
	}

	index := row.GetIndex()
	if index < 0 || index >= int(e.listBox.GetChildren().Length()) {
		// TODO : MessageBox "Pleade select an application!"
		return nil, -1
	}

	app := e.config.ExternalApplications[index]
	return &app, index
}
