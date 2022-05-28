package components

import (
	"strconv"

	"github.com/ferealqq/wienerlist/front/actions"
	"github.com/ferealqq/wienerlist/front/dispatcher"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/ferealqq/wienerlist/front/store/state"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/style"
	router "marwan.io/vecty-router"
)

// WorkspaceList is a vecty.Component which represents a single item in the TODO
type WorkspaceList struct {
	vecty.Core

	Index      int                  `vecty:"prop"`
	Workspace  *state.WorkspaceData `vecty:"prop"`
	editing    bool
	showBoards bool
	editTitle  string
	input      *vecty.HTML
}

// Key implements the vecty.Keyer interface.
func (p *WorkspaceList) Key() interface{} {
	return p.Index
}

func (p *WorkspaceList) toggleBoards(event *vecty.Event) {
	p.showBoards = !p.showBoards
	vecty.Rerender(p)
}

func (p *WorkspaceList) onDestroy(event *vecty.Event) {
	dispatcher.Dispatch(&actions.DestroyItem{
		Index: p.Index,
	})
}

func (p *WorkspaceList) onToggleCompleted(event *vecty.Event) {
	dispatcher.Dispatch(&actions.SetCompleted{
		Index:     p.Index,
		Completed: event.Target.Get("checked").Bool(),
	})
}

func (p *WorkspaceList) onStartEdit(event *vecty.Event) {
	p.editing = true
	p.editTitle = p.Workspace.Title
	vecty.Rerender(p)
	p.input.Node().Call("focus")
}

func (p *WorkspaceList) onEditInput(event *vecty.Event) {
	p.editTitle = event.Target.Get("value").String()
	vecty.Rerender(p)
}

func (p *WorkspaceList) onStopEdit(event *vecty.Event) {
	p.editing = false
	vecty.Rerender(p)
	dispatcher.Dispatch(&actions.SetTitle{
		Index: p.Index,
		Title: p.editTitle,
	})
}

// Render implements the vecty.Component interface.
func (p *WorkspaceList) Render() vecty.ComponentOrHTML {
	return elem.ListItem(
		vecty.Markup(
			vecty.ClassMap{
				"editing": p.editing,
			},
			vecty.Class("py-1"),
		),

		elem.Div(
			vecty.Markup(
				vecty.Class("view"),
			),
			elem.Button(
				vecty.Markup(
					vecty.Class("btn"),
					vecty.Class("btn-primary"),
					event.Click(p.toggleBoards),
				),
				vecty.Text(p.Workspace.Title),
			),
		),
		vecty.If(p.showBoards,
			p.renderBoardList(),
		),
		elem.Form(
			vecty.Markup(
				style.Margin(style.Px(0)),
				event.Submit(p.onStopEdit).PreventDefault(),
			),
			p.input,
		),
	)
}

func (p *WorkspaceList) renderBoardList() *vecty.HTML {
	var items vecty.List
	for i := range p.Workspace.Boards {
		b := p.Workspace.Boards[i]
		items = append(items, &boardItem{Index: i, Board: b})
	}
	return elem.UnorderedList(
		vecty.Markup(
			vecty.Class("todo-list"),
		),
		items,
	)
}

type boardItem struct {
	vecty.Core

	Index int          `vecty:"prop"`
	Board *model.Board `vecty:"prop"`
}

func (p *boardItem) Render() vecty.ComponentOrHTML {

	return elem.ListItem(
		vecty.Markup(
			vecty.Class("border-1"),
			vecty.Class("rounded-1"),
		),
		elem.Button(
			vecty.Markup(
				vecty.Class("btn"),
				vecty.Class("btn-link"),
				event.Click(func(_ *vecty.Event) {
					router.Redirect("/boards/" + strconv.Itoa(int(p.Board.ID)))
				}),
			),
			vecty.Text(p.Board.Title),
		),
	)
}
