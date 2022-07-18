import {
  Binary,
  BlockKind,
  ComparisonLevel,
  DifficultyTag,
  ProofAssertionKind,
  SignSymbol,
  TextKind,
  VectorPairCriterion,
  Visibility,
  type Block,
  type CoordExpression,
  type ExpressionFieldBlock,
  type FigureAffineLineFieldBlock,
  type FigureBlock,
  type FigurePointFieldBlock,
  type FigureVectorFieldBlock,
  type FigureVectorPairFieldBlock,
  type FormulaBlock,
  type FunctionPointsFieldBlock,
  type FunctionsGraphBlock,
  type NumberFieldBlock,
  type OrderedListFieldBlock,
  type Origin,
  type ProofAssertion,
  type ProofFieldBlock,
  type QuestionHeader,
  type RadioFieldBlock,
  type SignTableBlock,
  type TableBlock,
  type TableFieldBlock,
  type TextBlock,
  type TextPart,
  type TreeFieldBlock,
  type Variable,
  type VariationTableBlock,
  type VariationTableFieldBlock,
  type VectorFieldBlock,
} from "./api_gen";
import { LevelTag } from "./exercice_gen";

export const ExpressionColor = "orange";

export const colorByKind: { [key in TextKind]: string } = {
  [TextKind.Text]: "",
  [TextKind.StaticMath]: "green",
  [TextKind.Expression]: ExpressionColor,
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
    (chunk) => {
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
    BlockKind.FunctionsGraphBlock,
    {
      label: "Graphes de fonctions",
      isAnswerField: false,
    },
  ],
  [BlockKind.TableBlock, { label: "Tableau", isAnswerField: false }],
  [
    BlockKind.SignTableBlock,
    {
      label: "Tableau de signes",
      isAnswerField: false,
    },
  ],
  [
    BlockKind.VariationTableBlock,
    {
      label: "Tableau de variations",
      isAnswerField: false,
    },
  ],
  [BlockKind.NumberFieldBlock, { label: "Nombre", isAnswerField: true }],
  [
    BlockKind.ExpressionFieldBlock,
    {
      label: "Expression",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.OrderedListFieldBlock,
    {
      label: "Liste ordonnée",
      isAnswerField: true,
    },
  ],
  [BlockKind.RadioFieldBlock, { label: "QCM", isAnswerField: true }],
  [
    BlockKind.FigurePointFieldBlock,
    {
      label: "Point (sur une figure)",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.FigureVectorFieldBlock,
    {
      label: "Vecteur (sur une figure)",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.VectorFieldBlock,
    { label: "Vecteur (numérique)", isAnswerField: true },
  ],
  [
    BlockKind.FunctionPointsFieldBlock,
    {
      label: "Construction de fonction",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.VariationTableFieldBlock,
    {
      label: "Tableau de variations",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.FigureAffineLineFieldBlock,
    {
      label: "Droite (affine)",
      isAnswerField: true,
    },
  ],

  [
    BlockKind.FigureVectorPairFieldBlock,
    {
      label: "Construction de vecteurs",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.ProofFieldBlock,
    { label: "Preuve (à compléter)", isAnswerField: true },
  ],
  [BlockKind.TableFieldBlock, { label: "Tableau", isAnswerField: true }],
  [BlockKind.TreeFieldBlock, { label: "Arbre", isAnswerField: true }],
] as const;

export const BlockKindLabels: {
  [T in BlockKind]: { label: string; isAnswerField: boolean };
} = Object.fromEntries(sortedBlockKindLabels);

interface BlockKindTypes {
  [BlockKind.FigureBlock]: FigureBlock;
  [BlockKind.FormulaBlock]: FormulaBlock;
  [BlockKind.ExpressionFieldBlock]: ExpressionFieldBlock;
  [BlockKind.FunctionsGraphBlock]: FunctionsGraphBlock;
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
  [BlockKind.VectorFieldBlock]: VectorFieldBlock;
  [BlockKind.ProofFieldBlock]: ProofFieldBlock;
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
  vars.forEach((v) => {
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

  return Object.keys(points).filter((name) => {
    const point = points[name];
    return point.x && point.y;
  });
}

export function saveData<T>(data: T, fileName: string) {
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
          Smaller: false,
        },
      };
      return out;
    }
    case BlockKind.FormulaBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Parts: "",
        },
      };
      return out;
    }
    case BlockKind.FigureBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          ShowGrid: true,
          ShowOrigin: true,
          Bounds: {
            Width: 10,
            Height: 10,
            Origin: { X: 3, Y: 3 },
          },
          Drawings: {
            Lines: [],
            Points: [],
            Segments: [],
            Areas: [],
          },
        },
      };
      return out;
    }
    case BlockKind.FunctionsGraphBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          FunctionExprs: [
            {
              Function: "abs(x) + sin(x)",
              Decoration: {
                Label: "f",
                Color: "",
              },
              Variable: { Name: xRune, Indice: "" },
              From: "-5",
              To: "5",
            },
          ],
          FunctionVariations: [],
          Areas: [],
        },
      };
      return out;
    }
    case BlockKind.VariationTableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Label: "f(x)",
          Xs: ["-5", "0", "5"],
          Fxs: ["-3", "2", "-1"],
        },
      };
      return out;
    }
    case BlockKind.SignTableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Label: "f(x)",
          FxSymbols: [
            SignSymbol.Nothing,
            SignSymbol.Zero,
            SignSymbol.ForbiddenValue,
            SignSymbol.Nothing,
          ],
          Xs: ["\\infty", "3", "5", "+\\infty"],
          Signs: [true, false, true],
        },
      };
      return out;
    }
    case BlockKind.TableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          VerticalHeaders: [
            { Kind: TextKind.Text, Content: "Ligne 1" },
            { Kind: TextKind.Text, Content: "Ligne 2" },
          ],
          HorizontalHeaders: [
            { Kind: TextKind.Text, Content: "Colonne 1" },
            { Kind: TextKind.Text, Content: "Colonne 2" },
          ],
          Values: [
            [
              { Kind: TextKind.Text, Content: "Case" },
              { Kind: TextKind.StaticMath, Content: "\\frac{1}{2}" },
            ],
            [
              { Kind: TextKind.Expression, Content: "2x + 3" },
              { Kind: TextKind.StaticMath, Content: "18" },
            ],
          ],
        },
      };
      return out;
    }
    case BlockKind.NumberFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Expression: "1",
        },
      };
      return out;
    }
    case BlockKind.ExpressionFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Label: { Kind: TextKind.Text, Content: "" },
          Expression: "x^2 + 2x + 1",
          ComparisonLevel: ComparisonLevel.SimpleSubstitutions,
        },
      };
      return out;
    }
    case BlockKind.RadioFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: "1",
          Proposals: ["Oui", "Non"],
          AsDropDown: false,
        },
      };
      return out;
    }
    case BlockKind.OrderedListFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: ["$\\{$", "-12", ";", "30", "$\\}$"],
          AdditionalProposals: [],
          Label: "x \\in ",
        },
      };
      return out;
    }
    case BlockKind.FigurePointFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
            ShowGrid: true,
            ShowOrigin: true,
            Bounds: {
              Width: 10,
              Height: 10,
              Origin: { X: 3, Y: 3 },
            },
            Drawings: {
              Lines: [],
              Points: [],
              Segments: [],
              Areas: [],
            },
          },
          Answer: {
            X: "3",
            Y: "4",
          },
        },
      };
      return out;
    }
    case BlockKind.FigureVectorFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
            ShowGrid: true,
            ShowOrigin: true,
            Bounds: {
              Width: 10,
              Height: 10,
              Origin: { X: 3, Y: 3 },
            },
            Drawings: {
              Lines: [],
              Points: [],
              Segments: [],
              Areas: [],
            },
          },
          Answer: {
            X: "3",
            Y: "4",
          },
          MustHaveOrigin: false,
          AnswerOrigin: { X: "", Y: "" },
        },
      };
      return out;
    }
    case BlockKind.VariationTableFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: {
            Label: "f(x)",
            Xs: ["-5", "0", "5"],
            Fxs: ["-3", "2", "-1"],
          },
        },
      };
      return out;
    }
    case BlockKind.FunctionPointsFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Function: "(x/2)^2",
          Label: "C_f",
          Variable: { Name: xRune, Indice: "" },
          XGrid: ["-4", "-2", "0", "2", "4"],
        },
      };
      return out;
    }
    case BlockKind.FigureVectorPairFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
            ShowGrid: true,
            ShowOrigin: true,
            Bounds: {
              Width: 10,
              Height: 10,
              Origin: { X: 3, Y: 3 },
            },
            Drawings: { Points: [], Lines: [], Segments: [], Areas: [] },
          },
          Criterion: VectorPairCriterion.VectorColinear,
        },
      };
      return out;
    }
    case BlockKind.FigureAffineLineFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Figure: {
            ShowGrid: true,
            ShowOrigin: true,
            Bounds: {
              Width: 10,
              Height: 10,
              Origin: { X: 3, Y: 3 },
            },
            Drawings: { Points: [], Lines: [], Segments: [], Areas: [] },
          },
          Label: "f",
          A: "1",
          B: "3",
        },
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
                  { Value: 1, Children: [], Probabilities: [] },
                ],
                Probabilities: ["0.7", "0.3"],
                Value: 0,
              },
              {
                Children: [
                  { Value: 0, Children: [], Probabilities: [] },
                  { Value: 1, Children: [], Probabilities: [] },
                ],
                Probabilities: ["0.7", "0.3"],
                Value: 1,
              },
            ],
            Probabilities: ["0.7", "0.3"],
            Value: 0,
          },
        },
      };
      return out;
    }
    case BlockKind.TableFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          VerticalHeaders: [
            { Kind: TextKind.Text, Content: "Ligne 1" },
            { Kind: TextKind.Text, Content: "Ligne 2" },
          ],
          HorizontalHeaders: [
            { Kind: TextKind.Text, Content: "Colonne 1" },
            { Kind: TextKind.Text, Content: "Colonne 2" },
          ],
          Answer: [
            ["0", "1"],
            ["2", "3"],
          ],
        },
      };
      return out;
    }
    case BlockKind.VectorFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: {
            X: "3.5",
            Y: "-4",
          },
          AcceptColinear: false,
          DisplayColumn: true,
        },
      };
      return out;
    }
    case BlockKind.ProofFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: {
            Parts: [
              {
                Kind: ProofAssertionKind.ProofNode,
                Data: {
                  Op: Binary.And,
                  Left: {
                    Kind: ProofAssertionKind.ProofStatement,
                    Data: { Content: "$n$ est pair" },
                  },
                  Right: {
                    Kind: ProofAssertionKind.ProofStatement,
                    Data: { Content: "$m$ est impair" },
                  },
                },
              },
              {
                Kind: ProofAssertionKind.ProofStatement,
                Data: { Content: "$n+m$ est pair" },
              },
            ],
          },
        },
      };
      return out;
    }
    default:
      throw "Unexpected Kind";
  }
}

/** update the Y field of `point` if it is empty and `s`
 * has the form 'x_A'
 */
export function completePoint(s: string, point: CoordExpression) {
  s = s.trim();
  if (s.startsWith("x_") && s.length > 2) {
    const name = s.substring(2);
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

/** tagString returns a normalized version of the `tag` */
export function tagString(tag: string) {
  return tag
    .trim()
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .toUpperCase();
}

export function tagColor(tag: string) {
  if (
    tag == DifficultyTag.Diff1 ||
    tag == DifficultyTag.Diff2 ||
    tag == DifficultyTag.Diff3
  ) {
    return "secondary";
  }
  if (
    tag == LevelTag.Seconde ||
    tag == LevelTag.Premiere ||
    tag == LevelTag.Terminale
  ) {
    return "pink";
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
  if (payload.dataTransfer == null) return;
  payload.dataTransfer.setData("text/json", JSON.stringify({ index: index }));
  payload.dataTransfer.dropEffect = "move";
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

/** return the list of tags shared by all the list */
export function commonTags(tags: string[][]) {
  const crible: { [key: string]: number } = {};
  tags.forEach((l) =>
    l.forEach((tag) => (crible[tag] = (crible[tag] || 0) + 1))
  );
  return Object.entries(crible)
    .filter((entry) => entry[1] == tags.length)
    .map((entry) => entry[0]);
}

/** return the list of tags shared by all the questions */
export function commonGroupTags(questions: QuestionHeader[]) {
  return commonTags(questions.map((qu) => qu.Tags || []));
}

/** `visiblityColors` exposes the colors used to differentiate ressource visiblity */
export const visiblityColors: { [key in Visibility]: string } = {
  [Visibility.Admin]: "yellow-lighten-3",
  [Visibility.Personnal]: "white",
};

export function removeDuplicates(tags: string[][]) {
  const unique: string[][] = [];
  tags.forEach((l) => {
    if (unique.map((l) => JSON.stringify(l)).includes(JSON.stringify(l))) {
      return;
    }
    unique.push(l);
  });
  return unique;
}

export function personnalOrigin(): Origin {
  return {
    IsPublic: false,
    AllowPublish: false,
    Visibility: Visibility.Personnal,
  };
}

/** lastColorUsed is a shared global parameter used by color pickers */
export const lastColorUsed = { color: "#FF0000" };

export function emptyAssertion(): ProofAssertion {
  return {
    Kind: ProofAssertionKind.ProofInvalid,
    Data: {},
  };
}
