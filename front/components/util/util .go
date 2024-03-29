package util

import (
	"syscall/js"

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

// remove a class from element without triggering rerender
func RemoveClassById(id string, class string) {
	el := js.Global().Get("document").Call("getElementById", id)
	if !el.IsNull() && !el.IsUndefined() {
		if l := el.Get("classList"); !l.IsNull() && !l.IsUndefined() && l.Length() > 0 {
			l.Call("remove", class)
		}
	}
}

// add a class to element without triggering rerender
func AddClassById(id string, class string) {
	el := js.Global().Get("document").Call("getElementById", id)
	if !el.IsNull() && !el.IsUndefined() {
		if l := el.Get("classList"); !l.IsNull() && !l.IsUndefined() && l.Length() > 0 {
			l.Call("add", class)
		}
	}
}
