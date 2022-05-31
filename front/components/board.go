package components

import (
	"strconv"
	"time"

	"github.com/ferealqq/wienerlist/front/store"
	router "marwan.io/vecty-router"

	"github.com/ferealqq/wienerlist/front/components/bs"
	u "github.com/ferealqq/wienerlist/front/components/util"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/hexops/vecty"
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type BoardComponent struct {
	v.Mounter
	v.Unmounter
	v.Keyer
	v.RenderSkipper
	v.Core
	// Index => BoardId
	Index int `vecty:"prop"`
	secs  map[int]*model.Section
	board *model.Board
}

func (b *BoardComponent) Key() interface{} {
	return b.Index
}

func (b *BoardComponent) Mount() {
	store.SectionState.Listeners.Add(b, func() {
		b.secs = store.SectionState.BoardSections[b.Index]
		v.Rerender(b)
	})
}

func (b *BoardComponent) SkipRender(prev v.Component) bool {
	if rs, ok := prev.(v.Keyer); ok {
		// if the index changes we need to fetch all the sections for this board
		if rs.Key() != b.Index {
			store.SectionState.Listeners.Remove(b)
			store.SectionState.Listeners.Add(b, func() {
				b.secs = store.SectionState.BoardSections[b.Index]
				v.Rerender(b)
			})
			b.secs = store.SectionState.BoardSections[b.Index]
		}
	}
	return false
}

func (b *BoardComponent) Render() v.ComponentOrHTML {
	// we can just spam this action call because it only fetches when the data dosn't exist
	go store.FetchBoardSectionsIfNeeded(b.Index)

	var secs v.List
	for _, sec := range b.secs {
		secs = append(secs, &sectionItem{section: sec})
	}

	return bs.ContainerFluid(
		bs.Row(
			v.Markup(
				v.Class("p-3"),
			),

			v.If(len(b.secs) > 0,
				secs,
			),
		),
		// TODO Fix this
		// router.NewRoute("/item/{itemId}", new(Test), router.NewRouteOpts{ExactMatch: true}),
	)
}

type sectionItem struct {
	v.Core

	section *model.Section `vecty:"prop"`
}

func (s *sectionItem) Render() v.ComponentOrHTML {
	var list v.List
	for j := range s.section.Items {
		list = append(list, &itemComponent{item: &s.section.Items[j]})
	}

	return e.Div(
		v.Markup(
			v.Class("col"),
		),

		v.Text(s.section.Title),
		v.If(len(s.section.Items) > 0,
			bs.List(
				list,
			),
		),
	)
}

type itemComponent struct {
	v.Core
	item   *model.Item `vecty:"prop"`
	toggle bool

	created string
	updated string
}

const (
	day   = time.Hour * 24
	week  = day * 7
	month = day * 30
	year  = day * 365
)

// time.Duration in Nanoseconds
func timeExt(na int64) string {
	d := na / int64(day)
	s := na / int64(time.Second)
	if d == 365 {
		return "year"
	} else if d > 365 {
		return strconv.Itoa(int(d/365)) + " years"
	} else if d == 28 || d == 29 || d == 30 || d == 31 {
		return "month"
	} else if d > 30 {
		// FIX this
		return strconv.Itoa(int(d/30)) + " months"
	} else if d < 14 {
		return "week"
	} else if d > 7 {
		return strconv.Itoa(int(d/7)) + " weeks"
	} else if d == 1 {
		return "day"
	} else if d > 1 {
		return strconv.Itoa(int(d)) + " days"
	} else if s == (60 * 60) {
		return "hour"
	} else if s > (60 * 60) {
		return strconv.Itoa(int(s/(60*60))) + " hours"
	} else if s == 60 {
		return "minute"
	} else if s > 60 {
		return strconv.Itoa(int(s/60)) + " minutes"
	} else if s == 1 {
		return "second"
	} else {
		return strconv.Itoa(int(s/60)) + " seconds"
	}
}

func since(s string) string {
	if t, err := time.Parse(time.RFC3339Nano, s); err != nil {
		return ""
	} else {
		sin := time.Since(t)
		return timeExt(sin.Nanoseconds()) + " ago"
	}
}

func (i *itemComponent) Render() v.ComponentOrHTML {
	id := "item-" + strconv.Itoa(int(i.item.ID))
	i.created = since(i.item.CreatedAt)
	i.updated = since(i.item.UpdatedAt)
	// i.toggle = true
	return bs.ListItem(
		v.Markup(
			event.Click(func(_ *vecty.Event) {
				router.Redirect("/boards/1/item/" + strconv.Itoa(int(i.item.ID)))
			}),
			event.MouseEnter(func(_ *v.Event) {
				i.toggle = !i.toggle
				v.Rerender(i)
			}),
			event.MouseLeave(func(_ *v.Event) {
				i.toggle = !i.toggle
				v.Rerender(i)
			}),
		),
		bs.Active(i.toggle),
		e.Div(
			u.Classes(
				"d-flex",
				"w-100",
				"justify-content-between",
			),
			u.Atr("id", id),
			e.Heading5(
				u.Classes("mb-1"),
				v.Text(i.item.Title),
			),
			v.If(
				!i.toggle,
				e.Small(v.Text(i.created)),
			),
			v.If(
				i.toggle,
				bs.Icon("github"),
			),
		),
		v.If(i.toggle,
			e.Div(
				u.Classes("col", "pt-1", "d-flex", "justify-content-between"),
				e.Paragraph(
					u.Classes("small", "my-auto"),
					bs.Icon("calendar2-fill", u.Classes("pe-1")),
					vecty.Text(i.created),
				),
				e.Paragraph(
					u.Classes("small", "my-auto"),
					bs.Icon("calendar2-plus-fill", u.Classes("pe-1")),
					vecty.Text(i.updated),
				),
			),
		),
	)
}
