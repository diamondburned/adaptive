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
	page := &StatusPage{}

	page.Icon = gtk.NewImage()
	page.Icon.SetIconSize(gtk.IconSizeLarge)
	page.Icon.AddCSSClass("adaptive-statuspage-icon")

	page.Title = gtk.NewLabel("")
	page.Title.AddCSSClass("adaptive-statuspage-title")
	page.Title.SetEllipsize(pango.EllipsizeEnd)
	page.Title.SetSingleLineMode(true)

	page.Grid = gtk.NewGrid()
	page.Grid.SetHExpand(true)
	page.Grid.SetVExpand(true)
	page.Grid.SetHAlign(gtk.AlignCenter)
	page.Grid.SetVAlign(gtk.AlignCenter)
	page.Grid.AddCSSClass("adaptive-statuspage")

	return page
}

func (p *StatusPage) ensureIcon() {
	p.Grid.Attach(p.Icon, 0, 0, 1, 1)
}

func (p *StatusPage) ensureTitle() {
	p.Grid.Attach(p.Title, 0, 1, 1, 1)
}

func (p *StatusPage) ensureDescription(desc gtk.Widgetter) {
	p.Grid.Attach(desc, 0, 2, 1, 1)
}

// SetTitle ensures the title is in the page and sets its content.
func (p *StatusPage) SetTitle(title string) {
	p.ensureTitle()
	p.Title.SetText(title)
	p.Title.SetTooltipText(title)
}

// SetDescription ensures the description is in the page and sets its content.
func (p *StatusPage) SetDescription(desc gtk.Widgetter) {
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
	p.ensureIcon()
	p.Icon.SetFromIconName(icon)
}
