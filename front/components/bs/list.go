package bs

import (
	u "github.com/ferealqq/wienerlist/front/components/util"
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
)

// Bootstrap list group
func List(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, u.Classes("list-group"))
	return e.Div(markup...)
}

func ListItem(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, u.Classes("list-group-item"))
	return e.Div(markup...)
}
