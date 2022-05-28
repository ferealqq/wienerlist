package components

import (
	"github.com/ferealqq/wienerlist/front/store"
	"github.com/ferealqq/wienerlist/front/store/state"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type LeftPanel struct {
	vecty.Mounter
	vecty.Core

	ws state.WorkspaceStore
}

func (l *LeftPanel) Mount() {
	store.WorkspaceState.Listeners.Add(l, func() {
		l.ws = store.WorkspaceState.Workspaces
		vecty.Rerender(l)
	})
}

func (l *LeftPanel) Render() vecty.ComponentOrHTML {
	// maybe this should be moved to mount? but it probably won't matter because it does only one operations and exitsy
	go store.FetchWorkspacesIfNeeded()

	return elem.Div(
		renderHeader(),
		vecty.If(len(l.ws) > 0,
			renderWorkspaceList(l.ws),
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

func renderWorkspaceList(a state.WorkspaceStore) *vecty.HTML {
	var wsItems vecty.List
	for i := range a {
		wsItems = append(wsItems, &WorkspaceList{Index: i, Workspace: a[i]})
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
