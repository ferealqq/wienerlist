package bs

import (
	u "github.com/ferealqq/wienerlist/front/components/util"
	v "github.com/hexops/vecty"
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
