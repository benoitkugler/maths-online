import { TextKind, type TextPart } from "@/controller/api_gen";
import { colorByKind } from "@/controller/editor";

/** Token generates a span with given style and text content.
 * The content may contain line breaks
 */
export interface Token {
  Content: string;
  Kind: string; // usually a style string
}

const reLaTeX = /\$([^$]+)\$/g;
const reExpression = /&([^&]+)&/g;
const reNumberField = /#([^#]+)#/g;

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

export function defautTokenize(input: string, allowNumberField = false): Token[] {
  if (allowNumberField) {
    const chunks = splitByRegexp(reNumberField, input, true, false);
    const out: Token[] = []
    for (const chunk of chunks) {
      if (chunk.Kind) { // number field
        out.push({ Content: chunk.Content, Kind: styles.numberField })
      } else { // regular 
        out.push(...itemize(chunk.Content).map(partToToken))
      }
    }
    return out;
  } else {
    return itemize(input).map(partToToken);
  }
}

const styles = {
  [TextKind.Text]: "",
  [TextKind.StaticMath]: `color: ${colorByKind[TextKind.StaticMath]}`,
  [TextKind.Expression]: `color: ${colorByKind[TextKind.Expression]
    }; font-weight: bold;`,
  numberField: "color: pink; font-weight: bold;",
} as const;

const TokenNewLine = "__newLine" as const;

/** return a list of lines, where each line does not contain line break anymore */
export function splitNewLines(tokens: Token[]): Token[][] {
  const out: Token[][] = [];
  let currentLine: Token[] = [];
  tokens.forEach(origin => {
    const splitted = splitOneNewLines(origin);
    splitted.forEach(token => {
      if (token.Kind == TokenNewLine) {
        out.push(currentLine);
        currentLine = [];
      } else {
        currentLine.push(token);
      }
    });
  });
  if (currentLine.length) {
    out.push(currentLine);
  }
  return out;
}

function splitOneNewLines(token: Token): Token[] {
  const out: Token[] = [];
  token.Content.split("\n").forEach(text =>
    out.push(
      {
        Content: text,
        Kind: token.Kind
      },
      { Content: "", Kind: TokenNewLine }
    )
  );
  return out.slice(0, out.length - 1);
}
