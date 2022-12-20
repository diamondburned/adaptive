package adaptive

import (
	"log"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

// TransformInitials transforms the given input string into initials.
func TransformInitials(in string) string {
	if in == "" {
		return ""
	}

	runes := []rune(in)

	for _, r := range runes {
		if unicode.IsUpper(r) {
			return transformInitialsCap(runes)
		}
	}

	return transformInitials(runes)
}

func transformInitialsCap(runes []rune) string {
	return runesMap(runes, 2, func(r rune) rune {
		if unicode.IsUpper(r) {
			return r
		}
		return -1
	})
}

func transformInitials(runes []rune) string {
	return runesMap(runes, 1, func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return -1
	})
}

func runesMap(runes []rune, until int, f func(rune) rune) string {
	b := strings.Builder{}
	b.Grow(until * utf8.UTFMax)

	for i, l := 0, 0; i < len(runes) && l < until; i++ {
		if r := f(runes[i]); r > -1 {
			l++
			b.WriteRune(r)
		}
	}

	return b.String()
}

// Avatar wraps around an Image and makes it appear round.
type Avatar struct {
	*Bin
	Image *gtk.Image
	Label *gtk.Label // for initials

	labelAttrs   *pango.AttrList
	initialsFunc func(string) string
	typeClass    string
}

// NewAvatar creates a new round image. If radius is 0, then it will be half the
// dimensions. If the radius is less than 0, then nothing is rounded.
func NewAvatar(size int) *Avatar {
	avatar := &Avatar{
		Bin:          NewBin(),
		Image:        gtk.NewImage(),
		Label:        gtk.NewLabel(""),
		initialsFunc: TransformInitials,
	}

	avatar.Image.SetOverflow(gtk.OverflowHidden)
	avatar.Label.SetOverflow(gtk.OverflowHidden)

	avatar.SetHExpand(false)
	avatar.SetVExpand(false)
	avatar.SetHAlign(gtk.AlignCenter)
	avatar.SetVAlign(gtk.AlignCenter)
	avatar.SetSizeRequest(size)
	avatar.SetChild(avatar.Image)
	avatar.AddCSSClass("adaptive-avatar")
	avatar.Connect("notify::root", avatar.updateFontSize)

	avatar.updateBin(false)
	return avatar
}

// SetSizeRequest sets the avatar size.
func (a *Avatar) SetSizeRequest(size int) {
	a.Bin.SetSizeRequest(size, size)
	a.Image.SetSizeRequest(size, size)
	a.Label.SetSizeRequest(size, size)
	a.updateFontSize()
}

// SizeRequest gets the avatar's size request.
func (a *Avatar) SizeRequest() int {
	w, h := a.Bin.SizeRequest()
	size := w

	if w != h {
		if w > h {
			size = h
		}
		a.Bin.SetSizeRequest(size, size)
	}

	return size
}

// SetFromFile sets the avatar from the given filename.
func (a *Avatar) SetFromFile(file string) {
	a.Image.SetFromFile(file)
	a.updateBin(file != "")
}

// SetFromIconName sets the avatar from the given icon name.
func (a *Avatar) SetFromIconName(iconName string) {
	a.Image.SetFromIconName(iconName)
	a.updateBin(false)
}

// SetFromPixbuf sets the avatar from the given pixbuf.
func (a *Avatar) SetFromPixbuf(p *gdkpixbuf.Pixbuf) {
	a.Image.SetFromPixbuf(p)
	a.updateBin(p != nil)
}

// SetFromPaintable sets the avatar from the given paintable.
func (a *Avatar) SetFromPaintable(p gdk.Paintabler) {
	a.Image.SetFromPaintable(p)
	a.updateBin(p != nil)
}

// SetInitialsTransformer sets the initials transformer function for the Avatar.
// The function will be called to get the initials from the set string.
func (a *Avatar) SetInitialsTransformer(initialsFn func(string) string) {
	a.initialsFunc = initialsFn
}

// ConnectLabel connects the Avatar's initials field to the Label's visible
// text.
func (a *Avatar) ConnectLabel(l *gtk.Label) func() {
	var id glib.SignalHandle
	onMapped := func() {
		f := func() {
			a.SetInitials(l.Text())
		}
		id = l.Connect("notify::label", f)
		f()
	}

	if a.Mapped() {
		onMapped()
	}

	id1 := a.ConnectMap(onMapped)
	id2 := a.ConnectUnmap(func() {
		l.HandlerDisconnect(id)
		id = 0
	})

	return func() {
		a.HandlerDisconnect(id1)
		a.HandlerDisconnect(id2)

		if id > 0 {
			l.HandlerDisconnect(id)
			id = 0
		}
	}
}

// Initials returns the full initials string.
func (a *Avatar) Initials() string {
	return a.Label.Text()
}

// SetInitials sets the string to be displayed as initials.
func (a *Avatar) SetInitials(initials string) {
	a.Label.SetText(a.initialsFunc(initials))

	switch t := a.Image.StorageType(); t {
	case gtk.ImageEmpty:
		a.updateBin(false)
	case gtk.ImageIconName, gtk.ImagePaintable:
		a.updateBin(true)
	default:
		log.Panicln("unknown avatar image type", t)
	}
}

// SetAttributes sets the initial label's Pango attributes.
func (a *Avatar) SetAttributes(attrs *pango.AttrList) {
	a.labelAttrs = attrs
	a.updateFontSize()
}

func (a *Avatar) updateBin(hasImage bool) {
	if a.typeClass != "" {
		a.RemoveCSSClass(a.typeClass)
		a.typeClass = ""
	}

	switch {
	case !hasImage && a.Label.Text() != "":
		a.updateFontSize()
		a.SetChild(a.Label)
		a.typeClass = "adaptive-avatar-label"
	case a.Image.StorageType() == gtk.ImageIconName:
		a.SetChild(a.Image)
		a.typeClass = "adaptive-avatar-icon"
	default:
		a.SetChild(a.Image)
		a.typeClass = "adaptive-avatar-image"
	}

	if a.typeClass != "" {
		a.AddCSSClass(a.typeClass)
	}
}

func (a *Avatar) updateFontSize() {
	// Code adapted from libadwaita's adw-avatar.c.

	// Reset font size first to avoid rounding errors.
	var attrs *pango.AttrList
	if a.labelAttrs != nil {
		attrs = a.labelAttrs.Copy()
	} else {
		attrs = pango.NewAttrList()
	}
	a.Label.SetAttributes(attrs)

	w, h := a.Label.Layout().PixelSize()

	size := float64(a.SizeRequest())
	// This is the size of the biggest square fitting inside the circle.
	squareSize := size / 1.1412
	// The padding has to be a function of the overall size. The 0.4 is how
	// steep the linear function grows and the -5 is just an adjustment for
	// smaller sizes which doesn't have a big impact on bigger sizes. Make also
	// sure we don't have a negative padding.
	padding := math.Max(size*0.4-5, 0)
	maxSize := squareSize - padding
	newFontSize := float64(h) * (maxSize / float64(w))

	attrs.Change(pango.NewAttrSizeAbsolute(clampInt(newFontSize, 0, maxSize) * pango.SCALE))
	a.Label.SetAttributes(attrs)
}

func clampInt(f, min, max float64) int {
	if f < min {
		return int(math.Round(min))
	}
	if f > max {
		return int(math.Round(max))
	}
	return int(math.Round(f))
}
