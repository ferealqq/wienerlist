package components

import (
	"log"
	"strconv"

	router "marwan.io/vecty-router"

	"github.com/ferealqq/wienerlist/front/components/bs"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type BoardContainer struct {
	vecty.Core
}

func (b *BoardContainer) Render() vecty.ComponentOrHTML {
	id, err := strconv.Atoi(router.GetNamedVar(b)["id"])
	if err != nil {
		return vecty.Text("Invalid board id")
	}
	var secs model.ListSections
	// TODO Create wrapper actions to get section data
	if err := api.Params("board_id", id).Get("/sections/").BindModel(&secs); err != nil {
		log.Println(err.Error())
		return vecty.Text("Something went wrong try again later!")
	}
	var sectionItems vecty.List
	for _, section := range secs.Sections {
		sectionItems = append(sectionItems, &sectionItem{section: &section})
	}
	return bs.Row(
		vecty.Markup(
			vecty.Class("p-3"),
		),

		sectionItems,
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
