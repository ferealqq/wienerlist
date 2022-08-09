package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"
	"unicode/utf8"

	"github.com/hexops/vecty"
)

func pathname() string {
	return js.Global().Get("location").Get("pathname").String()
}

// FIXME delete pathname function rename it to this
func LocationPathname() string {
	return pathname()
}

func search() string {
	return js.Global().Get("location").Get("search").String()
}

func GetSearch() string {
	return search()
}

// search param regex
var searchRe = regexp.MustCompile(`(?m).+\/?\w{1,100}=(\w{1,1000})$`)
var ErrSearchEmpty = errors.New("location.search didn't contain question mark, nothing to parse")
var ErrInvalidSearchParam = errors.New("location.search is malformed")

func GetSearchParams() (map[string]string, error) {
	if s := search(); searchRe.MatchString(s) {
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

func RerenderRoute(route string) {
	for _, r := range routes {
		// this would probably be better if it would use strict match
		if r.pattern.MatchString(route) {
			vecty.Rerender(r)
			return
		}
	}
}

// back to the previous url
func Back() {
	println("before " + pathname())
	js.Global().Get("history").Call(
		"go",
		"-1",
	)
	RerenderRoute(pathname())
}

func Redirect(route string) {
	js.Global().Get("history").Call(
		"pushState",
		map[string]interface{}{"redirectRoute": route},
		route,
		route,
	)
	RerenderRoute(route)
}

type defaultNotFound struct {
	vecty.Core
}

func (d *defaultNotFound) Render() vecty.ComponentOrHTML {
	return vecty.Text("Path does not exist")
}

var regexNamedVar = regexp.MustCompile("{[^/]+}")

var notFoundComponent *vecty.Component = new(vecty.Component)

var routes = []*Route{}

type Route struct {
	vecty.Core

	comp    vecty.Component
	pattern *regexp.Regexp
	path    string
}

func NewRoute(path string, c vecty.Component) *Route {
	r := &Route{
		comp: c,
		path: path,
	}

	if notFoundComponent == nil {
		*notFoundComponent = &defaultNotFound{}
	}

	pattern := path

	pattern = fmt.Sprintf("^%v$", pattern)

	if regexNamedVar.MatchString(path) {
		pattern = regexNamedVar.ReplaceAllString(path, "([^/]+)")
	}

	r.pattern = regexp.MustCompile(pattern)

	addRoute(r)

	return r
}

func NotFoundComponent(c vecty.Component) {
	*notFoundComponent = c
}

func addRoute(r *Route) {
	routes = append(routes, r)
}

func (r *Route) Render() vecty.ComponentOrHTML {
	path := pathname()
	if r.pattern.MatchString(path) {
		return r.comp
	}
	return *notFoundComponent
}

type VarGetter struct {
	Path    string
	Pattern *regexp.Regexp
}

func (v *VarGetter) All() map[string]string {
	var vars = make(map[string]string)

	if v == nil {
		return vars
	}
	// extract the named vars from url: "/users/{id}/{dog}" => ["{id}", "{dog}"]
	namedVars := regexNamedVar.FindAllString(v.Path, -1)

	// remove the surrounding brackets from each named var
	for i := 0; i < len(namedVars); i++ {
		namedVars[i] = strings.Replace(
			strings.Replace(namedVars[i], "{", "", 1),
			"}",
			"",
			1,
		)
	}

	namedValues := v.Pattern.FindAllStringSubmatch(pathname(), -1)[0][1:]

	for i := 0; i < len(namedVars); i++ {
		vars[namedVars[i]] = namedValues[i]
	}

	return vars
}

func (v *VarGetter) Get(key string) (string, error) {
	if s, ok := v.All()[key]; ok {
		return s, nil
	}
	return "", errors.New("named variable not found")
}

func (v *VarGetter) GetInt(key string) (int, error) {
	if s, ok := v.All()[key]; ok {
		return strconv.Atoi(s)
	}
	return -1, errors.New("named variable not found")
}

func GetVar(name string) (string, error) {
	for _, r := range routes {
		if r.pattern.MatchString(pathname()) {
			g := VarGetter{
				Path:    r.path,
				Pattern: r.pattern,
			}

			return g.Get(name)
		}
	}
	return "", errors.New("no url value found")
}

func GetIntVar(name string) (int, error) {
	for _, r := range routes {
		if r.pattern.MatchString(pathname()) {
			g := VarGetter{
				Path:    r.path,
				Pattern: r.pattern,
			}

			return g.GetInt(name)
		}
	}
	return 0, errors.New("no url value found")
}

func GetVarComp(c vecty.Component) *VarGetter {
	for i := range routes {
		if routes[i].comp == c {
			return &VarGetter{Path: routes[i].path, Pattern: routes[i].pattern}
		}
	}
	return nil
}
