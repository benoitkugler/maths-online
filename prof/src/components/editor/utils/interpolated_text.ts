import { TextKind, type TextPart } from "@/controller/api_gen";
import { colorByKind } from "@/controller/editor";

/** Token generates a span with given style and text content */
export interface Token {
  Content: string;
  Kind: string;
}

const reLaTeX = /\$([^$]+)\$/g;
const reExpression = /&([^&]+)&/g;

export function splitByRegexp<T>(
  re: RegExp,
  s: string,
  kindMatch: T,
  kindDefault: T
): { Content: string; Kind: T }[] {
  const out: { Content: string; Kind: T }[] = [];
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

export function partToToken(part: TextPart): Token {
  return {
    Content: part.Content,
    Kind: styles[part.Kind]
  };
}

export function defautTokenize(input: string): Token[] {
  return itemize(input).map(partToToken);
}

const styles = {
  [TextKind.Text]: "",
  [TextKind.StaticMath]: `color: ${colorByKind[TextKind.StaticMath]}`,
  [TextKind.Expression]: `color: ${
    colorByKind[TextKind.Expression]
  }; font-weight: bold;`
} as const;
