import {
  Binary,
  BlockKind,
  ComparisonLevel,
  DifficultyTag,
  ProofAssertionKind,
  Section,
  SignSymbol,
  TextKind,
  Visibility,
  type Block,
  type CoordExpression,
  type ExpressionFieldBlock,
  type FigureBlock,
  type FormulaBlock,
  type FunctionPointsFieldBlock,
  type FunctionsGraphBlock,
  type NumberFieldBlock,
  type OrderedListFieldBlock,
  type Origin,
  type ProofAssertion,
  type ProofFieldBlock,
  type RadioFieldBlock,
  type SignTableBlock,
  type SignTableFieldBlock,
  type TableBlock,
  type TableFieldBlock,
  type Tags,
  type TagsDB,
  type TextBlock,
  type TreeFieldBlock,
  type Variable,
  type VariationTableBlock,
  type VariationTableFieldBlock,
  type VectorFieldBlock,
  type TreeNodeAnswer,
  type TreeBlock,
  type GeometricConstructionFieldBlock,
  GeoFieldKind,
  FiguresOrGraphsKind,
  PublicStatus,
  type QuestiongroupExt,
  type ExercicegroupExt,
  Int,
  Parameters,
  Enonce,
  ErrQuestionInvalid,
  SetFieldBlock,
} from "./api_gen";
import { LoopbackServerEvent } from "./loopback_gen";
import { copy } from "./utils";

export const ExpressionColor = "orange";

export const colorByKind: { [key in TextKind]: string } = {
  [TextKind.Text]: "",
  [TextKind.StaticMath]: "green",
  [TextKind.Expression]: ExpressionColor,
};

export const sortedBlockKindLabels = [
  [BlockKind.TextBlock, { label: "Texte", isAnswerField: false }],
  [BlockKind.FormulaBlock, { label: "Formule", isAnswerField: false }],
  [
    BlockKind.FigureBlock,
    { label: "Figure (droites, cercles)", isAnswerField: false },
  ],
  [
    BlockKind.FunctionsGraphBlock,
    {
      label: "Graphes (fonctions, suites)",
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
  [
    BlockKind.TreeBlock,
    {
      label: "Arbre (de probas)",
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
    BlockKind.GeometricConstructionFieldBlock,
    {
      label: "Constructions (point, vecteur, droite)",
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
      label: "Graphe de fonction (ou suite)",
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
    BlockKind.SignTableFieldBlock,
    {
      label: "Tableau de signes",
      isAnswerField: true,
    },
  ],
  [
    BlockKind.ProofFieldBlock,
    { label: "Preuve (à compléter)", isAnswerField: true },
  ],
  [BlockKind.TableFieldBlock, { label: "Tableau", isAnswerField: true }],
  [BlockKind.TreeFieldBlock, { label: "Arbre", isAnswerField: true }],
  [BlockKind.SetFieldBlock, { label: "Ensembles", isAnswerField: true }],
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
  [BlockKind.GeometricConstructionFieldBlock]: GeometricConstructionFieldBlock;
  [BlockKind.VariationTableFieldBlock]: VariationTableFieldBlock;
  [BlockKind.SignTableFieldBlock]: SignTableFieldBlock;
  [BlockKind.FunctionPointsFieldBlock]: FunctionPointsFieldBlock;
  [BlockKind.TreeBlock]: TreeBlock;
  [BlockKind.TreeFieldBlock]: TreeFieldBlock;
  [BlockKind.TableFieldBlock]: TableFieldBlock;
  [BlockKind.VectorFieldBlock]: VectorFieldBlock;
  [BlockKind.ProofFieldBlock]: ProofFieldBlock;
  [BlockKind.SetFieldBlock]: SetFieldBlock;
}

export interface TypedBlock<K extends BlockKind> {
  Kind: K;
  Data: BlockKindTypes[K];
}

export const xRune = "x".codePointAt(0)! as Int;
export const yRune = "y".codePointAt(0)! as Int;
export const nRune = "n".codePointAt(0)! as Int;

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

const signTableExample: SignTableBlock = {
  Xs: ["-inf", "3", "5", "+inf"],
  Functions: [
    {
      Label: "f(x)",
      FxSymbols: [
        SignSymbol.Nothing,
        SignSymbol.Zero,
        SignSymbol.ForbiddenValue,
        SignSymbol.Nothing,
      ],
      Signs: [true, false, true],
    },
    {
      Label: "g(x)",
      FxSymbols: [
        SignSymbol.ForbiddenValue,
        SignSymbol.Nothing,
        SignSymbol.Zero,
        SignSymbol.Nothing,
      ],
      Signs: [false, false, true],
    },
  ],
};

function newTreeNode(value: number): TreeNodeAnswer {
  return {
    Value: value as Int,
    Children: [],
    Probabilities: [],
  };
}

function newTree() {
  return {
    EventsProposals: ["P", "F", "?"],
    AnswerRoot: {
      Children: [
        {
          Children: [newTreeNode(0), newTreeNode(1), newTreeNode(2)],
          Probabilities: ["0.7", "0.2", "0.1"],
          Value: 0 as Int,
        },
        {
          Children: [newTreeNode(0), newTreeNode(1), newTreeNode(2)],
          Probabilities: ["0.7", "0.2", "0.1"],
          Value: 1 as Int,
        },
      ],
      Probabilities: ["0.7", "0.3"],
      Value: 0 as Int,
    },
  };
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
            Width: 10 as Int,
            Height: 10 as Int,
            Origin: { X: 3, Y: 3 },
          },
          Drawings: {
            Lines: [],
            Points: [],
            Segments: [],
            Circles: [],
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
          SequenceExprs: [
            {
              Function: "(-1)^n * (n/2)^2",
              Decoration: {
                Label: "u_n",
                Color: "#00FF00",
              },
              Variable: { Name: nRune, Indice: "" },
              From: "-6",
              To: "6",
            },
          ],
          FunctionVariations: [],
          Areas: [],
          Points: [],
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
        Data: copy(signTableExample),
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
          Label: "",
          Expression: "x^2 + 2x + 1",
          ComparisonLevel: ComparisonLevel.SimpleSubstitutions,
          ShowFractionHelp: false,
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
    case BlockKind.GeometricConstructionFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Field: {
            Kind: GeoFieldKind.GFPoint,
            Data: {
              Answer: {
                X: "3",
                Y: "4",
              },
            },
          },
          Background: {
            Kind: FiguresOrGraphsKind.FigureBlock,
            Data: {
              ShowGrid: true,
              ShowOrigin: true,
              Bounds: {
                Width: 10 as Int,
                Height: 10 as Int,
                Origin: { X: 3, Y: 3 },
              },
              Drawings: {
                Lines: [],
                Points: [],
                Segments: [],
                Circles: [],
                Areas: [],
              },
            },
          },
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
    case BlockKind.SignTableFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: copy(signTableExample),
        },
      };
      return out;
    }
    case BlockKind.FunctionPointsFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          IsDiscrete: false,
          Function: "(x/2)^2",
          Label: "C_f",
          Variable: { Name: xRune, Indice: "" },
          XGrid: ["-4", "-2", "0", "2", "4"],
        },
      };
      return out;
    }
    case BlockKind.TreeFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: { Answer: newTree() },
      };
      return out;
    }
    case BlockKind.TreeBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: newTree(),
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
                Data: { Content: "$n+m$ est impair" },
              },
            ],
          },
        },
      };
      return out;
    }
    case BlockKind.SetFieldBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Answer: "¬A ∩ B ",
          AdditionalSets: ["C"],
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

export const MatiereColor = "orange-lighten-2";
export const LevelColor = "pink";
export const SubLevelColor = "purple";
export const ChapterColor = "primary-darken-1";
export const IsyTrivColor = "brown";
export const DifficultyColor = "secondary-darken-1";

export function tagColor(tag: Section | undefined) {
  switch (tag) {
    case Section.Level:
      return LevelColor;
    case Section.Chapter:
      return ChapterColor;
    case Section.TrivMath:
      return IsyTrivColor;
    case Section.SubLevel:
      return SubLevelColor;
    case Section.Matiere:
      return MatiereColor;
  }
  // difficulty
  return DifficultyColor;
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
    PublicStatus: PublicStatus.NotAdmin,
    Visibility: Visibility.Personnal,
    IsInReview: { InReview: false, Id: -1 as Int },
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

export interface VariantG {
  Id: Int;
  Subtitle: string;
  Difficulty: DifficultyTag;
  HasCorrection?: boolean;
}

/** filterTags returns the first tags matching query.
 * `blackList` may be provided to exclude results.
 * If `query` is empty, the first items are not returned.
 */
export function filterTags(
  candidates: string[],
  query: string,
  blackList: string[]
) {
  const pagination = 6;
  const blackListSet = new Set(blackList);
  const out: string[] = [];
  query = query.toUpperCase();
  for (const candidate of candidates) {
    if (out.length >= pagination) {
      break;
    }
    if (blackListSet.has(candidate)) {
      continue;
    }
    const start = candidate.indexOf(query);
    if (query == "" || start != -1) {
      out.push(candidate);
    }
  }
  return out;
}

export function emptyTagsDB(): TagsDB {
  return {
    Levels: [],
    ChaptersByLevel: {},
    TrivByChapters: {},
    SubLevelsByLevel: {},
  };
}

/** areTagsEquals compares the tags without taking order in account */
export function areTagsEquals(tags1: Tags, tags2: Tags) {
  const l1 = (tags1 || []).map((ts) => `${ts.Section}--${ts.Tag}`);
  l1.sort();
  const l2 = (tags2 || []).map((ts) => `${ts.Section}--${ts.Tag}`);
  l2.sort();
  return l1.join(";") == l2.join(";");
}

/** either a questiongroup or an exercicegroup */
export interface ResourceGroup {
  Id: Int;
  Title: string;
  Variants: VariantG[];
  Tags: Tags;
  Origin: Origin;
}

export function questionToResource(group: QuestiongroupExt): ResourceGroup {
  return {
    Id: group.Group.Id,
    Title: group.Group.Title,
    Tags: group.Tags,
    Variants: group.Variants || [],
    Origin: group.Origin,
  };
}

export function exerciceToResource(group: ExercicegroupExt): ResourceGroup {
  return {
    Id: group.Group.Id,
    Title: group.Group.Title,
    Tags: group.Tags,
    Variants: group.Variants || [],
    Origin: group.Origin,
  };
}

export interface QuestionPage {
  id: number;
  parameters: Parameters;
  sharedParameters: Parameters;
  enonce: Enonce;
  correction: Enonce;
}

export interface SaveQuestionOut {
  IsValid: boolean;
  Error: ErrQuestionInvalid;
  Preview: LoopbackServerEvent;
}
