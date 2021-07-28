package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
)

type aboutDialog struct {
	dialog *gtk.AboutDialog
	parent *gtk.ApplicationWindow
	logger *logrus.Logger
}

func NewAboutDialog(logger *logrus.Logger, parent *gtk.ApplicationWindow) *aboutDialog{
	about := new(aboutDialog)
	about.parent = parent
	about.logger = logger
	return about
}

func (m *aboutDialog) openAboutDialog() {
	if m.dialog == nil {
		// Create a new softBuilder
		builder := NewGtkBuilder("about.ui", m.logger)
		about := builder.getObject("aboutDialog").(*gtk.AboutDialog)

		about.SetDestroyWithParent(true)
		about.SetTransientFor(m.parent)
		about.SetProgramName(ApplicationTitle)
		about.SetVersion(ApplicationVersion)
		about.SetCopyright(ApplicationCopyRight)
		about.SetComments("Discover your GIT repositories...")
		about.SetModal(true)
		about.SetPosition(gtk.WIN_POS_CENTER)

		_ = about.Connect("response", func(dialog *gtk.AboutDialog, responseId gtk.ResponseType) {
			if responseId == gtk.RESPONSE_CANCEL || responseId == gtk.RESPONSE_DELETE_EVENT {
				about.Hide()
			}
		})

		m.dialog = about
	}

	m.dialog.ShowAll()
}

