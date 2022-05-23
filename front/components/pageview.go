package components

import (
	"github.com/ferealqq/wienerlist/front/store/model"
	services "github.com/ferealqq/wienerlist/front/store/services"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

var api = services.NewApi("http://localhost:4000/api/v1")

// PageView is a vecty.Component which represents the entire page.
type PageView struct {
	vecty.Core
}

func (p *PageView) Render() vecty.ComponentOrHTML {
	// Render implements the vecty.Component interface.
	var allWs model.ListWorkspace
	if err := api.Get("/workspaces/").BindModel(&allWs); err != nil {
		//FIXME Handle errors?
		panic(err)
	}

	l := len(allWs.Workspaces)
	// List of pointers to workspace
	wsps := make([]*model.Workspace, 0, l)
	for i := 0; i != l; i++ {
		wsps = append(wsps, &allWs.Workspaces[i])
	}

	return elem.Body(
		vecty.Markup(
			vecty.Class("px-2"),
		),
		elem.Section(
			vecty.Markup(
				vecty.Class("todoapp"),
			),

			p.renderHeader(),
			vecty.If(len(allWs.Workspaces) > 0,
				p.renderWorkspaceList(wsps),
				p.renderFooter(),
			),
		),
	)
}

func (p *PageView) renderHeader() *vecty.HTML {
	return elem.Header(
		vecty.Markup(
			vecty.Class("header"),
		),

		elem.Heading1(
			vecty.Markup(
				vecty.Class("mun"),
			),
      
			vecty.Text("Workspaces!"),
		),
	)
}

func (p *PageView) renderFooter() *vecty.HTML {
	return elem.Footer(
		vecty.Markup(
			vecty.Class("footer"),
		),
	)
}

func (p *PageView) renderWorkspaceList(a []*model.Workspace) *vecty.HTML {
	var wsItems vecty.List
	for i, ws := range a {
		wsItems = append(wsItems, &WorkspaceView{Index: i, Workspace: ws})
	}

	return elem.Section(
		vecty.Markup(
			vecty.Class("main"),
		),

		elem.UnorderedList(
			vecty.Markup(
				vecty.Class("todo-list"),
			),
			wsItems,
		),
	)
}
