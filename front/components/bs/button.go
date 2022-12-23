package bs

import (
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
)

func Button(m ...v.MarkupOrChild) *v.HTML {
	m = append(m, v.Markup(v.Class("btn")))
	return e.Button(m...)
}

// ButtonPrimary
func ButtonPry(m ...v.MarkupOrChild) *v.HTML {
	m = append(m, v.Markup(v.Class("btn-primary")))
	return Button(m...)
}

func Button2ry(m ...v.MarkupOrChild) *v.HTML {
	m = append(m, v.Markup(v.Class("btn-secondary")))
	return Button(m...)
}
