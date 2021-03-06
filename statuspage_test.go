package adaptive_test

import (
	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleStatusPage() {
	testapp.Run("status-page", func(app *gtk.Application) {
		adaptive.Init()

		status := adaptive.NewStatusPage()
		status.SetTitle("Uh oh!")
		status.SetIconName("computer-fail-symbolic")
		status.SetDescriptionText("An oopsie-whoopsie has occured. Please throw your computer out the window.")

		w := testapp.NewWindow(app, "Status Page", 350, 200)
		w.SetChild(status)
		w.Show()
	})
	// Output:
}
