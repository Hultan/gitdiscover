package gitdiscover_gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type InfoBar struct {
	infoBar      *gtk.InfoBar
	labelInfoBar *gtk.Label
}

func NewInfoBar(infoBar *gtk.InfoBar, labelInfoBar *gtk.Label) *InfoBar {
	info := new(InfoBar)
	info.infoBar = infoBar
	info.labelInfoBar = labelInfoBar
	return info
}

func (i *InfoBar) ShowInfo(text string) {
	i.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()
}

func (i *InfoBar) ShowInfoWithTimeout(text string, seconds uint) {
	i.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()

	glib.TimeoutSecondsAdd(seconds, func() {
		i.HideInfoBar()
	})
}

func (i *InfoBar) ShowError(text string) {
	i.infoBar.SetMessageType(gtk.MESSAGE_ERROR)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()
}

func (i *InfoBar) HideInfoBar() {
	i.infoBar.Hide()
}
