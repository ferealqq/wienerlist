package components

import (
	"log"
	"strconv"

	"github.com/ferealqq/wienerlist/front/components/bs"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type BoardContainer struct {
	vecty.Core

	Index     string `vecty:"prop"`
	hasLoaded bool
	secs      []*model.Section
}

func (b *BoardContainer) Key() interface{} {
	return b.Index
}

func (b *BoardContainer) fetchSections() {
	id, err := strconv.Atoi(b.Index)
	if err != nil {
		log.Fatal(err.Error())
		// return vecty.Text("Invalid board id")
	}
	var secs model.ListSections
	// TODO Create wrapper actions to get section data
	if err := api.Params("board_id", id).Get("/sections/").BindModel(&secs); err != nil {
		log.Fatal(err.Error())
		// return vecty.Text("Something went wrong try again later!")
	}
	size := len(secs.Sections)

	wsps := make([]*model.Section, 0, size)
	for i := 0; i != size; i++ {
		wsps = append(wsps, &secs.Sections[i])
	}

	b.secs = wsps
	b.hasLoaded = true
	vecty.Rerender(b)
}

func (b *BoardContainer) Render() vecty.ComponentOrHTML {
	if !b.hasLoaded {
		go b.fetchSections()
	}
	return bs.Row(
		vecty.Markup(
			vecty.Class("p-3"),
		),

		vecty.If(b.hasLoaded,
			renderSections(b.secs),
		),
	)
}

func renderSections(sections []*model.Section) *vecty.HTML {
	var secs vecty.List
	for _, sec := range sections {
		secs = append(secs, &sectionItem{section: sec})
	}
	return elem.Div(
		secs,
	)
}

type sectionItem struct {
	vecty.Core

	section *model.Section
}

func (s *sectionItem) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("col"),
		),

		vecty.Text(s.section.Title),
	)
}
