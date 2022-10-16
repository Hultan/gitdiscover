package gitdiscover_gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
)

type aboutDialog struct {
	dialog *gtk.AboutDialog
	parent *gtk.ApplicationWindow
	logger *logrus.Logger
}

func newAboutDialog(logger *logrus.Logger, parent *gtk.ApplicationWindow) *aboutDialog {
	about := new(aboutDialog)
	about.parent = parent
	about.logger = logger
	return about
}

func (m *aboutDialog) openAboutDialog() {
	if m.dialog == nil {
		// Create a new softBuilder
		builder, err := fw.Gtk.CreateBuilder("about.ui")
		if err != nil {
			panic(err)
		}
		about := builder.GetObject("aboutDialog").(*gtk.AboutDialog)

		about.SetDestroyWithParent(true)
		about.SetTransientFor(m.parent)
		about.SetProgramName(applicationTitle)
		about.SetVersion(applicationVersion)
		about.SetCopyright(applicationCopyRight)
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
