package adaptive

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// Bin is a widget that holds a single widget.
type Bin struct {
	*gtk.Box
	child gtk.Widgetter
}

// NewBin creates a new bin.
func NewBin() *Bin {
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.SetHomogeneous(true)
	box.SetHExpand(true)
	box.SetVExpand(true)
	return &Bin{box, nil}
}

// SetChild sets the child in the bin. If child is nil, then the box is cleared.
func (b *Bin) SetChild(child gtk.Widgetter) {
	if child == b.child {
		return
	}
	if b.child != nil {
		b.Remove(b.child)
	}
	b.child = child
	if child != nil {
		b.Append(child)
	}
}

// Child returns the Bin's child.
func (b *Bin) Child() gtk.Widgetter {
	return b.child
}

// IsChild returns true if the given child is the bin's child.
func (b *Bin) IsChild(child gtk.Widgetter) bool {
	return b.child == child
}
