package adaptive

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Bin is a widget that holds a single widget.
type Bin struct {
	*gtk.Box
	child *gtk.Widget
}

// NewBin creates a new bin.
func NewBin() *Bin {
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.SetLayoutManager(gtk.NewBinLayout())
	return &Bin{box, nil}
}

// SetChild sets the child in the bin. If child is nil, then the box is cleared.
func (b *Bin) SetChild(child gtk.Widgetter) {
	if child == b.child {
		return
	}

	if b.child != nil {
		b.child.Unparent()
		b.child = nil
	}

	if child != nil {
		b.child = gtk.BaseWidget(child)
		b.child.SetParent(b)
	}
}

// Child returns the Bin's child.
func (b *Bin) Child() gtk.Widgetter {
	return b.child
}

// IsChild returns true if the given child is the bin's child.
func (b *Bin) IsChild(child gtk.Widgetter) bool {
	return glib.ObjectEq(b.child, child)
}
