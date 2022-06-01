package util

import (
	"errors"
	"strings"
	"syscall/js"
	"unicode/utf8"
)

func search() string {
	return js.Global().Get("location").Get("search").String()
}

func GetSearch() string {
	return search()
}

var ErrSearchEmpty = errors.New("location.search didn't contain question mark, nothing to parse")
var ErrInvalidSearchParam = errors.New("location.search is malformed")

func GetSearchParams() (map[string]string, error) {
	if s := search(); strings.Contains(s, "?") {
		// remove the first character from the string
		s := trimFirstRune(s)
		r := make(map[string]string)
		for _, v := range strings.Split(s, "&") {
			if a := strings.Split(v, "="); len(a) != 2 {
				return nil, ErrInvalidSearchParam
			} else {
				r[a[0]] = a[1]
			}
		}
		return r, nil
	} else {
		return nil, ErrSearchEmpty
	}
}

var ErrSearchParamNotFound = errors.New("search param not found")

func GetSearchParam(key string) (string, error) {
	if p, e := GetSearchParams(); e != nil {
		return "", e
	} else {
		if val, ok := p[key]; ok {
			return val, nil
		} else {
			return "", ErrSearchParamNotFound
		}
	}
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// push state event name
var eventName = "pushstate"

func CreateSearchParamListener(f func()) {
	// custom js event
	e := js.Global().Get("document").Call("createEvent", "Event")
	// initialize event
	e.Call("initEvent", eventName, true, true)

	onPushState := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f()
		return nil
	})

	js.Global().Get("document").Call("addEventListener", eventName, onPushState)
	// js window.history
	h := js.Global().Get("window").Get("history")
	og := h.Get("pushState")
	pushState := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// for some reason this can't be invoked with the value of object?
		return og.Invoke(h.Type(), args)
	})
	h.Set("pushState", pushState)
}
