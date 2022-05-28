package components

import (
	"github.com/ferealqq/wienerlist/front/store"

	"github.com/ferealqq/wienerlist/front/components/bs"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type Board struct {
	vecty.Mounter
	vecty.Unmounter
	vecty.Keyer
	vecty.RenderSkipper
	vecty.Core
	// Index => BoardId
	Index int `vecty:"prop"`
	secs  []*model.Section
}

func (b *Board) Key() interface{} {
	return b.Index
}

func (b *Board) Mount() {
	store.Listeners.Add(b, func() {
		b.secs = store.SectionState.BoardSections[b.Index]
		vecty.Rerender(b)
	})
}

func (b *Board) SkipRender(prev vecty.Component) bool {
	if rs, ok := prev.(vecty.Keyer); ok {
		// if the index changes we need to fetch all the sections for this board
		if rs.Key() != b.Index {
			store.Listeners.Remove(b)
			store.Listeners.Add(b, func() {
				b.secs = store.SectionState.BoardSections[b.Index]
				vecty.Rerender(b)
			})
			b.secs = store.SectionState.BoardSections[b.Index]
		}
	}
	return false
}

func (b *Board) Render() vecty.ComponentOrHTML {
	// we can just spam this action call because it only fetches when the data dosn't exist
	go store.FetchBoardSectionsIfNeeded(b.Index)
	return bs.Row(
		vecty.Markup(
			vecty.Class("p-3"),
		),

		vecty.If(len(b.secs) > 0,
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
