package components

import (
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func LeftPanel() vecty.ComponentOrHTML {
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

	return elem.Div(
		renderHeader(),
		vecty.If(len(allWs.Workspaces) > 0,
			renderWorkspaceList(wsps),
			renderFooter(),
		),
	)
}

func renderHeader() *vecty.HTML {
	return elem.Heading1(
		vecty.Markup(
			vecty.Class("display-6"),
			vecty.Class("bold"),
		),

		vecty.Text("Workspaces!"),
	)
}
func renderFooter() *vecty.HTML {
	return elem.Footer(
		vecty.Markup(
			vecty.Class("footer"),
		),
	)
}

func renderWorkspaceList(a []*model.Workspace) *vecty.HTML {
	var wsItems vecty.List
	for i, ws := range a {
		wsItems = append(wsItems, &WorkspaceList{Index: i, Workspace: ws})
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
