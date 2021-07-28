package gui

import "github.com/gotk3/gotk3/gtk"

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

func (i *InfoBar) ShowError(text string) {
	i.infoBar.SetMessageType(gtk.MESSAGE_ERROR)
	i.labelInfoBar.SetText(text)
	i.infoBar.ShowAll()
}

func (i *InfoBar) hideInfoBar() {
	i.infoBar.Hide()
}
