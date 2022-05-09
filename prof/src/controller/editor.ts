import type { QuestionHeader } from "./api_gen";
import type {
  Block,
  CoordExpression,
  FigureAffineLineFieldBlock,
  FigureBlock,
  FigurePointFieldBlock,
  FigureVectorFieldBlock,
  FigureVectorPairFieldBlock,
  FormulaBlock,
  FormulaFieldBlock,
  FunctionGraphBlock,
  FunctionPointsFieldBlock,
  FunctionVariationGraphBlock,
  NumberFieldBlock,
  OrderedListFieldBlock,
  RadioFieldBlock,
  SignTableBlock,
  TableBlock,
  TableFieldBlock,
  TextBlock,
  TextPart,
  TreeFieldBlock,
  Variable,
  VariationTableBlock,
  VariationTableFieldBlock
} from "./exercice_gen";
import {
  BlockKind,
  ComparisonLevel,
  DifficultyTag,
  SignSymbol,
  TextKind,
  VectorPairCriterion
} from "./exercice_gen";

export const ExpressionColor = "orange";

export const colorByKind: { [key in TextKind]: string } = {
  [TextKind.Text]: "",
  [TextKind.StaticMath]: "green",
  [TextKind.Expression]: ExpressionColor
};

const reLaTeX = /\$([^$]+)\$/g;
const reExpression = /&([^&]+)&/g;

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

export const sortedBlockKindLabels = [
  [BlockKind.TextBlock, { label: "Texte", isAnswerField: false }],
  [BlockKind.FormulaBlock, { label: "Formule", isAnswerField: false }],
  [BlockKind.FigureBlock, { label: "Figure", isAnswerField: false }],
  [
    BlockKind.FunctionGraphBlock,
    {
      label: "Graphe (expression)",
      isAnswerField: false
    }
  ],
  [
    BlockKind.FunctionVariationGraphBlock,
    {
      label: "Graphe (variations)",
      isAnswerField: false
    }
  ],
  [BlockKind.TableBlock, { label: "Tableau", isAnswerField: false }],
  [
    BlockKind.SignTableBlock,
    {
      label: "Tableau de signes",
      isAnswerField: false
    }
  ],
  [
    BlockKind.VariationTableBlock,
    {
      label: "Tableau de variations",
      isAnswerField: false
    }
  ],
  [BlockKind.NumberFieldBlock, { label: "Nombre", isAnswerField: true }],
  [
    BlockKind.FormulaFieldBlock,
    {
      label: "Expression",
      isAnswerField: true
    }
  ],
  [
    BlockKind.OrderedListFieldBlock,
    {
      label: "Liste ordonnée",
      isAnswerField: true
    }
  ],
  [BlockKind.RadioFieldBlock, { label: "QCM", isAnswerField: true }],
  [
    BlockKind.FigurePointFieldBlock,
    {
      label: "Point (sur une figure)",
      isAnswerField: true
    }
  ],
  [
    BlockKind.FigureVectorFieldBlock,
    {
      label: "Vecteur (sur une figure)",
      isAnswerField: true
    }
  ],
  [
    BlockKind.FunctionPointsFieldBlock,
    {
      label: "Construction de fonction",
      isAnswerField: true
    }
  ],
  [
    BlockKind.VariationTableFieldBlock,
    {
      label: "Tableau de variations",
      isAnswerField: true
    }
  ],
  [
    BlockKind.FigureAffineLineFieldBlock,
    {
      label: "Fonction affine",
      isAnswerField: true
    }
  ],
  [
    BlockKind.FigureVectorPairFieldBlock,
    {
      label: "Construction de vecteurs",
      isAnswerField: true
    }
  ],
  [BlockKind.TableFieldBlock, { label: "Tableau", isAnswerField: true }],
  [BlockKind.TreeFieldBlock, { label: "Arbre", isAnswerField: true }]
] as const;

export const BlockKindLabels: {
  [T in BlockKind]: { label: string; isAnswerField: boolean };
} = Object.fromEntries(sortedBlockKindLabels);

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
  [BlockKind.FigurePointFieldBlock]: FigurePointFieldBlock;
  [BlockKind.FigureVectorFieldBlock]: FigureVectorFieldBlock;
  [BlockKind.VariationTableFieldBlock]: VariationTableFieldBlock;
  [BlockKind.FunctionPointsFieldBlock]: FunctionPointsFieldBlock;
  [BlockKind.FigureAffineLineFieldBlock]: FigureAffineLineFieldBlock;
  [BlockKind.FigureVectorPairFieldBlock]: FigureVectorPairFieldBlock;
  [BlockKind.TreeFieldBlock]: TreeFieldBlock;
  [BlockKind.TableFieldBlock]: TableFieldBlock;
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
          Parts: "",
          Bold: false,
          Italic: false,
          Smaller: false
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
          Functions: [
            {
              Function: "abs(x) + sin(x)",
              Decoration: {
                Label: "f",
                Color: ""
              },
              Variable: { Name: xRune, Indice: "" },
              Range: [-5, 5]
            }
          ]
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
    case BlockKind.FigurePointFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
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
          },
          Answer: {
            X: "3",
            Y: "4"
          }
        }
      };
      return out;
    }
    case BlockKind.FigureVectorFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
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
          },
          Answer: {
            X: "3",
            Y: "4"
          },
          MustHaveOrigin: false,
          AnswerOrigin: { X: "", Y: "" }
        }
      };
      return out;
    }
    case BlockKind.VariationTableFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: {
            Xs: ["-5", "0", "5"],
            Fxs: ["-3", "2", "-1"]
          }
        }
      };
      return out;
    }
    case BlockKind.FunctionPointsFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Function: "(x/2)^2",
          Label: "f",
          Variable: { Name: xRune, Indice: "" },
          XGrid: [-4, -2, 0, 2, 4]
        }
      };
      return out;
    }
    case BlockKind.FigureVectorPairFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
            ShowGrid: true,
            Bounds: {
              Width: 10,
              Height: 10,
              Origin: { X: 3, Y: 3 }
            },
            Drawings: { Points: [], Lines: [], Segments: [] }
          },
          Criterion: VectorPairCriterion.VectorColinear
        }
      };
      return out;
    }
    case BlockKind.FigureAffineLineFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
            ShowGrid: true,
            Bounds: {
              Width: 10,
              Height: 10,
              Origin: { X: 3, Y: 3 }
            },
            Drawings: { Points: [], Lines: [], Segments: [] }
          },
          Label: "f",
          A: "1",
          B: "3"
        }
      };
      return out;
    }
    case BlockKind.TreeFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          EventsProposals: ["P", "F", "?"],
          AnswerRoot: {
            Children: [
              {
                Children: [
                  { Value: 0, Children: [], Probabilities: [] },
                  { Value: 1, Children: [], Probabilities: [] }
                ],
                Probabilities: ["0.7", "0.3"],
                Value: 0
              },
              {
                Children: [
                  { Value: 0, Children: [], Probabilities: [] },
                  { Value: 1, Children: [], Probabilities: [] }
                ],
                Probabilities: ["0.7", "0.3"],
                Value: 1
              }
            ],
            Probabilities: ["0.7", "0.3"],
            Value: 0
          }
        }
      };
      return out;
    }
    case BlockKind.TableFieldBlock: {
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
          Answer: [
            ["0", "1"],
            ["2", "3"]
          ]
        }
      };
      return out;
    }
    default:
      throw "Unexpected Kind";
  }
}

export function completePoint(s: string, point: CoordExpression) {
  s = s.trim();
  if (s.startsWith("x_") && s.length > 2) {
    const name = s.substr(2);
    if (!point.Y) {
      point.Y = "y_" + name;
    }
  }
}

export function variableToString(v: Variable) {
  let name = String.fromCodePoint(v.Name);
  if (v.Indice) {
    name += "_" + v.Indice;
  }
  return name;
}

export function tagColor(tag: string) {
  if (
    tag == DifficultyTag.Diff1 ||
    tag == DifficultyTag.Diff2 ||
    tag == DifficultyTag.Diff3
  ) {
    return "secondary";
  }
  return "primary";
}

// returns 0 for question without difficulty
export function questionDifficulty(tags: string[]): number {
  for (const tag of tags) {
    if (tag == DifficultyTag.Diff1) {
      return 1;
    } else if (tag == DifficultyTag.Diff2) {
      return 2;
    } else if (tag == DifficultyTag.Diff3) {
      return 3;
    }
  }
  return 0;
}

export function onDragListItemStart(payload: DragEvent, index: number) {
  payload.dataTransfer?.setData("text/json", JSON.stringify({ index: index }));
  payload.dataTransfer!.dropEffect = "move";
}

/** take the block at the index `origin` and insert it right before
the block at index `target` (which is between 0 and nbBlocks)
 */
export function swapItems<T>(origin: number, target: number, list: T[]) {
  if (target == origin || target == origin + 1) {
    // nothing to do
    return list;
  }

  if (origin < target) {
    const after = list.slice(target);
    const before = list.slice(0, target);
    const originRow = before.splice(origin, 1);
    before.push(...originRow);
    before.push(...after);
    return before;
  } else {
    const before = list.slice(0, target);
    const originRow = list.splice(origin, 1);
    const after = list.slice(target);
    before.push(...originRow);
    before.push(...after);
    return before;
  }
}

/** return the list of tags shared by all the questions */
export function commonTags(questions: QuestionHeader[]) {
  const crible: { [key: string]: number } = {};
  questions.forEach(qu =>
    (qu.Tags || []).forEach(tag => (crible[tag] = (crible[tag] || 0) + 1))
  );
  console.log(crible);

  return Object.entries(crible)
    .filter(entry => entry[1] == questions.length)
    .map(entry => entry[0]);
}
