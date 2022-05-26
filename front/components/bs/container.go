package bs

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Container(markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup, vecty.Markup(vecty.Class("container")))
	return elem.Div(markup...)
}

func ContainerFluid(markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup, vecty.Markup(vecty.Class("container-fluid")))
	return elem.Div(markup...)
}
