package main

import (
	"context"
	"os"

	"github.com/James-Trauger/Recipouir/templates"
	"github.com/a-h/templ"
)

func main() {
	var routes map[string]templ.SafeURL = map[string]templ.SafeURL{
		"home":       templ.SafeURL("/home"),
		"my recipes": templ.SafeURL("/my-recipes"),
	}
	comp := templates.Navbar(routes, "my recipes")
	comp.Render(context.Background(), os.Stdout)
}
