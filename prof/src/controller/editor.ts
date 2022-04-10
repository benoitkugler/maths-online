import type {
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
import { BlockKind, TextKind } from "./exercice_gen";

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
