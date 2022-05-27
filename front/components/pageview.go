package components

import (
	"github.com/ferealqq/wienerlist/front/components/bs"
	services "github.com/ferealqq/wienerlist/front/store/services"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	router "marwan.io/vecty-router"
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
					LeftPanel(),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("col-10"),
					),
					router.NewRoute("/", &HomeContainer{}, router.NewRouteOpts{ExactMatch: true}),
					router.NewRoute("/boards/{id}", new(Test), router.NewRouteOpts{ExactMatch: true}),
				),
			),
		),
	)
}

type Test struct {
	vecty.Core
}

func (t *Test) Render() vecty.ComponentOrHTML {
	sId := router.GetNamedVar(t)["id"]
	return &BoardContainer{Index: sId}
}

// HomeContainer is a vecty.Component which represents the entire page.
type HomeContainer struct {
	vecty.Core
}

func (h *HomeContainer) Render() vecty.ComponentOrHTML {
	return elem.Div(vecty.Text("Homepage"))
}
