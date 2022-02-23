package exercice

import "encoding/json"

func blockUnmarshallJSON(src []byte) (block, error) {
	type wrapper struct {
		Data json.RawMessage
		Kind int
	}
	var wr wrapper
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return nil, err
	}
	switch wr.Kind {
	case 0:
		var out Formula
		err = json.Unmarshal(wr.Data, &out)
		return out, err
	case 1:
		var out FormulaField
		err = json.Unmarshal(wr.Data, &out)
		return out, err
	case 2:
		var out ListField
		err = json.Unmarshal(wr.Data, &out)
		return out, err
	case 3:
		var out TextBlock
		err = json.Unmarshal(wr.Data, &out)
		return out, err

	default:
		panic("exhaustive switch")
	}
}

func blockMarshallJSON(item block) ([]byte, error) {
	type wrapper struct {
		Data interface{}
		Kind int
	}
	var out wrapper
	switch item.(type) {
	case Formula:
		out = wrapper{Kind: 0, Data: item}
	case FormulaField:
		out = wrapper{Kind: 1, Data: item}
	case ListField:
		out = wrapper{Kind: 2, Data: item}
	case TextBlock:
		out = wrapper{Kind: 3, Data: item}

	default:
		panic("exhaustive switch")
	}
	return json.Marshal(out)
}
