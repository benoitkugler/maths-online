package questions

import "encoding/json"

// Code generated by structgen/interfaces. DO NOT EDIT

// BlockWrapper may be used as replacements for Block
// when working with JSON
type BlockWrapper struct {
	Data Block
}

func (out *BlockWrapper) UnmarshalJSON(src []byte) error {
	var wr struct {
		Kind string
		Data json.RawMessage
	}
	err := json.Unmarshal(src, &wr)
	if err != nil {
		return err
	}
	switch wr.Kind {
	case "ExpressionFieldBlock":
		var data ExpressionFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FigureAffineLineFieldBlock":
		var data FigureAffineLineFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FigureBlock":
		var data FigureBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FigurePointFieldBlock":
		var data FigurePointFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FigureVectorFieldBlock":
		var data FigureVectorFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FigureVectorPairFieldBlock":
		var data FigureVectorPairFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FormulaBlock":
		var data FormulaBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FunctionPointsFieldBlock":
		var data FunctionPointsFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "FunctionsGraphBlock":
		var data FunctionsGraphBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "NumberFieldBlock":
		var data NumberFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "OrderedListFieldBlock":
		var data OrderedListFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "ProofFieldBlock":
		var data ProofFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "RadioFieldBlock":
		var data RadioFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "SignTableBlock":
		var data SignTableBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "TableBlock":
		var data TableBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "TableFieldBlock":
		var data TableFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "TextBlock":
		var data TextBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "TreeFieldBlock":
		var data TreeFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "VariationTableBlock":
		var data VariationTableBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "VariationTableFieldBlock":
		var data VariationTableFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data
	case "VectorFieldBlock":
		var data VectorFieldBlock
		err = json.Unmarshal(wr.Data, &data)
		out.Data = data

	default:
		panic("exhaustive switch")
	}
	return err
}

func (item BlockWrapper) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		Data interface{}
		Kind string
	}
	var wr wrapper
	switch data := item.Data.(type) {
	case ExpressionFieldBlock:
		wr = wrapper{Kind: "ExpressionFieldBlock", Data: data}
	case FigureAffineLineFieldBlock:
		wr = wrapper{Kind: "FigureAffineLineFieldBlock", Data: data}
	case FigureBlock:
		wr = wrapper{Kind: "FigureBlock", Data: data}
	case FigurePointFieldBlock:
		wr = wrapper{Kind: "FigurePointFieldBlock", Data: data}
	case FigureVectorFieldBlock:
		wr = wrapper{Kind: "FigureVectorFieldBlock", Data: data}
	case FigureVectorPairFieldBlock:
		wr = wrapper{Kind: "FigureVectorPairFieldBlock", Data: data}
	case FormulaBlock:
		wr = wrapper{Kind: "FormulaBlock", Data: data}
	case FunctionPointsFieldBlock:
		wr = wrapper{Kind: "FunctionPointsFieldBlock", Data: data}
	case FunctionsGraphBlock:
		wr = wrapper{Kind: "FunctionsGraphBlock", Data: data}
	case NumberFieldBlock:
		wr = wrapper{Kind: "NumberFieldBlock", Data: data}
	case OrderedListFieldBlock:
		wr = wrapper{Kind: "OrderedListFieldBlock", Data: data}
	case ProofFieldBlock:
		wr = wrapper{Kind: "ProofFieldBlock", Data: data}
	case RadioFieldBlock:
		wr = wrapper{Kind: "RadioFieldBlock", Data: data}
	case SignTableBlock:
		wr = wrapper{Kind: "SignTableBlock", Data: data}
	case TableBlock:
		wr = wrapper{Kind: "TableBlock", Data: data}
	case TableFieldBlock:
		wr = wrapper{Kind: "TableFieldBlock", Data: data}
	case TextBlock:
		wr = wrapper{Kind: "TextBlock", Data: data}
	case TreeFieldBlock:
		wr = wrapper{Kind: "TreeFieldBlock", Data: data}
	case VariationTableBlock:
		wr = wrapper{Kind: "VariationTableBlock", Data: data}
	case VariationTableFieldBlock:
		wr = wrapper{Kind: "VariationTableFieldBlock", Data: data}
	case VectorFieldBlock:
		wr = wrapper{Kind: "VectorFieldBlock", Data: data}

	default:
		panic("exhaustive switch")
	}
	return json.Marshal(wr)
}

const (
	ExpressionFieldBlockBlKind       = "ExpressionFieldBlock"
	FigureAffineLineFieldBlockBlKind = "FigureAffineLineFieldBlock"
	FigureBlockBlKind                = "FigureBlock"
	FigurePointFieldBlockBlKind      = "FigurePointFieldBlock"
	FigureVectorFieldBlockBlKind     = "FigureVectorFieldBlock"
	FigureVectorPairFieldBlockBlKind = "FigureVectorPairFieldBlock"
	FormulaBlockBlKind               = "FormulaBlock"
	FunctionPointsFieldBlockBlKind   = "FunctionPointsFieldBlock"
	FunctionsGraphBlockBlKind        = "FunctionsGraphBlock"
	NumberFieldBlockBlKind           = "NumberFieldBlock"
	OrderedListFieldBlockBlKind      = "OrderedListFieldBlock"
	ProofFieldBlockBlKind            = "ProofFieldBlock"
	RadioFieldBlockBlKind            = "RadioFieldBlock"
	SignTableBlockBlKind             = "SignTableBlock"
	TableBlockBlKind                 = "TableBlock"
	TableFieldBlockBlKind            = "TableFieldBlock"
	TextBlockBlKind                  = "TextBlock"
	TreeFieldBlockBlKind             = "TreeFieldBlock"
	VariationTableBlockBlKind        = "VariationTableBlock"
	VariationTableFieldBlockBlKind   = "VariationTableFieldBlock"
	VectorFieldBlockBlKind           = "VectorFieldBlock"
)

func (ct Enonce) MarshalJSON() ([]byte, error) {
	tmp := make([]BlockWrapper, len(ct))
	for i, v := range ct {
		tmp[i].Data = v
	}
	return json.Marshal(tmp)
}

func (ct *Enonce) UnmarshalJSON(data []byte) error {
	var tmp []BlockWrapper
	err := json.Unmarshal(data, &tmp)
	*ct = make(Enonce, len(tmp))
	for i, v := range tmp {
		(*ct)[i] = v.Data
	}
	return err
}
