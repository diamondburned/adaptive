package adaptive

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

// StatusPage is a widget component that contains an icon, a title and a
// description, which are all optional.
type StatusPage struct {
	*gtk.Grid
	Icon  *gtk.Image
	Title *gtk.Label
}

// NewStatusPage creates a new empty status page. All its widgets are properly
// initialized, but they're not added into the box until set.
func NewStatusPage() *StatusPage {
	p := &StatusPage{}

	p.Icon = gtk.NewImage()
	p.Icon.Hide()
	p.Icon.SetIconSize(gtk.IconSizeLarge)
	p.Icon.AddCSSClass("adaptive-statuspage-icon")

	p.Title = gtk.NewLabel("")
	p.Title.AddCSSClass("adaptive-statuspage-title")
	p.Title.Hide()
	p.Title.SetEllipsize(pango.EllipsizeEnd)
	p.Title.SetSingleLineMode(true)

	p.Grid = gtk.NewGrid()
	p.Grid.AddCSSClass("adaptive-statuspage")
	p.Grid.SetHExpand(false)
	p.Grid.SetVExpand(false)
	p.Grid.SetHAlign(gtk.AlignCenter)
	p.Grid.SetVAlign(gtk.AlignCenter)
	p.Grid.Attach(p.Icon, 0, 0, 1, 1)
	p.Grid.Attach(p.Title, 0, 1, 1, 1)

	return p
}

func (p *StatusPage) ensureIcon() {
	p.Icon.Show()
}

func (p *StatusPage) ensureTitle() {
	p.Title.Show()
}

func (p *StatusPage) ensureDescription(desc gtk.Widgetter) {
	p.Grid.RemoveRow(2)
	p.Grid.Attach(desc, 0, 2, 1, 1)
}

// SetTitle ensures the title is in the page and sets its content.
func (p *StatusPage) SetTitle(title string) {
	if title == "" {
		if p.Title != nil {
			p.Grid.Remove(p.Title)
			p.Title = nil
		}
		return
	}

	p.ensureTitle()
	p.Title.SetText(title)
	p.Title.SetTooltipText(title)
}

// SetDescription ensures the description is in the page and sets its content.
func (p *StatusPage) SetDescription(desc gtk.Widgetter) {
	if desc == nil {
		p.Grid.RemoveRow(2)
		return
	}
	p.ensureDescription(desc)
}

// SetDescriptionText calls SetDescription with a new description label. The
// label is justified to the middle and has a 50 characters wide width cap.
func (p *StatusPage) SetDescriptionText(desc string) {
	description := gtk.NewLabel(desc)
	description.AddCSSClass("adaptive-statuspage-description")
	description.SetSelectable(true)
	description.SetWrap(true)
	description.SetWrapMode(pango.WrapWordChar)
	description.SetJustify(gtk.JustifyCenter)
	description.SetMaxWidthChars(50)

	p.SetDescription(description)
}

// SetIconName ensures the icon is in the page and sets its icon name.
func (p *StatusPage) SetIconName(icon string) {
	if icon == "" {
		if p.Icon != nil {
			p.Grid.Remove(p.Icon)
			p.Icon = nil
		}
		return
	}

	p.ensureIcon()
	p.Icon.SetFromIconName(icon)
}
