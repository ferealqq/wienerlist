package components

import (
	"github.com/ferealqq/wienerlist/front/store"

	"github.com/ferealqq/wienerlist/front/components/bs"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type BoardComponent struct {
	vecty.Mounter
	vecty.Unmounter
	vecty.Keyer
	vecty.RenderSkipper
	vecty.Core
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
		vecty.Rerender(b)
	})
}

func (b *BoardComponent) SkipRender(prev vecty.Component) bool {
	if rs, ok := prev.(vecty.Keyer); ok {
		// if the index changes we need to fetch all the sections for this board
		if rs.Key() != b.Index {
			store.SectionState.Listeners.Remove(b)
			store.SectionState.Listeners.Add(b, func() {
				b.secs = store.SectionState.BoardSections[b.Index]
				vecty.Rerender(b)
			})
			b.secs = store.SectionState.BoardSections[b.Index]
		}
	}
	return false
}

func (b *BoardComponent) Render() vecty.ComponentOrHTML {
	// we can just spam this action call because it only fetches when the data dosn't exist
	go store.FetchBoardSectionsIfNeeded(b.Index)

	var secs vecty.List
	for _, sec := range b.secs {
		secs = append(secs, &sectionItem{section: sec})
	}

	return bs.ContainerFluid(
		bs.Row(
			vecty.Markup(
				vecty.Class("p-3"),
			),

			vecty.If(len(b.secs) > 0,
				secs,
			),
		),
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
