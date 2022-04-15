package gitdiscover_gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type infoBar struct {
	infoBar      *gtk.InfoBar
	labelInfoBar *gtk.Label
}

func newInfoBar(bar *gtk.InfoBar, labelInfoBar *gtk.Label) *infoBar {
	info := new(infoBar)
	info.infoBar = bar
	info.labelInfoBar = labelInfoBar
	return info
}

func (i *infoBar) showInfo(text string) {
	i.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()
}

func (i *infoBar) showInfoWithTimeout(text string, seconds uint) {
	i.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()

	glib.TimeoutSecondsAdd(seconds, func() {
		i.hideInfoBar()
	})
}

func (i *infoBar) showError(text string) {
	i.infoBar.SetMessageType(gtk.MESSAGE_ERROR)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()
}

func (i *infoBar) hideInfoBar() {
	i.infoBar.Hide()
}
