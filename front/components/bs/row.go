package bs

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Row(markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup, vecty.Markup(vecty.Class("row")))
	return elem.Div(markup...)
}
