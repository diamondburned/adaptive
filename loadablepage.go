package adaptive

import (
	"context"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// LoadablePage wraps a child that can be loading at times and error out.
type LoadablePage struct {
	*gtk.Stack
	Spinner   *gtk.Spinner
	ErrorPage *StatusPage
	content   *Bin

	ctx context.Context
}

// NewLoadablePage creates a new LoadablePage widget.
func NewLoadablePage() *LoadablePage {
	p := &LoadablePage{}
	p.Spinner = gtk.NewSpinner()
	p.Spinner.SetVExpand(true)
	p.Spinner.SetHExpand(true)
	p.Spinner.SetVAlign(gtk.AlignCenter)
	p.Spinner.SetHAlign(gtk.AlignCenter)
	p.Spinner.SetSizeRequest(24, 24)
	p.Spinner.Start()

	p.ErrorPage = NewStatusPage()
	p.ErrorPage.AddCSSClass("adaptive-busybox-errorpage")
	p.ErrorPage.SetIconName("dialog-error-symbolic")

	p.content = NewBin()

	p.Stack = gtk.NewStack()
	p.Stack.AddChild(p.Spinner)
	p.Stack.AddChild(p.ErrorPage)
	p.Stack.AddChild(p.content)
	p.Stack.SetVisibleChild(p.Spinner)
	p.Stack.SetTransitionType(gtk.StackTransitionTypeCrossfade)

	ctx, cancel := context.WithCancel(context.Background())
	p.Stack.ConnectDestroy(cancel)
	p.ctx = ctx

	return p
}

// SetError shows an error in the busy box.
func (p *LoadablePage) SetError(err error) {
	p.content.SetChild(nil)
	p.Spinner.Stop()

	errLabel := NewErrorLabel(err)
	errLabel.Reveal.SetTransitionType(gtk.RevealerTransitionTypeCrossfade)
	p.ErrorPage.SetDescription(errLabel)

	p.SetVisibleChild(p.ErrorPage)
}

// SetLoading shows a loading animation in the busy box.
func (p *LoadablePage) SetLoading() {
	p.content.SetChild(nil)
	p.Spinner.Start()
	p.SetVisibleChild(p.Spinner)
}

// SetChild sets the main child of the busy box.
func (p *LoadablePage) SetChild(child gtk.Widgetter) {
	p.content.SetChild(nil)
	p.Spinner.Stop()
	p.content.Append(child)
	p.SetVisibleChild(p.content)
}
