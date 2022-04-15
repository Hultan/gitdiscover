package gitdiscover_gui

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/softteam/framework"

	"github.com/sirupsen/logrus"
)

// externalApplicationsWindow represents the window for external applications (like nemo, terminal, etc...)
type externalApplicationsWindow struct {
	window  *gtk.Window
	builder *framework.GtkBuilder
	config  *gitdiscover.Config
	logger  *logrus.Logger
	listBox *gtk.ListBox
	refresh func()
}

// newExternalApplicationsWindow creates a new external applications window
func newExternalApplicationsWindow(logger *logrus.Logger, config *gitdiscover.Config) *externalApplicationsWindow {
	window := new(externalApplicationsWindow)
	window.config = config
	window.logger = logger
	return window
}

func (e *externalApplicationsWindow) openWindow(refresh func()) {
	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("externalApplicationsWindow.ui")
	if err != nil {
		panic(err)
	}
	e.builder = builder

	window := e.builder.GetObject("externalApplicationsWindow").(*gtk.Window)
	window.Connect("destroy", e.closeWindow)
	window.SetTitle("External Applications...")
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := e.builder.GetObject("closeButton").(*gtk.Button)
	button.Connect("clicked", e.closeWindow)

	// Toolbar
	tool := e.builder.GetObject("toolbarAddApplication").(*gtk.ToolButton)
	tool.Connect("clicked", e.addExternalApplication)
	tool = e.builder.GetObject("toolbarRemoveApplication").(*gtk.ToolButton)
	tool.Connect("clicked", e.removeExternalApplication)
	tool = e.builder.GetObject("toolbarEditApplication").(*gtk.ToolButton)
	tool.Connect("clicked", e.editExternalApplication)

	e.listBox = e.builder.GetObject("externalApplicationsList").(*gtk.ListBox)
	e.listBox.SetActivateOnSingleClick(false)
	e.listBox.Connect("row-activated", func(listbox *gtk.ListBox, row *gtk.ListBoxRow) {
		index := row.GetIndex()
		e.editExternalApplicationByIndex(index)
	})

	e.fillExternalApplicationsList()

	e.window = window
	e.refresh = refresh
	window.ShowAll()
}

func (e *externalApplicationsWindow) closeWindow() {
	e.window.Hide()
	e.refresh()
}

func (e *externalApplicationsWindow) fillExternalApplicationsList() {
	e.clearListBox()

	sgName, _ := gtk.SizeGroupNew(gtk.SIZE_GROUP_BOTH)
	sgCommand, _ := gtk.SizeGroupNew(gtk.SIZE_GROUP_BOTH)
	sgArgument, _ := gtk.SizeGroupNew(gtk.SIZE_GROUP_BOTH)

	for _, application := range e.config.ExternalApplications {
		item := e.createListItem(application, sgName, sgCommand, sgArgument)
		e.listBox.Add(item)
	}

	e.listBox.ShowAll()
}

func (e *externalApplicationsWindow) addExternalApplication() {
	dialog := newExternalApplicationDialog(e.logger, e.config)
	dialog.mode = externalApplicationModeNew
	dialog.openDialog(e.window, func() bool {
		app := gitdiscover.ExternalApplication{
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

func (e *externalApplicationsWindow) removeExternalApplication() {
	app, index := e.getSelectedApplication()
	if app == nil {
		// TODO Please select an application
		return
	}
	// Remove external application from config
	e.config.ExternalApplications = append(e.config.ExternalApplications[:index], e.config.ExternalApplications[index+1:]...)
	e.config.Save()
	e.fillExternalApplicationsList()
}

func (e *externalApplicationsWindow) editExternalApplication() {
	_, index := e.getSelectedApplication()
	if index == -1 {
		// TODO Please select an application
		return
	}
	e.editExternalApplicationByIndex(index)
}

func (e *externalApplicationsWindow) editExternalApplicationByIndex(index int) {
	dialog := newExternalApplicationDialog(e.logger, e.config)
	dialog.externalApplication = e.config.ExternalApplications[index]
	dialog.mode = externalApplicationModeEdit
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

func (e *externalApplicationsWindow) clearListBox() {
	children := e.listBox.GetChildren()
	if children == nil {
		return
	}
	var i uint = 0
	for i < children.Length() {
		widget, _ := children.NthData(i).(*gtk.Widget)
		e.listBox.Remove(widget)
		i++
	}
}

func (e *externalApplicationsWindow) createListItem(application gitdiscover.ExternalApplication,
	sgName, sgCommand, sgArgument *gtk.SizeGroup) *gtk.Box {

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}

	// Name
	labelName, err := gtk.LabelNew("")
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	applicationName := `<span font="Sans Regular 10" foreground="#44DD44">` + application.Name + `</span>`
	labelName.SetMarkup(applicationName)
	labelName.SetXAlign(0.0)
	sgName.AddWidget(labelName)
	box.PackStart(labelName, true, true, 10)

	// Command
	labelCommand, err := gtk.LabelNew("")
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	labelCommand.SetText(application.Command)
	labelCommand.SetXAlign(0.0)
	sgCommand.AddWidget(labelCommand)
	box.PackStart(labelCommand, true, true, 10)

	// Argument
	labelArgument, err := gtk.LabelNew("")
	if err != nil {
		e.logger.Error(err)
		panic(err)
	}
	labelArgument.SetText(application.Argument)
	labelArgument.SetXAlign(0.0)
	sgArgument.AddWidget(labelArgument)
	box.PackStart(labelArgument, true, true, 10)
	return box
}

func (e *externalApplicationsWindow) getSelectedApplication() (*gitdiscover.ExternalApplication, int) {
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
