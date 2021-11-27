package adaptive_test

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/adaptive/internal/testapp"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleFold() {
	testapp.Run("fold", func(app *gtk.Application) {
		stack := newStack()
		stack.SetHExpand(true)

		stackside := gtk.NewStackSidebar()
		stackside.SetStack(stack)

		fold := adaptive.NewFold(gtk.PosLeft)
		fold.SetSideChild(stackside)
		fold.SetChild(stack)

		foldButton := adaptive.NewFoldRevealButton()
		foldButton.ConnectFold(fold)

		h := gtk.NewHeaderBar()
		h.PackStart(foldButton)

		w := testapp.NewWindow(app, "Example Sidebar", 450, 300)
		w.SetChild(fold)
		w.SetTitlebar(h)
		w.Show()
	})
	// Output:
}

func newStack() *gtk.Stack {
	stack := gtk.NewStack()
	stack.SetTransitionType(gtk.StackTransitionTypeSlideUpDown)

	for i := 0; i < 5; i++ {
		istr := strconv.Itoa(i)
		content := gtk.NewLabel(fmt.Sprintf("You're in stack number %s.", istr))
		stack.AddTitled(content, "stack-"+istr, "Stack "+istr)
	}

	return stack
}
