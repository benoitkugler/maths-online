package trivial


// Code generated by structgen/interfaces. DO NOT EDIT

// // GroupStrategyWrapper may be used as replacements for GroupStrategy
// // when working with JSON
// type GroupStrategyWrapper struct {
// 	Data GroupStrategy
// }

// func (out *GroupStrategyWrapper) UnmarshalJSON(src []byte) error {
// 	var wr struct {
// 		Data json.RawMessage
// 		Kind int
// 	}
// 	err := json.Unmarshal(src, &wr)
// 	if err != nil {
// 		return err
// 	}
// 	switch wr.Kind {
// 	case 0:
// 		var data FixedSizeGroupStrategy
// 		err = json.Unmarshal(wr.Data, &data)
// 		out.Data = data
// 	case 1:
// 		var data RandomGroupStrategy
// 		err = json.Unmarshal(wr.Data, &data)
// 		out.Data = data

// 	default:
// 		panic("exhaustive switch")
// 	}
// 	return err
// }

// func (item GroupStrategyWrapper) MarshalJSON() ([]byte, error) {
// 	type wrapper struct {
// 		Data interface{}
// 		Kind int
// 	}
// 	var wr wrapper
// 	switch data := item.Data.(type) {
// 	case FixedSizeGroupStrategy:
// 		wr = wrapper{Kind: 0, Data: data}
// 	case RandomGroupStrategy:
// 		wr = wrapper{Kind: 1, Data: data}

// 	default:
// 		panic("exhaustive switch")
// 	}
// 	return json.Marshal(wr)
// }

const (
	FixedSizeGroupStrategyGrKind = iota
	RandomGroupStrategyGrKind
)
