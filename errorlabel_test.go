package adaptive_test

import (
	"errors"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleErrorLabel() {
	testapp.Run("error-label", func(app *gtk.Application) {
		err := errors.New("failed to open hello.txt: filesystem error: missing hard drive")

		status := adaptive.NewErrorLabel(err)
		status.SetMarginTop(8)
		status.SetMarginBottom(8)
		status.SetMarginStart(8)
		status.SetMarginEnd(8)

		w := testapp.NewWindow(app, "Error", 150, 250)
		w.SetChild(status)
		w.Show()
	})
	// Output:
}
