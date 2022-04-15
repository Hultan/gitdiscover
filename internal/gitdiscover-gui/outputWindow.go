package gitdiscover_gui

import (
	"bufio"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"

	"github.com/hultan/softteam/framework"
)

type outputWindow struct {
	builder *framework.GtkBuilder
	logger  *logrus.Logger
	window  *gtk.Window
}

func newOutputWindow(builder *framework.GtkBuilder, logger *logrus.Logger) *outputWindow {
	output := new(outputWindow)
	output.builder = builder
	output.logger = logger
	return output
}

func (o *outputWindow) openWindow(header, text string, gitCommand gitCommandType) {
	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("outputWindow.ui")
	if err != nil {
		panic(err)
	}
	o.builder = builder

	window := o.builder.GetObject("outputWindow").(*gtk.Window)
	window.Connect("destroy", o.closeWindow)
	window.SetTitle("Output window...")
	window.HideOnDelete()
	window.SetModal(true)
	window.SetKeepAbove(true)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	button := o.builder.GetObject("closeButton").(*gtk.Button)
	button.Connect("clicked", o.closeWindow)

	label := o.builder.GetObject("labelHeader").(*gtk.Label)
	if header == "" {
		header = o.getHeader(gitCommand)
	}
	label.SetText(header)

	textView := o.builder.GetObject("textView").(*gtk.TextView)
	buffer, err := gtk.TextBufferNew(nil)
	if err != nil {
		o.logger.Error(err)
		return
	}

	textView.SetBuffer(buffer)
	o.setTextForTextView(text, gitCommand, buffer)
	textView.SetEditable(false)

	o.window = window
	window.ShowAll()
}

func (o *outputWindow) closeWindow() {
	o.window.Hide()
	o.window = nil
}

func (o *outputWindow) setTextForTextView(text string, gitCommand gitCommandType, buffer *gtk.TextBuffer) {
	// Remove illegal characters
	text = o.formatTextGeneral(text)

	// Fix specific formatting
	switch gitCommand {
	case outputGitStatus:
		text = o.formatTextGitStatus(text)
		break
	case outputGitDiff:
		text = o.formatTextGitDiff(text)
		break
	case outputGitLog:
		text = o.formatTextGitLog(text)
		break

	}
	buffer.InsertMarkup(buffer.GetStartIter(), text)
}

func (o *outputWindow) getHeader(gitCommand gitCommandType) string {
	switch gitCommand {
	case outputGitStatus:
		return "git status"
	case outputGitLog:
		return "git log"
	case outputGitDiff:
		return "git diff"
	default:
		return ""
	}
}

func (o *outputWindow) formatTextGeneral(text string) string {
	text = strings.Replace(text, "&", "&amp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)

	return text
}

func (o *outputWindow) formatTextGitStatus(text string) string {
	var result = ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "modified:") {
			result += `<span color="red">` + line + "</span>\n"
			continue
		}
		if strings.HasPrefix(trimmedLine, "deleted:") {
			result += `<span color="red">` + line + "</span>\n"
			continue
		}
		if strings.HasPrefix(trimmedLine, "new file:") {
			result += `<span color="green">` + line + "</span>\n"
			continue
		}
		result += line + "\n"
	}

	if strings.Index(result, `nothing to commit, working tree clean`) >= 0 {
		return result
	}

	if strings.Index(result, `no changes added to commit (use "git add" and/or "git commit -a")`) >= 0 {
		return result
	}

	result = strings.Replace(result, `(use "git add &lt;file&gt;..." to include in what will be committed)`,
		`(use "git add &lt;file&gt;..." to include in what will be committed)<span color="red">`, -1)

	result += "</span>"

	return result
}

func (o *outputWindow) formatTextGitDiff(text string) string {
	var result = ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "+++") {
			result += `<span color="red">` + line + "</span>\n"
			continue
		}
		if strings.HasPrefix(trimmedLine, "-") {
			result += `<span color="red">` + line + "</span>\n"
			continue
		}
		if strings.HasPrefix(trimmedLine, "+") {
			result += `<span color="green">` + line + "</span>\n"
			continue
		}
		result += line + "\n"
	}

	return result
}

func (o *outputWindow) formatTextGitLog(text string) string {
	var result = ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "commit ") {
			result += `<span color="yellow">` + line + "</span>\n"
			continue
		}
		result += line + "\n"
	}

	return result
}
