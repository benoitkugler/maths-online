import type { TextPart } from "./exercice_gen";
import { BlockKind, TextKind } from "./exercice_gen";

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
    console.log(s.substring(outerStart, outerEnd));

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
  [BlockKind.FormulaFieldBlock]: "FormulaFieldBlock",
  [BlockKind.FunctionGraphBlock]: "FunctionGraphBlock",
  [BlockKind.FunctionVariationGraphBlock]: "FunctionVariationGraphBlock",
  [BlockKind.ListField]: "ListField",
  [BlockKind.NumberFieldBlock]: "NumberFieldBlock",
  [BlockKind.RadioFieldBlock]: "RadioFieldBlock",
  [BlockKind.SignTableBlock]: "SignTableBlock",
  [BlockKind.TableBlock]: "TableBlock",
  [BlockKind.VariationTableBlock]: "VariationTableBlock"
};
