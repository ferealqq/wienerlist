package bs

import (
	u "github.com/ferealqq/wienerlist/front/components/util"
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
)

func Active(b bool) v.MarkupList {
	if b {
		return u.Classes("active")
	}
	return v.Markup()
}

func Icon(name string, ms ...v.MarkupOrChild) *v.HTML {
	ms = append(ms, u.Classes("bi", "bi-"+name))
	return v.Tag("i", ms...)
}

func Label(forId string, text string, m ...v.MarkupOrChild) *v.HTML {
	m = append(m, v.Markup(v.Attribute("for", forId)))
	m = append(m, v.Text(text))
	return e.Label(m...)
}
