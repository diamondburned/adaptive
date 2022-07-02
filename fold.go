package adaptive

import (
	"log"
	"math"

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

	fold.NotifyRevealed(func(revealed bool) {
		b.Button.SetActive(revealed)
	})

	fold.NotifyFolded(func(folded bool) {
		b.SetRevealChild(folded)
		b.Button.SetActive(fold.SideIsRevealed())
		b.Button.SetSensitive(folded)
	})
}

// BindFolds binds the given folds to have synchronized fold and reveal states.
// The first fold is used as the basis for the width.
func BindFolds(folds ...*Fold) {
	for _, fold := range folds[1:] {
		fold.SetShouldFoldFunc(func() bool { return folds[0].fold })
	}

	folds[0].NotifyFolded(func(bool) {
		for _, fold := range folds[1:] {
			fold.updateLayout()
		}
	})

	var mutex bool
	do := func(f func()) {
		if !mutex {
			mutex = true
			f()
			mutex = false
		}
	}

	// We need to do this because we want to unreveal everything if even one
	// fold gets collapsed.
	for i := range folds {
		i := i
		folds[i].NotifyRevealed(func(revealed bool) {
			do(func() {
				for j, fold := range folds {
					if i == j {
						continue
					}
					fold.reveal = revealed
					fold.doRevealSide()
				}
			})
		})
	}
}

// Fold is a component that acts similar to libadwaita's AdwFlap.
type Fold struct {
	*gtk.Widget
	overlay *gtk.Overlay
	main    *gtk.Box

	dimming    *gtk.Box
	siderev    *gtk.Revealer
	sidebox    *Bin
	contentbox *gtk.Overlay

	onFold     func(bool)
	shouldFold func() bool

	fpos   gtk.PositionType
	fthres int
	fwidth int

	fold   bool
	reveal bool
}

// Fold threshold constants that determine when swiping velocities should be
// handled.
var (
	FoldXThreshold = [2]float64{800, math.Inf(+1)}
	FoldYThreshold = [2]float64{0, 4000}
)

func isInThreshold(f float64, thres [2]float64) bool {
	return false ||
		+thres[0] <= f && f <= +thres[1] ||
		-thres[0] >= f && f >= -thres[1]
}

const (
	defaultFoldWidth     = 200
	defaultFoldThreshold = 400
)

// NewFold creates a new sidebar.
func NewFold(position gtk.PositionType) *Fold {
	f := &Fold{
		fpos:   position,
		fthres: defaultFoldThreshold,
		fwidth: defaultFoldWidth,
		fold:   false,
	}

	f.sidebox = NewBin()
	f.sidebox.SetSizeRequest(f.fwidth, -1)
	f.sidebox.AddCSSClass("adaptive-sidebar-side")
	f.sidebox.SetVExpand(true)

	f.siderev = gtk.NewRevealer()
	f.siderev.AddCSSClass("adaptive-sidebar-revealer")
	f.siderev.SetChild(f.sidebox)
	f.siderev.SetVExpand(true)
	f.siderev.SetHExpand(false)

	f.dimming = gtk.NewBox(gtk.OrientationHorizontal, 0)
	f.dimming.AddCSSClass("adaptive-sidebar-dimming")
	f.dimming.SetCanTarget(false)
	f.dimming.SetCanFocus(false)
	f.dimming.SetVExpand(true)
	f.dimming.SetHExpand(true)
	f.dimming.SetVisible(false)

	f.contentbox = gtk.NewOverlay()
	f.contentbox.AddOverlay(f.dimming)
	f.contentbox.SetClipOverlay(f.dimming, true)
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

	f.Widget = gtk.BaseWidget(f.overlay)
	f.shouldFold = func() bool { return f.overlay.AllocatedWidth() < f.fthres }
	f.bind()
	f.updateLayout()

	// Bind handlers that will blur the content box if the revealer is over it.
	f.NotifyFolded(func(folded bool) { f.updateState() })
	f.siderev.Connect("notify::reveal-child", func() { f.updateState() })

	// Controller for clicking on the background.
	bgclicker := gtk.NewGestureClick()
	bgclicker.SetExclusive(true)
	bgclicker.ConnectPressed(func(n int, x, y float64) {
		if f.fold && f.siderev.RevealChild() {
			f.siderev.SetRevealChild(false)
		}
	})
	// Bind it to the main widget. Note that f.main will be underneath the
	// revealer overlay, so we can assume that if it's clicked, it's ever going
	// to be clicked outside the revealer.
	f.main.AddController(bgclicker)

	// Controller for swiping.
	swiper := gtk.NewGestureSwipe()
	swiper.SetExclusive(true)
	swiper.SetTouchOnly(true)
	swiper.ConnectSwipe(func(velX, velY float64) {
		if !f.fold {
			return
		}
		if isInThreshold(velX, FoldXThreshold) && isInThreshold(velY, FoldYThreshold) {
			// Determine the orientation of the swiping by inspecting the sign
			// of the X (horizontal) velocity.
			// Negative is right-to-left, and positive is left-to-right.
			f.siderev.SetRevealChild(velX > 0)
		}
	})
	f.overlay.AddController(swiper)

	return f
}

// SetWidthFunc sets the function to get the width to determine the fold
// threshold.
func (f *Fold) SetWidthFunc(widthFunc func() int) {
	f.shouldFold = func() bool { return widthFunc() < f.fthres }
}

// SetShouldFoldFunc sets the callback to determine whether or not fold should
// be folded. It overrides SetWidthFunc.
func (f *Fold) SetShouldFoldFunc(shouldFold func() bool) {
	f.shouldFold = shouldFold
}

// SetFoldThreshold sets the width threshold that the sidebar will determine
// whether or not to fold.
func (f *Fold) SetFoldThreshold(threshold int) {
	f.fthres = threshold
	f.updateLayout()
}

// FoldThreshold returns the fold width.
func (f *Fold) FoldThreshold() int {
	return f.fthres
}

// SetFoldWidth sets the width of the sidebar. The width must be lower than the
// fold threshold.
func (f *Fold) SetFoldWidth(width int) {
	f.sidebox.SetSizeRequest(width, -1)
	f.updateLayout()
}

// FoldWidth returns the width of the sidebar. It is calculated from the fold
// threshold.
func (f *Fold) FoldWidth() int {
	w, _ := f.sidebox.SizeRequest()
	return w
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
	f.setFold(folded)
}

// SetRevealSide sets whether or not the sidebar is revealed. It does not
// change if the sidebar isn't currently folded.
func (f *Fold) SetRevealSide(reveal bool) {
	f.reveal = reveal
	f.doRevealSide()
}

func (f *Fold) doRevealSide() {
	reveal := f.reveal || !f.fold
	f.siderev.SetRevealChild(reveal)
}

// SideIsRevealed returns true if the sidebar is revealed. If the sidebar is not
// folded, then true is returned regardless of what's given into SetRevealSide.
func (f *Fold) SideIsRevealed() bool {
	return f.siderev.RevealChild()
}

// NotifyRevealed subscribes fn to be called if the sidebar is revealed or not.
func (f *Fold) NotifyRevealed(fn func(revealed bool)) {
	f.siderev.Connect("notify::reveal-child", func() {
		fn(f.siderev.RevealChild())
	})
}

// NotifyFolded subscribes fn to be called if the sidebar is folded or unfolded.
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

// QueueResize should be called when Fold's parent widths are changed.
func (f *Fold) QueueResize() {
	f.updateLayout()
	gtk.BaseWidget(f).QueueResize()
}

func (f *Fold) bind() {
	var handle glib.SignalHandle
	var surface *gdk.Surface

	w := f.overlay

	// Hack to resize the first time the widget has a size.
	w.AddTickCallback(func(gtk.Widgetter, gdk.FrameClocker) bool {
		if f.AllocatedWidth() > 0 {
			f.updateLayout()
			return false
		}
		// Retry on the next frame.
		return true
	})

	w.ConnectRealize(func() {
		// TODO: this doesn't cover the page where the inside is changed without
		// the window being resized. It might be worth it to have a slow path
		// that checks the width and updates the size every 1000/30ms or so.
		surface = gdk.BaseSurface(w.Native().Surface())
		handle = surface.Connect("notify::width", func() { f.updateLayout() })
	})

	w.ConnectUnrealize(func() {
		surface.HandlerDisconnect(handle)
		surface = nil
	})
}

func (f *Fold) updateLayout() {
	f.setFold(f.shouldFold())
}

func (f *Fold) updateState() {
	reveal := f.siderev.RevealChild()

	// If we're folded, then the user shouldn't be able to target the
	// content box behind the revealer.
	f.contentbox.SetCanTarget(!f.fold || !reveal)
	// Only show the dimming overlay if we're folded.
	f.dimming.SetVisible(f.fold)

	if reveal {
		f.overlay.AddCSSClass("adaptive-sidebar-open")
	} else {
		f.overlay.RemoveCSSClass("adaptive-sidebar-open")
	}
}

func (f *Fold) setFold(fold bool) {
	if f.fold == fold {
		return
	}
	f.fold = fold

	if fold {
		f.main.Remove(f.siderev)
		f.overlay.AddOverlay(f.siderev)
		f.overlay.SetMeasureOverlay(f.siderev, true)

		f.doRevealSide()
		f.notifyFolded()
	} else {
		f.overlay.RemoveOverlay(f.siderev)
		switch f.fpos {
		case gtk.PosLeft:
			f.main.Prepend(f.siderev)
		case gtk.PosRight:
			f.main.Append(f.siderev)
		}

		f.doRevealSide()
		f.notifyFolded()
	}
}
