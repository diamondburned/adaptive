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

// NewWindow creates a new window.
func NewWindow(app *gtk.Application, title string, w, h int) *gtk.ApplicationWindow {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle(title)
	window.SetDefaultSize(w, h)
	return window
}

// Run runs an app with the given window.
func Run(name string, f func(*gtk.Application)) {
	runApp(name, f, func(err error) {
		log.Fatalln(err)
	})
}

// RunTest runs a test with the given window.
func RunTest(t *testing.T, f func(*gtk.Application)) {
	runApp(slugify(t.Name()), f, func(err error) {
		t.Fatal(err)
	})
}

func runApp(name string, f func(*gtk.Application), errFn func(error)) {
	id := "com.github.diamondburned.adaptive.tests." + name
	app := gtk.NewApplication(id, 0)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		_, ok := <-ch
		if ok {
			glib.IdleAdd(func() { app.Quit() })
		}
	}()

	app.ConnectActivate(func() {
		f(app)
	})

	code := app.Run([]string{os.Args[0]})

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
