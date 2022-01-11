package adaptive_test

import (
	"context"
	"errors"
	"time"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleLoadablePage() {
	const margin = 8

	testapp.Run("loadable-page", func(app *gtk.Application) {
		adaptive.Init()

		errorCheck := gtk.NewCheckButtonWithLabel("Erroneous")
		errorCheck.SetHExpand(true)

		loadButton := gtk.NewButtonWithLabel("Load")
		loadButton.SetHExpand(true)

		child := gtk.NewBox(gtk.OrientationVertical, margin)
		child.SetMarginTop(margin)
		child.SetMarginBottom(margin)
		child.SetMarginStart(margin)
		child.SetMarginEnd(margin)
		child.SetHAlign(gtk.AlignCenter)
		child.SetVAlign(gtk.AlignCenter)
		child.SetSizeRequest(250, -1)
		child.Append(errorCheck)
		child.Append(loadButton)

		main := adaptive.NewLoadablePage()
		main.SetHExpand(true)
		main.SetChild(child)
		main.SetRetryFunc(func() { main.SetChild(child) })

		loadButton.ConnectClicked(func() {
			erroneous := errorCheck.Active()

			ctx := main.SetCancellableLoading(context.Background())

			go func() {
				select {
				case <-ctx.Done():
					glib.IdleAdd(func() { main.SetError(ctx.Err()) })
				case <-time.After(2 * time.Second):
					glib.IdleAdd(func() {
						if erroneous {
							main.SetError(errors.New("failed to load busy box: checkmark was active"))
						} else {
							main.SetChild(child)
						}
					})
				}
			}()
		})

		w := testapp.NewWindow(app, "Loadable", 250, -1)
		w.SetChild(main)
		w.SetDefaultSize(300, 450)
		w.Show()
	})
	// Output:
}
