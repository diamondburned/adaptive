package testapp

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"testing"
	"unicode"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func init() {
	gtk.Init()
}

// NewWindow creates a new window.
func NewWindow(title string, w, h int) *gtk.Window {
	window := gtk.NewWindow()
	window.SetTitle(title)
	window.SetDefaultSize(w, h)
	return window
}

// Run runs an app with the given window.
func Run(w *gtk.Window) {
	runApp(newApp(w), func(err error) {
		log.Fatalln(err)
	})
}

// RunTest runs a test with the given window.
func RunTest(t *testing.T, w *gtk.Window) {
	runApp(newApp(w), func(err error) {
		t.Fatal(err)
	})
}

func newApp(w *gtk.Window) *gtk.Application {
	id := "com.github.diamondburned.adaptive.tests." + slugify(w.Title())

	app := gtk.NewApplication(id, 0)
	app.ConnectActivate(func() {
		w.SetApplication(app)
		w.Show()

		if !testing.Verbose() {
			glib.TimeoutSecondsAdd(1, func() { w.Close() })
		}
	})

	return app
}

func runApp(app *gtk.Application, errFn func(error)) {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		_, ok := <-ch
		if ok {
			glib.IdleAdd(func() { app.Quit() })
		}
	}()

	code := app.Run(nil)

	signal.Stop(ch)
	close(ch)

	if code != 0 {
		errFn(fmt.Errorf("application exited with status %d", code))
	}
}

func slugify(name string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLower(r), unicode.IsDigit(r):
			return r
		case unicode.IsUpper(r):
			return unicode.ToLower(r)
		default:
			return '_'
		}
	}, name)
}
