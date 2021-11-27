// Package adaptive provides adaptive GTK4 widget components. It's mainly for
// use in applications that aim to support both mobile and desktop viewports.
// It is an alternative to libadwaita.
package adaptive

import (
	_ "embed"
	"log"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var styleCSS string

// Init initializes package adaptive. The caller should call this on application
// activation.
func Init() {
	css := gtk.NewCSSProvider()
	css.Connect("parsing-error", func(section *gtk.CSSSection, err error) {
		loc := section.StartLocation()
		log.Printf(
			"adaptive: bug: error parsing CSS at style.css:%d:%d",
			loc.Lines(), loc.LineChars(),
		)
	})
	css.LoadFromData(styleCSS)

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(), css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
