package adaptive

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

// ErrorLabel is a label that displays the short form of an error but allows the
// user to get the full error from the UI directly.
type ErrorLabel struct {
	*gtk.Box

	Short      *gtk.ToggleButton
	ShortLabel *gtk.Label
	ShortIcon  *gtk.Image

	Reveal *gtk.Revealer
	Full   *gtk.Label

	RevealedIcon  string
	CollapsedIcon string
}

// NewErrorLabel creates a new error label from the given error. If err is nil,
// then the function panics.
func NewErrorLabel(err error) *ErrorLabel {
	if err == nil {
		panic("unexpected nil error given to NewErrorLabel")
	}

	return NewErrorLabelFull("Error: "+shortError(err), expandError(err))
}

// NewErrorLabelFull creates a new error label from two strings.
func NewErrorLabelFull(short, full string) *ErrorLabel {
	l := &ErrorLabel{
		RevealedIcon:  "pan-down-symbolic",
		CollapsedIcon: "pan-end-symbolic",
	}

	l.Full = gtk.NewLabel(full)
	l.Full.AddCSSClass("adaptive-errorlabel-full")
	l.Full.SetSelectable(true)
	l.Full.SetXAlign(0)
	l.Full.SetWrap(true)
	l.Full.SetWrapMode(pango.WrapWordChar)

	l.Reveal = gtk.NewRevealer()
	l.Reveal.SetChild(l.Full)
	l.Reveal.SetTransitionType(gtk.RevealerTransitionTypeSlideDown)
	l.Reveal.SetTransitionDuration(250)
	l.Reveal.SetRevealChild(false)

	l.ShortLabel = gtk.NewLabel(short)
	l.ShortLabel.SetXAlign(0)
	l.ShortLabel.SetWrap(true)
	l.ShortLabel.SetWrapMode(pango.WrapWordChar)

	l.ShortIcon = gtk.NewImageFromIconName(l.CollapsedIcon)
	l.ShortIcon.SetVAlign(gtk.AlignStart)
	l.ShortIcon.SetIconSize(gtk.IconSizeNormal)

	shortBox := gtk.NewBox(gtk.OrientationHorizontal, 0)
	shortBox.Append(l.ShortIcon)
	shortBox.Append(l.ShortLabel)

	l.Short = gtk.NewToggleButtonWithLabel(short)
	l.Short.AddCSSClass("adaptive-errorlabel-button")
	l.Short.SetTooltipText("")
	l.Short.SetHasFrame(false)
	l.Short.SetChild(shortBox)

	l.Short.ConnectClicked(func() {
		reveal := l.Short.Active()
		l.Reveal.SetRevealChild(reveal)
		if reveal {
			l.ShortIcon.SetFromIconName(l.RevealedIcon)
		} else {
			l.ShortIcon.SetFromIconName(l.CollapsedIcon)
		}
	})

	l.Box = gtk.NewBox(gtk.OrientationVertical, 0)
	l.Box.AddCSSClass("adaptive-errorlabel")
	l.Box.Append(l.Short)
	l.Box.Append(l.Reveal)

	return l
}

func shortError(err error) string {
	error := err.Error()
	if error == "" {
		return ""
	}
	return strings.Split(error, ": ")[0]
}

func expandError(err error) string {
	error := err.Error()
	parts := strings.SplitAfter(error, ": ")
	if len(parts) > 1 {
		parts = parts[1:]
	}

	b := strings.Builder{}
	b.Grow(len(error) * 2)

	for i, part := range parts {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(strings.Repeat("    ", i))
		b.WriteString("└ ")

		part = strings.TrimSpace(part)
		b.WriteString(capitalize(part))

		if i == len(parts)-1 {
			if !strings.HasSuffix(part, ".") {
				b.WriteByte('.')
			}
		} else {
			if !strings.HasSuffix(part, ":") {
				b.WriteByte(':')
			}
		}
	}

	return b.String()
}

func capitalize(sentence string) string {
	r, sz := utf8.DecodeRuneInString(sentence)
	if sz > -1 {
		return string(unicode.ToUpper(r)) + sentence[sz:]
	}
	return sentence
}
