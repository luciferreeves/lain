package tags

import (
	"log"

	"github.com/flosch/pongo2/v6"
)

type templateTag struct {
	Name string
	Fn   pongo2.TagParser
}

func Initialize() {

	tags := []templateTag{
		{"url", url},
		// {"static", static},
	}

	for _, t := range tags {
		if err := pongo2.RegisterTag(t.Name, t.Fn); err != nil {
			log.Println("Failed to register tag:", t.Name, "Error:", err)
		}
	}
}
