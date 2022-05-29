package state

import (
	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/ferealqq/wienerlist/front/store/storeutil"
)

// First int represents parent id second int represents the Models id
type GenDataStore[M interface{}] map[int]map[int]*M

type defState struct {
	IsFetching       bool
	LastActionFailed bool
	Error            error
	Listeners        *storeutil.ListenerRegistry
}

func defaultState() defState {
	return defState{
		IsFetching:       false,
		LastActionFailed: false,
		Error:            nil,
		Listeners:        storeutil.NewListenerRegistry(),
	}
}

type SectionState struct {
	// int represents board id
	defState
	BoardSections GenDataStore[model.Section]
}

func NewSectionState() *SectionState {
	return &SectionState{
		defState:      defaultState(),
		BoardSections: make(map[int]map[int]*model.Section),
	}
}

type WorkspaceStore map[int]*WorkspaceData

type WorkspaceState struct {
	defState
	// int represents workspace id
	Workspaces WorkspaceStore
}

type WorkspaceData struct {
	*model.Workspace
	// int represents board id
	Boards map[int]*model.Board
}

func NewWorkspaceState() *WorkspaceState {
	return &WorkspaceState{
		defState:   defaultState(),
		Workspaces: make(map[int]*WorkspaceData),
	}
}
