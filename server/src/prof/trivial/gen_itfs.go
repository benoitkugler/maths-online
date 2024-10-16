package trivial

import (
	"encoding/json"

	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
)

// Code generated by gomacro/generator/gounions. DO NOT EDIT

// GroupsStrategyWrapper may be used as replacements for GroupsStrategy
// when working with JSON
type GroupsStrategyWrapper struct {
	Data GroupsStrategy
}

func (out *GroupsStrategyWrapper) UnmarshalJSON(src []byte) error {
	var wr struct {
		Kind string
		Data json.RawMessage
	}
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return err
	}
	switch wr.Kind {
	case "GroupsStrategyAuto":
		var data GroupsStrategyAuto
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "GroupsStrategyManual":
		var data GroupsStrategyManual
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data

	default:
		panic("exhaustive switch")
	}
	return err
}

func (item GroupsStrategyWrapper) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		Data interface{}
		Kind string
	}
	var wr wrapper
	switch data := item.Data.(type) {
	case GroupsStrategyAuto:
		wr = wrapper{Kind: "GroupsStrategyAuto", Data: data}
	case GroupsStrategyManual:
		wr = wrapper{Kind: "GroupsStrategyManual", Data: data}

	default:
		panic("exhaustive switch")
	}
	return json.Marshal(wr)
}

const (
	GroupsStrategyAutoGrKind   = "GroupsStrategyAuto"
	GroupsStrategyManualGrKind = "GroupsStrategyManual"
)

func (item LaunchSessionIn) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		IdConfig trivial.IdTrivial
		Groups   GroupsStrategyWrapper
	}
	wr := wrapper{
		IdConfig: item.IdConfig,
		Groups:   GroupsStrategyWrapper{item.Groups},
	}
	return json.Marshal(wr)
}

func (item *LaunchSessionIn) UnmarshalJSON(src []byte) error {
	type wrapper struct {
		IdConfig trivial.IdTrivial
		Groups   GroupsStrategyWrapper
	}
	var wr wrapper
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return err
	}
	item.IdConfig = wr.IdConfig
	item.Groups = wr.Groups.Data
	return nil
}
