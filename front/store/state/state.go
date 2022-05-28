package state

import "github.com/ferealqq/wienerlist/front/store/model"

type SectionState struct {
	// int represents board id
	BoardSections map[int][]*model.Section

	IsFetching       bool
	LastActionFailed bool
	Error            error
}

func NewSectionState() SectionState {
	return SectionState{
		BoardSections:    make(map[int][]*model.Section),
		IsFetching:       false,
		LastActionFailed: false,
		Error:            nil,
	}
}
