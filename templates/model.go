package templates

import "github.com/a-h/templ"

type Nav struct {
	active string
	pages  map[string]templ.SafeURL
}

func (n *Nav) Render() templ.Component {
	return Navbar(*n)
}

func NewNav(pg map[string]templ.SafeURL, act string) Nav {
	return Nav{
		active: act,
		pages:  pg,
	}
}

type Recipe struct {
	name     string
	preptime float32
	steps    []string
}
