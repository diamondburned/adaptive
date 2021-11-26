package adaptive

import (
	"log"
	"math"
	"strings"
	"unicode"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

// TransformInitials transforms the given input string into initials.
func TransformInitials(in string) string {
	if in == "" {
		return ""
	}

	r := []rune(in)

	b := strings.Builder{}
	b.Grow(len(in))
	b.WriteRune(r[0])

	for _, r := range r[1:] {
		if unicode.IsUpper(r) {
			b.WriteRune(r)
			break
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

	avatar.SetHAlign(gtk.AlignCenter)
	avatar.SetVAlign(gtk.AlignCenter)
	avatar.SetSizeRequest(size)
	avatar.SetChild(avatar.Image)
	avatar.AddCSSClass("adaptive-avatar")
	avatar.Connect("notify::root", avatar.updateFontSize)

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

// SetFromPixbuf sets the avatar from the given pixbuf.
func (a *Avatar) SetFromPixbuf(p *gdkpixbuf.Pixbuf) {
	a.Image.SetFromPixbuf(p)
	a.updateBin(p != nil)
}

// SetFromPaintable sets the avatar from the given paintable.
func (a *Avatar) SetFromPaintable(p *gdk.Paintable) {
	a.Image.SetFromPaintable(p)
	a.updateBin(p != nil)
}

// SetInitialsTransformer sets the initials transformer function for the Avatar.
// The function will be called to get the initials from the set string.
func (a *Avatar) SetInitialsTransformer(initialsFn func(string) string) {
	a.initialsFunc = initialsFn
}

// Initials returns the full initials string.
func (i *Avatar) Initials() string {
	return i.Label.Text()
}

// SetInitials sets the string to be displayed as initials.
func (a *Avatar) SetInitials(initials string) {
	a.Label.SetText(a.initialsFunc(initials))
	a.Label.SetTooltipText(initials)

	var hasImage bool
	switch t := a.Image.StorageType(); t {
	case gtk.ImageEmpty, gtk.ImageIconName:
		hasImage = false
	case gtk.ImagePaintable:
		hasImage = true
	default:
		log.Panicln("unknown avatar image type", t)
	}

	a.updateBin(hasImage)
}

// SetAttributes sets the initial label's Pango attributes.
func (a *Avatar) SetAttributes(attrs *pango.AttrList) {
	a.labelAttrs = attrs
	a.updateFontSize()
}

func (a *Avatar) updateBin(hasImage bool) {
	if hasImage || a.Label.Text() == "" {
		a.SetChild(a.Image)
		return
	}

	a.updateFontSize()
	a.SetChild(a.Label)
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
