import type {
  Block,
  FigureBlock,
  FormulaBlock,
  FormulaFieldBlock,
  FunctionGraphBlock,
  FunctionVariationGraphBlock,
  NumberFieldBlock,
  OrderedListFieldBlock,
  RadioFieldBlock,
  SignTableBlock,
  TableBlock,
  TextBlock,
  TextPart,
  Variable,
  VariationTableBlock
} from "./exercice_gen";
import {
  BlockKind,
  ComparisonLevel,
  SignSymbol,
  TextKind
} from "./exercice_gen";

export const ExpressionColor = "orange";

export const colorByKind: { [key in TextKind]: string } = {
  [TextKind.Text]: "",
  [TextKind.StaticMath]: "green",
  [TextKind.Expression]: ExpressionColor
};

const reLaTeX = /\$([^$]*)\$/g;
const reExpression = /#{([^}]*)}/g;

function splitByRegexp(
  re: RegExp,
  s: string,
  kindMatch: TextKind,
  kindDefault: TextKind
): TextPart[] {
  const out: TextPart[] = [];
  const matches = s.matchAll(re);
  let cursor = 0;
  for (const match of matches) {
    const outerStart = match.index!;
    const outerEnd = match.index! + match[0].length;

    if (outerStart > cursor) {
      out.push({ Kind: kindDefault, Content: s.substring(cursor, outerStart) });
    }

    out.push({ Kind: kindMatch, Content: s.substring(outerStart, outerEnd) });

    cursor = outerEnd;
  }

  if (s.length > cursor) {
    out.push({ Kind: kindDefault, Content: s.substring(cursor, s.length) });
  }

  return out;
}

export function itemize(s: string): TextPart[] {
  const out: TextPart[] = [];
  splitByRegexp(reLaTeX, s, TextKind.StaticMath, TextKind.Text).forEach(
    chunk => {
      out.push(
        ...splitByRegexp(
          reExpression,
          chunk.Content,
          TextKind.Expression,
          chunk.Kind
        )
      );
    }
  );
  return out;
}

export const BlockKindLabels: { [T in BlockKind]: string } = {
  [BlockKind.TextBlock]: "Texte",
  [BlockKind.FormulaBlock]: "Formule",
  [BlockKind.FigureBlock]: "Figure",
  [BlockKind.FunctionGraphBlock]: "Graphe (expression)",
  [BlockKind.FunctionVariationGraphBlock]: "Graphe (variations)",
  [BlockKind.VariationTableBlock]: "Tableau de variations",
  [BlockKind.SignTableBlock]: "Tableau de signes",
  [BlockKind.TableBlock]: "Tableau",
  [BlockKind.NumberFieldBlock]: "Nombre",
  [BlockKind.FormulaFieldBlock]: "Expression",
  [BlockKind.RadioFieldBlock]: "QCM",
  [BlockKind.OrderedListFieldBlock]: "Liste ordonnée"
};

interface BlockKindTypes {
  [BlockKind.FigureBlock]: FigureBlock;
  [BlockKind.FormulaBlock]: FormulaBlock;
  [BlockKind.FormulaFieldBlock]: FormulaFieldBlock;
  [BlockKind.FunctionGraphBlock]: FunctionGraphBlock;
  [BlockKind.FunctionVariationGraphBlock]: FunctionVariationGraphBlock;
  [BlockKind.NumberFieldBlock]: NumberFieldBlock;
  [BlockKind.RadioFieldBlock]: RadioFieldBlock;
  [BlockKind.SignTableBlock]: SignTableBlock;
  [BlockKind.TableBlock]: TableBlock;
  [BlockKind.TextBlock]: TextBlock;
  [BlockKind.VariationTableBlock]: VariationTableBlock;
  [BlockKind.OrderedListFieldBlock]: OrderedListFieldBlock;
}

export interface TypedBlock<K extends BlockKind> {
  Kind: K;
  Data: BlockKindTypes[K];
}

export const xRune = "x".codePointAt(0)!;
export const yRune = "y".codePointAt(0)!;

/** extractPoints returns the names of indices A for which 'x_A' and 'y_A' are defined */
export function extractPoints(vars: Variable[]) {
  const points: { [key: string]: { x: boolean; y: boolean } } = {};
  vars.forEach(v => {
    if (v.Indice != "") {
      const point = points[v.Indice] || {};
      if (v.Name == xRune) {
        point.x = true;
      } else if (v.Name == yRune) {
        point.y = true;
      }
      points[v.Indice] = point;
    }
  });

  return Object.keys(points).filter(name => {
    const point = points[name];
    return point.x && point.y;
  });
}

export function saveData(data: any, fileName: string) {
  const a = document.createElement("a");
  document.body.appendChild(a);
  a.setAttribute("style", "display: none");
  const json = JSON.stringify(data, null, "  "),
    blob = new Blob([json], { type: "octet/stream" }),
    url = window.URL.createObjectURL(blob);
  a.href = url;
  a.download = fileName;
  a.click();
  window.URL.revokeObjectURL(url);
}

export function newBlock(kind: BlockKind): Block {
  switch (kind) {
    case BlockKind.TextBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          IsHint: false,
          Parts: ""
        }
      };
      return out;
    }
    case BlockKind.FormulaBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Parts: ""
        }
      };
      return out;
    }
    case BlockKind.FigureBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          ShowGrid: true,
          Bounds: {
            Width: 10,
            Height: 10,
            Origin: { X: 3, Y: 3 }
          },
          Drawings: {
            Lines: [],
            Points: [],
            Segments: []
          }
        }
      };
      return out;
    }
    case BlockKind.FunctionGraphBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Function: "abs(x) + sin(x)",
          Label: "f",
          Variable: { Name: xRune, Indice: "" },
          Range: [-5, 5]
        }
      };
      return out;
    }
    case BlockKind.FunctionVariationGraphBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Xs: ["-5", "0", "5"],
          Fxs: ["-3", "2", "-1"]
        }
      };
      return out;
    }
    case BlockKind.VariationTableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Xs: ["-5", "0", "5"],
          Fxs: ["-3", "2", "-1"]
        }
      };
      return out;
    }
    case BlockKind.SignTableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          FxSymbols: [
            SignSymbol.Nothing,
            SignSymbol.Zero,
            SignSymbol.ForbiddenValue,
            SignSymbol.Nothing
          ],
          Xs: ["\\infty", "3", "5", "+\\infty"],
          Signs: [true, false, true]
        }
      };
      return out;
    }
    case BlockKind.TableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          VerticalHeaders: [
            { Kind: TextKind.Text, Content: "Ligne 1" },
            { Kind: TextKind.Text, Content: "Ligne 2" }
          ],
          HorizontalHeaders: [
            { Kind: TextKind.Text, Content: "Colonne 1" },
            { Kind: TextKind.Text, Content: "Colonne 2" }
          ],
          Values: [
            [
              { Kind: TextKind.Text, Content: "Case" },
              { Kind: TextKind.StaticMath, Content: "\\frac{1}{2}" }
            ],
            [
              { Kind: TextKind.Expression, Content: "2x + 3" },
              { Kind: TextKind.StaticMath, Content: "18" }
            ]
          ]
        }
      };
      return out;
    }
    case BlockKind.NumberFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Expression: "1"
        }
      };
      return out;
    }
    case BlockKind.FormulaFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Label: { Kind: TextKind.Text, Content: "" },
          Expression: "x^2 + 2x + 1",
          ComparisonLevel: ComparisonLevel.SimpleSubstitutions
        }
      };
      return out;
    }
    case BlockKind.RadioFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: "1",
          Proposals: ["Oui", "Non"],
          AsDropDown: false
        }
      };
      return out;
    }
    case BlockKind.OrderedListFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: [
            { Kind: TextKind.StaticMath, Content: "[" },
            { Kind: TextKind.StaticMath, Content: "-12" },
            { Kind: TextKind.StaticMath, Content: ";" },
            { Kind: TextKind.StaticMath, Content: "30" },
            { Kind: TextKind.StaticMath, Content: "]" }
          ],
          AdditionalProposals: [],
          Label: "x \\in "
        }
      };
      return out;
    }
    default:
      throw "Unexpected Kind";
  }
}
