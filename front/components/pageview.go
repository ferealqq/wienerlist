package components

import (
	"github.com/ferealqq/wienerlist/front/components/bs"
	services "github.com/ferealqq/wienerlist/front/store/services"
	"github.com/ferealqq/wienerlist/front/util"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

var api = services.NewApi("http://localhost:4000/api/v1")

type PageView struct {
	vecty.Core
}

func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		vecty.Markup(
			vecty.Class("px-2"),
		),
		bs.ContainerFluid(
			bs.Row(
				elem.Div(
					vecty.Markup(
						vecty.Class("col-2"),
					),
					&LeftPanel{},
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("col-10"),
					),
					util.NewRoute("/", &HomeContainer{}),
					util.NewRoute("/boards/{id}/item/{itemId}", new(ItemModal)),
					util.NewRoute("/boards/{id}", new(BoardContainer)),
				),
			),
		),
	)
}

type BoardContainer struct {
	vecty.Core
}

func (b *BoardContainer) Render() vecty.ComponentOrHTML {
	id, err := util.GetIntVar("id")
	if err != nil {
		return vecty.Text("Invalid board id")
	}
	return &BoardComponent{Index: id}
}

// HomeContainer is a vecty.Component which represents the entire page.
type HomeContainer struct {
	vecty.Core
}

func (h *HomeContainer) Render() vecty.ComponentOrHTML {
	return elem.Div(vecty.Text("Homepage"))
}
