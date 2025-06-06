package preview

import "encoding/json"

// Code generated by gomacro/generator/gounions. DO NOT EDIT

// LoopbackServerEventWrapper may be used as replacements for LoopbackServerEvent
// when working with JSON
type LoopbackServerEventWrapper struct {
	Data LoopbackServerEvent
}

func (out *LoopbackServerEventWrapper) UnmarshalJSON(src []byte) error {
	var wr struct {
		Kind string
		Data json.RawMessage
	}
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return err
	}
	switch wr.Kind {
	case "LoopbackPaused":
		var data LoopbackPaused
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "LoopbackShowCeinture":
		var data LoopbackShowCeinture
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "LoopbackShowExercice":
		var data LoopbackShowExercice
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "LoopbackShowQuestion":
		var data LoopbackShowQuestion
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data

	default:
		panic("exhaustive switch")
	}
	return err
}

func (item LoopbackServerEventWrapper) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		Data interface{}
		Kind string
	}
	var wr wrapper
	switch data := item.Data.(type) {
	case LoopbackPaused:
		wr = wrapper{Kind: "LoopbackPaused", Data: data}
	case LoopbackShowCeinture:
		wr = wrapper{Kind: "LoopbackShowCeinture", Data: data}
	case LoopbackShowExercice:
		wr = wrapper{Kind: "LoopbackShowExercice", Data: data}
	case LoopbackShowQuestion:
		wr = wrapper{Kind: "LoopbackShowQuestion", Data: data}

	default:
		panic("exhaustive switch")
	}
	return json.Marshal(wr)
}

const (
	LoopbackPausedLoKind       = "LoopbackPaused"
	LoopbackShowCeintureLoKind = "LoopbackShowCeinture"
	LoopbackShowExerciceLoKind = "LoopbackShowExercice"
	LoopbackShowQuestionLoKind = "LoopbackShowQuestion"
)
