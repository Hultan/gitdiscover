package gitdiscover_gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
	"github.com/hultan/softteam/framework"
)

type EditFolderWindow struct {
	mainWindow     *MainWindow
	window         *gtk.Window
	builder        *framework.GtkBuilder
	image          *gtk.Image
	folder         *gitdiscover.TrackedFolder
	folderIconPath *gtk.Entry
}

func NewEditFolderWindow(mainWindow *MainWindow) *EditFolderWindow {
	edit := new(EditFolderWindow)
	edit.mainWindow = mainWindow
	return edit
}

func (e *EditFolderWindow) openWindow(folder *gitdiscover.TrackedFolder) {
	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("editFolderWindow.ui")
	if err != nil {
		panic(err)
	}
	e.builder = builder

	window := e.builder.GetObject("editFolderWindow").(*gtk.Window)
	window.Connect("destroy", e.closeWindow)
	window.SetTitle("Edit folder...")
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := e.builder.GetObject("cancelButton").(*gtk.Button)
	button.Connect("clicked", e.closeWindow)

	button = e.builder.GetObject("saveButton").(*gtk.Button)
	button.Connect("clicked", e.save)

	e.folder = folder

	folderPath := e.builder.GetObject("folderPathEntry").(*gtk.Entry)
	folderPath.SetText(folder.Path())
	folderPath.SetSensitive(false)

	isGit := e.builder.GetObject("isGitFolderCheckBox").(*gtk.CheckButton)
	isGit.SetActive(folder.IsGit())
	isGit.SetSensitive(false)

	e.image = e.builder.GetObject("folderIconImage").(*gtk.Image)

	folderIconPath := e.builder.GetObject("folderIconPathEntry").(*gtk.Entry)
	folderIconPath.SetText(folder.ImagePath())
	e.folderIconPath = folderIconPath
	e.tryLoadIcon(folder.ImagePath())

	button = e.builder.GetObject("selectFolderIconPathButton").(*gtk.Button)
	button.Connect("clicked", func() {
		e.chooseIcon()

	})

	e.window = window
	window.ShowAll()
}

func (e *EditFolderWindow) save() {
	path, err := e.folderIconPath.GetText()
	if err != nil {
		e.mainWindow.logger.Error(err)
	} else {
		e.folder.SetPath(path)
		e.mainWindow.config.Save()
	}
	e.closeWindow()
}

func (e *EditFolderWindow) closeWindow() {
	e.window.Hide()
	e.window = nil
}

func (e *EditFolderWindow) tryLoadIcon(path string) {
	pix, err := gdk.PixbufNewFromFileAtSize(path, 16, 16)
	if err != nil {
		e.mainWindow.logger.Error(err)
		var iconPath = ""
		fw := framework.NewFramework()
		if e.folder.IsGit() {
			iconPath = fw.Resource.GetResourcePath("gitFolder.png")
		} else {
			iconPath = fw.Resource.GetResourcePath("folder.png")
		}
		pix, err = gdk.PixbufNewFromFileAtSize(iconPath, 16, 16)
		if err != nil {
			e.mainWindow.logger.Error(err)
			e.image.SetFromPixbuf(nil)
			return
		}
	}
	e.image.SetFromPixbuf(pix)
}

func (e *EditFolderWindow) chooseIcon() {
	fileDialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Choose file...",
		e.window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_DELETE_EVENT,
		"Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		// TODO
	}
	defer fileDialog.Destroy()

	fileFilter, err := gtk.FileFilterNew()
	if err != nil {
		// TODO
	}
	fileFilter.SetName("Image files")
	fileFilter.AddPattern("*.png")
	fileFilter.AddPattern("*.bmp")
	fileFilter.AddPattern("*.ico")
	fileDialog.AddFilter(fileFilter)
	fileDialog.SetCurrentFolder(e.folder.Path())

	if result := fileDialog.Run(); result == gtk.RESPONSE_ACCEPT {
		// Get selected filename.
		e.folderIconPath.SetText(fileDialog.GetFilename())
	}
}
