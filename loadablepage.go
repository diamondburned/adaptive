package adaptive

import (
	"context"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// LoadablePage wraps a child that can be loading at times and error out.
type LoadablePage struct {
	*gtk.Stack
	// ErrorPage is the page that shows the error.
	ErrorPage *StatusPage
	// Spinner is the spinner shown with SetLoading.
	Spinner *gtk.Spinner
	// StopButton is the button used with the Spinner when the user calls
	// SetCancellableLoading.
	StopButton *gtk.Button
	// RetryButton is the button used with the ErrorPage when the user calls
	// SetRetryableError.
	RetryButton *gtk.Button

	content  *Bin
	erroring *gtk.Box
	loading  *gtk.Box

	errLabel func(err error) gtk.Widgetter
	retry    func()
	cancel   context.CancelFunc
}

// NewLoadablePage creates a new LoadablePage widget.
func NewLoadablePage() *LoadablePage {
	p := &LoadablePage{}
	p.Spinner = gtk.NewSpinner()
	p.Spinner.SetSizeRequest(16, 16)
	p.Spinner.Start()

	p.StopButton = gtk.NewButtonWithMnemonic("_Stop")
	p.StopButton.ConnectClicked(func() {
		if p.cancel != nil {
			p.cancel()
			p.cancel = nil
			p.StopButton.SetSensitive(false)
		}
	})

	p.loading = gtk.NewBox(gtk.OrientationVertical, 0)
	p.loading.AddCSSClass("adaptive-loadablepage-loading")
	p.loading.SetVAlign(gtk.AlignCenter)
	p.loading.SetHAlign(gtk.AlignCenter)
	p.loading.Append(p.Spinner)
	p.loading.Append(p.StopButton)

	p.ErrorPage = NewStatusPage()
	p.ErrorPage.AddCSSClass("adaptive-loadablepage-errorpage")
	p.ErrorPage.SetIconName("dialog-error-symbolic")

	p.RetryButton = gtk.NewButtonWithMnemonic("_Retry")
	p.RetryButton.ConnectClicked(func() {
		if p.retry != nil {
			p.retry()
		}
	})

	p.erroring = gtk.NewBox(gtk.OrientationVertical, 0)
	p.erroring.AddCSSClass("adaptive-loadablepage-erroring")
	p.erroring.SetVAlign(gtk.AlignCenter)
	p.erroring.SetHAlign(gtk.AlignCenter)
	p.erroring.Append(p.ErrorPage)
	p.erroring.Append(p.RetryButton)

	p.content = NewBin()
	p.content.AddCSSClass("adaptive-loadablepage-content")

	p.Stack = gtk.NewStack()
	p.Stack.AddCSSClass("adaptive-loadablepage")
	p.Stack.AddChild(p.content)
	p.Stack.AddChild(p.erroring)
	p.Stack.AddChild(p.loading)
	p.Stack.SetVisibleChild(p.loading)
	p.Stack.SetTransitionType(gtk.StackTransitionTypeCrossfade)

	return p
}

// SetError shows an error in the loadable page.
func (p *LoadablePage) SetError(err error) {
	p.SetErrorWidget(NewErrorLabel(err))
}

// SetErrorWidget is like SetError, except the given widget is set instead of a
// new ErrorLabel.
func (p *LoadablePage) SetErrorWidget(w gtk.Widgetter) {
	p.ensureFresh()
	p.content.SetChild(nil)

	p.ErrorPage.SetDescription(w)
	p.RetryButton.SetVisible(p.retry != nil)

	p.SetVisibleChild(p.erroring)
}

// SetRetryFunc sets the function to be called when the user clicks the retry
// button in the error page. If the function is nil, then LoadablePage cannot be
// retried and the button is not shown.
func (p *LoadablePage) SetRetryFunc(fn func()) {
	p.retry = fn
	p.RetryButton.SetVisible(p.retry != nil)
}

// SetLoading shows a loading animation in the loadable page.
func (p *LoadablePage) SetLoading() {
	p.ensureFresh()
	p.content.SetChild(nil)

	p.Spinner.Start()
	p.SetVisibleChild(p.loading)
	p.StopButton.SetVisible(false)
}

// SetCancellableLoading shows a loading animation with a button in the busy
// box. If the button is clicked, then it's no longer clickable, and the
// returned context is interrupted. It does nothing more, and the user must call
// either SetChild or SetError to continue.
func (p *LoadablePage) SetCancellableLoading() context.Context {
	p.SetLoading()

	p.StopButton.SetVisible(true)
	p.StopButton.SetSensitive(true)

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	return ctx
}

// SetChild sets the main child of the loadable page.
func (p *LoadablePage) SetChild(child gtk.Widgetter) {
	p.ensureFresh()

	p.content.SetChild(child)
	p.SetVisibleChild(p.content)
}

func (p *LoadablePage) ensureFresh() {
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	p.Spinner.Stop()
	p.StopButton.SetSensitive(false)
}
