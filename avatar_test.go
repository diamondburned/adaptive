package adaptive_test

import (
	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleAvatar() {
	testapp.Run("avatar", func(app *gtk.Application) {
		sizes := []int{16, 24, 32, 48, 56, 64}

		box := gtk.NewBox(gtk.OrientationHorizontal, 8)
		box.SetHExpand(true)
		box.SetVExpand(true)
		box.SetHAlign(gtk.AlignCenter)
		box.SetVAlign(gtk.AlignCenter)
		box.SetMarginStart(8)
		box.SetMarginEnd(8)
		box.SetMarginTop(8)
		box.SetMarginBottom(8)

		for _, size := range sizes {
			avy := adaptive.NewAvatar(size)
			avy.SetInitials("Ferris Argyle")
			box.Append(avy)
		}

		w := testapp.NewWindow(app, "Avatars", -1, -1)
		w.SetChild(box)
		w.Show()
	})
	// Output:
}
