package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/gitdiscover/internal/gitdiscover"
)

type EditFolderWindow struct {
	mainWindow     *MainWindow
	window         *gtk.Window
	builder        *GtkBuilder
	image          *gtk.Image
	folder         *gitdiscover.RepositoryStatus
	folderIconPath *gtk.Entry
}

func NewEditFolderWindow(mainWindow *MainWindow) *EditFolderWindow {
	edit := new(EditFolderWindow)
	edit.mainWindow = mainWindow
	return edit
}

func (e *EditFolderWindow) openWindow(folder *gitdiscover.RepositoryStatus) {
	// Create a new softBuilder
	e.builder = NewGtkBuilder("editFolderWindow.ui", e.mainWindow.logger)

	window := e.builder.getObject("editFolderWindow").(*gtk.Window)
	window.Connect("destroy", e.closeWindow)
	window.SetTitle("Edit folder...")
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := e.builder.getObject("cancelButton").(*gtk.Button)
	button.Connect("clicked", e.closeWindow)

	button = e.builder.getObject("saveButton").(*gtk.Button)
	button.Connect("clicked", e.save)

	e.folder = folder

	folderPath := e.builder.getObject("folderPathEntry").(*gtk.Entry)
	folderPath.SetText(folder.Path)
	folderPath.SetSensitive(false)

	isGit := e.builder.getObject("isGitFolderCheckBox").(*gtk.CheckButton)
	isGit.SetActive(folder.IsGit)
	isGit.SetSensitive(false)

	e.image = e.builder.getObject("folderIconImage").(*gtk.Image)

	folderIconPath := e.builder.getObject("folderIconPathEntry").(*gtk.Entry)
	folderIconPath.SetText(folder.ImagePath)
	e.folderIconPath = folderIconPath
	e.tryLoadIcon(folder.ImagePath)

	button = e.builder.getObject("selectFolderIconPathButton").(*gtk.Button)
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
		e.folder.Path = path
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
		if e.folder.IsGit {
			iconPath, err = getResourcePath("gitFolder.png")
		} else {
			iconPath, err = getResourcePath("folder.png")
		}
		if err != nil {
			e.mainWindow.logger.Error(err)
			e.image.SetFromPixbuf(nil)
			return
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
	fileDialog.SetCurrentFolder(e.folder.Path)

	if result := fileDialog.Run(); result == gtk.RESPONSE_ACCEPT {
		// Get selected filename.
		e.folderIconPath.SetText(fileDialog.GetFilename())
	}
}
