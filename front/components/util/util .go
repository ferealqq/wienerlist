package util

import (
	v "github.com/hexops/vecty"
)

func Classes(s ...string) v.MarkupList {
	var l []v.Applyer
	for _, i := range s {
		l = append(l, v.Class(i))
	}

	return v.Markup(l...)
}

// Attributes for markup
func Atrs(s map[string]interface{}) v.MarkupList {
	var l []v.Applyer
	for k, val := range s {
		l = append(l, v.Attribute(k, val))
	}

	return v.Markup(l...)
}

func Atr(k string, val interface{}) v.MarkupList {
	return v.Markup(v.Attribute(k, val))
}
