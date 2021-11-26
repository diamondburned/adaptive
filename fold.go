package adaptive

import (
	"log"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// FoldRevealButtonIcon is the default icon name for a fold reveal button.
const FoldRevealButtonIcon = "open-menu-symbolic"

// FoldRevealButton is a button that toggles whether or not the fold's sidebar
// should be revealed.
type FoldRevealButton struct {
	*gtk.Revealer
	Button *gtk.ToggleButton
}

// NewFoldRevealButton creates a new fold reveal button. The button is hidden by
// default until a sidebar is connected to it.
func NewFoldRevealButton() *FoldRevealButton {
	button := gtk.NewToggleButton()
	button.SetIconName(FoldRevealButtonIcon)
	button.SetSensitive(false)

	revealer := gtk.NewRevealer()
	revealer.AddCSSClass("adaptive-sidebar-reveal-button")
	revealer.SetTransitionType(gtk.RevealerTransitionTypeCrossfade)
	revealer.SetChild(button)
	revealer.SetRevealChild(false)

	return &FoldRevealButton{
		Revealer: revealer,
		Button:   button,
	}
}

// SetIconName sets the reveal button's icon name.
func (b *FoldRevealButton) SetIconName(icon string) {
	b.Button.SetIconName(icon)
}

// ConnectFold connects the current sidebar reveal button to the given
// sidebar.
func (b *FoldRevealButton) ConnectFold(fold *Fold) {
	b.Button.ConnectClicked(func() {
		fold.SetRevealSide(b.Button.Active())
	})

	fold.NotifyFolded(func(folded bool) {
		b.SetRevealChild(folded)
		b.Button.SetActive(fold.SideIsRevealed())
		b.Button.SetSensitive(folded)
	})
}

// Fold is a component that acts similar to libadwaita's AdwFlap.
type Fold struct {
	gtk.Widgetter
	overlay *gtk.Overlay
	main    *gtk.Box

	siderev    *gtk.Revealer
	sidebox    *Bin
	contentbox *Bin

	onFold func(bool)

	fpos   gtk.PositionType
	fwidth int

	fold   bool
	reveal bool
}

const (
	defaultFoldWidth     = 500
	sidebarMarginTrigger = 50
)

// NewFold creates a new sidebar.
func NewFold(position gtk.PositionType) *Fold {
	f := &Fold{
		fpos:   position,
		fwidth: defaultFoldWidth,
		fold:   false,
	}

	f.sidebox = NewBin()
	f.sidebox.SetSizeRequest(defaultFoldWidth/2, -1)
	f.sidebox.AddCSSClass("adaptive-sidebar-side")
	f.sidebox.SetVExpand(true)

	f.siderev = gtk.NewRevealer()
	f.siderev.AddCSSClass("adaptive-sidebar-revealer")
	f.siderev.SetChild(f.sidebox)
	f.siderev.SetVExpand(true)
	f.siderev.SetHExpand(false)

	f.contentbox = NewBin()
	f.contentbox.AddCSSClass("adaptive-sidebar-child")
	f.contentbox.SetVExpand(true)
	f.contentbox.SetHExpand(true)

	f.main = gtk.NewBox(gtk.OrientationHorizontal, 0)
	f.main.SetVExpand(true)

	switch position {
	case gtk.PosLeft:
		f.siderev.SetHAlign(gtk.AlignStart)
		f.siderev.SetTransitionType(gtk.RevealerTransitionTypeSlideRight)
		f.main.Append(f.siderev)
		f.main.Append(f.contentbox)
	case gtk.PosRight:
		f.siderev.SetHAlign(gtk.AlignEnd)
		f.siderev.SetTransitionType(gtk.RevealerTransitionTypeSlideLeft)
		f.main.Append(f.contentbox)
		f.main.Append(f.siderev)
	default:
		log.Panicln("invalid position given:", position)
	}

	f.overlay = gtk.NewOverlay()
	f.overlay.SetChild(f.main)
	f.overlay.AddCSSClass("adaptive-sidebar")
	f.overlay.SetVExpand(true)

	f.Widgetter = f.overlay
	f.bind(f.Widgetter)

	return f
}

// SetFoldThreshold sets the width threshold that the sidebar will determine
// whether or not to fold.
func (f *Fold) SetFoldThreshold(fwidth int) {
	f.fwidth = fwidth
	f.updateLayout()
}

// FoldThreshold returns the fold width.
func (f *Fold) FoldThreshold() int {
	return f.fwidth
}

// FoldWidth returns the width of the sidebar. It is calculated from the fold
// threshold.
func (f *Fold) FoldWidth() int {
	return f.fwidth / 2
}

// SetSideChild sets the sidebar's side content.
func (f *Fold) SetSideChild(child gtk.Widgetter) {
	f.sidebox.SetChild(child)
}

// SetChild sets the sidebar's main content.
func (f *Fold) SetChild(child gtk.Widgetter) {
	f.contentbox.SetChild(child)
}

// SetFolded sets whether or not the sidebar is folded.
func (f *Fold) SetFolded(folded bool) {
	if folded {
		f.doFold()
	} else {
		f.doUnfold()
	}
}

// SetRevealSide sets whether or not the sidebar is revealed. It does not
// change if the sidebar isn't currently folded.
func (f *Fold) SetRevealSide(reveal bool) {
	f.reveal = reveal
	f.updateRevealSide()
}

// SideIsRevealed returns true if the sidebar is revealed. If the sidebar is not
// folded, then true is returned regardless of what's given into SetRevealSide.
func (f *Fold) SideIsRevealed() bool {
	return f.siderev.RevealChild()
}

// NotifyFolded subscribes f to be called if the sidebar is folded or unfolded.
func (f *Fold) NotifyFolded(fn func(folded bool)) {
	defer f.notifyFolded()

	if f.onFold == nil {
		f.onFold = fn
		return
	}

	old := f.onFold
	f.onFold = func(folded bool) {
		old(folded)
		fn(folded)
	}
}

func (f *Fold) notifyFolded() {
	if f.fold {
		f.overlay.AddCSSClass("adaptive-sidebar-folded")
	} else {
		f.overlay.RemoveCSSClass("adaptive-sidebar-folded")
	}
	if f.onFold != nil {
		f.onFold(f.fold)
	}
}

func (f *Fold) bind(widget gtk.Widgetter) {
	var handle glib.SignalHandle
	var surface *gdk.Surface

	w := gtk.BaseWidget(widget)

	w.ConnectRealize(func() {
		f.updateLayout()

		surface = gdk.BaseSurface(w.GetNative().Surface())
		handle = surface.ConnectLayout(func(int, int) { f.updateLayout() })
	})
	w.ConnectUnrealize(func() {
		surface.HandlerDisconnect(handle)
		surface = nil
	})
}

func (f *Fold) updateLayout() {
	if (f.fwidth + sidebarMarginTrigger) <= f.overlay.AllocatedWidth() {
		f.doUnfold()
	} else {
		f.doFold()
	}
}

func (f *Fold) updateRevealSide() {
	reveal := f.reveal || !f.fold
	f.siderev.SetRevealChild(reveal)

	if reveal {
		f.overlay.AddCSSClass("adaptive-sidebar-open")
	} else {
		f.overlay.RemoveCSSClass("adaptive-sidebar-open")
	}
}

func (f *Fold) doFold() {
	if f.fold {
		return
	}
	f.fold = true

	f.main.Remove(f.siderev)
	f.overlay.AddOverlay(f.siderev)
	f.overlay.SetMeasureOverlay(f.siderev, true)

	f.updateRevealSide()
	f.notifyFolded()
}

func (f *Fold) doUnfold() {
	if !f.fold {
		return
	}
	f.fold = false

	f.overlay.RemoveOverlay(f.siderev)
	switch f.fpos {
	case gtk.PosLeft:
		f.main.Prepend(f.siderev)
	case gtk.PosRight:
		f.main.Append(f.siderev)
	}

	f.updateRevealSide()
	f.notifyFolded()
}
