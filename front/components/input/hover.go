package input

import (
	"github.com/ferealqq/wienerlist/front/components/bs"
	v "github.com/hexops/vecty"
	// e "github.com/hexops/vecty/elem"
)

func HoverTextInput(m ...v.MarkupOrChild) *v.HTML {
	m = append(m, v.Markup(
		v.Class("border", "border-hover"),
	))
	return bs.Input(m...)
}

func HoverTextAreaInput(m ...v.MarkupOrChild) *v.HTML {
	m = append(m, v.Markup(
		v.Class("border", "border-hover"),
	))
	return bs.TextArea(m...)
}
