package adaptive_test

import (
	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/adaptive/internal/testdata"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleAvatar() {
	testapp.Run("avatar", func(app *gtk.Application) {
		adaptive.Init()

		sizes := []int{16, 24, 32, 48, 56, 64}

		main := gtk.NewBox(gtk.OrientationVertical, 8)
		main.SetMarginStart(8)
		main.SetMarginEnd(8)
		main.SetMarginTop(8)
		main.SetMarginBottom(8)

		avatarFns := []func(*adaptive.Avatar){
			func(a *adaptive.Avatar) { a.SetInitials("Ferris Argyle") },
			func(a *adaptive.Avatar) { a.SetFromPixbuf(testdata.MustAvatarPixbuf()) },
		}

		for _, avatarFn := range avatarFns {
			box := gtk.NewBox(gtk.OrientationHorizontal, 8)
			box.SetHExpand(true)
			box.SetVExpand(true)
			box.SetHAlign(gtk.AlignCenter)
			box.SetVAlign(gtk.AlignCenter)

			for _, size := range sizes {
				avy := adaptive.NewAvatar(size)
				avatarFn(avy)
				box.Append(avy)
			}

			main.Append(box)
		}

		w := testapp.NewWindow(app, "Avatars", -1, -1)
		w.SetChild(main)
		w.Show()
	})
	// Output:
}
