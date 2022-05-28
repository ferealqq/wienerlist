package store

import (
	"github.com/ferealqq/wienerlist/front/actions"
	"github.com/ferealqq/wienerlist/front/dispatcher"
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/ferealqq/wienerlist/front/store/services"
	"github.com/ferealqq/wienerlist/front/store/state"
)

var api = services.NewApi("http://localhost:4000/api/v1")

var (
	// Items represents all of the TODO items in the store.
	Items []*model.Item

	// Filter represents the active viewing filter for items.
	Filter         = model.All
	SectionState   = state.NewSectionState()
	WorkspaceState = state.NewWorkspaceState()
)

func init() {
	registers := []func(action interface{}){onWorkspaceAction, onSectionAction}
	for _, r := range registers {
		dispatcher.Register(r)
	}
}

func FetchBoardSectionsIfNeeded(boardId int) error {
	if _, ok := SectionState.BoardSections[boardId]; !ok && !SectionState.IsFetching {
		dispatcher.Dispatch(&actions.FetchSectionsRequest{})
		var secs model.ListSections
		// TODO Create wrapper actions to get section data
		if err := api.Params("board_id", boardId).Get("/sections/").BindModel(&secs); err != nil {
			dispatcher.Dispatch(&actions.FetchSectionsResponseError{Error: err})
			return err
		}

		dispatcher.Dispatch(&actions.FetchSectionsResponse{
			Sections: secs.Sections,
			BoardId:  boardId,
		})
		return nil
	}

	return nil
}

func FetchWorkspacesIfNeeded() error {
	if len(WorkspaceState.Workspaces) == 0 && !WorkspaceState.IsFetching {
		dispatcher.Dispatch(&actions.FetchWorkspacesRequest{})
		// Render implements the vecty.Component interface.
		var allWs model.ListWorkspace
		if err := api.Get("/workspaces/").BindModel(&allWs); err != nil {
			dispatcher.Dispatch(&actions.FetchWorkspacesResponseError{Error: err})
			return err
		}
		dispatcher.Dispatch(&actions.FetchWorkspacesResponse{ListWorkspace: allWs})
	}

	return nil
}

// ActiveItemCount returns the current number of items that are not completed.
func ActiveItemCount() int {
	return count(false)
}

// CompletedItemCount returns the current number of items that are completed.
func CompletedItemCount() int {
	return count(true)
}

func count(completed bool) int {
	count := 0
	for _, item := range Items {
		if item.Completed == completed {
			count++
		}
	}
	return count
}

func onSectionAction(action interface{}) {
	switch a := action.(type) {
	case *actions.FetchSectionsRequest:
		SectionState.IsFetching = true
		SectionState.LastActionFailed = false

	case *actions.FetchSectionsResponseError:
		SectionState.IsFetching = false
		SectionState.LastActionFailed = true
		SectionState.Error = a.Error

	case *actions.FetchSectionsResponse:
		l := len(a.Sections)
		// List of pointers
		secs := make(map[int]*model.Section, l)
		for i := 0; i != l; i++ {
			c := a.Sections[i]
			secs[int(c.ID)] = &c
		}
		SectionState.BoardSections[a.BoardId] = secs
		SectionState.LastActionFailed = false
		SectionState.IsFetching = false
	default:
		return // don't fire listeners
	}

	SectionState.Listeners.Fire()
}

func onWorkspaceAction(action interface{}) {
	switch a := action.(type) {
	case *actions.FetchWorkspacesRequest:
		WorkspaceState.IsFetching = true
		WorkspaceState.LastActionFailed = false

	case *actions.FetchWorkspacesResponseError:
		WorkspaceState.IsFetching = false
		WorkspaceState.LastActionFailed = true
		WorkspaceState.Error = a.Error

	case *actions.FetchWorkspacesResponse:
		// List of pointers
		ws := make(state.WorkspaceStore, len(a.Workspaces))
		for i := range a.Workspaces {
			w := a.Workspaces[i]
			bs := make(map[int]*model.Board)
			for j := range w.Boards {
				bs[int(w.Boards[j].ID)] = &w.Boards[j]
			}
			ws[int(w.ID)] = &state.WorkspaceData{
				Workspace: &w,
				Boards:    bs,
			}
		}
		WorkspaceState.Workspaces = ws
		WorkspaceState.LastActionFailed = false
		WorkspaceState.IsFetching = false
	default:
		return // don't fire listeners
	}

	WorkspaceState.Listeners.Fire()
}
