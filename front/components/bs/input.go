package bs

import (
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
)

func Input(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("form-control")))
	return e.Input(markup...)
}

func TextArea(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("form-control")))
	return e.TextArea(markup...)
}
