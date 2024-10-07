package templates

import (
	"github.com/a-h/templ"
)

// a page included in the navbar
type Nav struct {
	// current page the user is on
	active string
	// link to the page that the nav element is linked to
	pages map[string]templ.SafeURL
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
