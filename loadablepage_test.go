package adaptive_test

import (
	"errors"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleLoadablePage() {
	const margin = 8

	testapp.Run("loadable-page", func(app *gtk.Application) {
		errorCheck := gtk.NewCheckButtonWithLabel("Erroneous")
		errorCheck.SetHExpand(true)

		loadButton := gtk.NewButtonWithLabel("Load")
		loadButton.SetHExpand(true)

		child := gtk.NewBox(gtk.OrientationVertical, margin)
		child.SetMarginTop(margin)
		child.SetMarginBottom(margin)
		child.SetMarginStart(margin)
		child.SetMarginEnd(margin)
		child.Append(errorCheck)
		child.Append(loadButton)

		main := adaptive.NewLoadablePage()
		main.SetChild(child)

		loadButton.ConnectClicked(func() {
			erroneous := errorCheck.Active()
			main.SetLoading()

			glib.TimeoutSecondsAdd(5, func() {
				if erroneous {
					main.SetError(errors.New("failed to load busy box: checkmark was active"))
				} else {
					main.SetChild(child)
				}
			})
		})

		w := testapp.NewWindow(app, "Loadable", 250, -1)
		w.SetChild(main)
		w.Show()
	})
	// Output:
}
