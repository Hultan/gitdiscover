package gitdiscover_gui

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/gitdiscover/internal/gitdiscover"
)

func (m *MainWindow) getMarkup(text, color string) string {
	markup := fmt.Sprintf(`<span font="Sans Regular 10" foreground="#%s">`, color)
	markup += text
	markup += `</span>`
	return markup
}

func (m *MainWindow) getSelectedRepo() *gitdiscover.TrackedFolder {
	row := m.repositoryListBox.GetSelectedRow()
	if row == nil {
		return nil
	}
	boxObj, err := row.GetChild()
	if err != nil {
		m.infoBar.ShowError(err.Error())
		return nil
	}
	box, ok := boxObj.(*gtk.Box)
	if !ok {
		m.infoBar.ShowError("Failed to convert to *gtk.Widget")
		return nil
	}
	name, err := box.GetName()
	if err != nil {
		m.infoBar.ShowError(err.Error())
		return nil
	}
	if name == "sep" {
		return nil
	}
	indexString := name[4:]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		m.infoBar.ShowError(err.Error())
		return nil
	}
	repo := m.tracker.Folders[index]

	return repo
}

func (m *MainWindow) openConfig() {
	// Open the config file in the text editor
	go func() {
		m.executeCommand("xed", m.config.GetConfigPath())
	}()
}

func (m *MainWindow) openLog() {
	// Open the log file in the text editor
	go func() {
		m.executeCommand("xed", m.ApplicationLogPath)
	}()
}

func (m *MainWindow) executeCommand(command, arguments string) string {
	cmd := exec.Command(command, arguments)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	// https://stackoverflow.com/questions/62853835/how-to-use-syscall-sysprocattr-struct-fields-for-windows-when-os-is-set-for-linu
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	//	Setpgid: true,
	//	Pgid:    0,
	// }

	// set the output to our variable
	out, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.Error("Failed to open external application: ", command, " ", arguments)
		m.logger.Error(err)
		m.infoBar.ShowError(err.Error())
		return ""
	}

	return string(out)
}
