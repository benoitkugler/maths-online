import {
  ParameterEntryKind,
  type ErrParameters,
  type ParameterEntry,
  type Parameters,
  type Rp
} from "@/controller/api_gen";
import { ExpressionColor, variableToString } from "@/controller/editor";
import type { Token } from "../utils/interpolated_text";

function tokenizeLine(line: string): Token[] {
  if (line.startsWith("#")) {
    // we have a comment
    return [{ Content: line + "\n", Kind: "color: green" }];
  }

  const i = line.indexOf("=");
  if (i == -1) {
    // invalid line
    return [{ Content: line + "\n", Kind: "" }];
  }

  const vars = line.substring(0, i);
  const expression = line.substring(i + 1);
  // differentiate between regular var and intrisic
  // by finding the number of variables
  const isIntrinsic = vars.split(",").length > 1;
  const color = isIntrinsic ? "purple" : ExpressionColor;
  return [
    {
      Content: vars,
      Kind: ""
    },
    {
      Content: "=",
      Kind: "font-weight: bold"
    },
    {
      Content: `${expression}\n`,
      Kind: `color: ${color}`
    }
  ];
}

/** Highlights the parameters text, accepting invalid entries */
export function tokenize(text: string): Token[] {
  const out: Token[] = [];
  const lines = text.split("\n");
  for (const line of lines) {
    out.push(...tokenizeLine(line));
  }

  // remove the last \n
  const tk = out[out.length - 1];
  tk.Content = tk.Content.substring(0, tk.Content.length - 1);

  return out;
}

/** Transforms the raw text into a list of parameter entries */
export function parseParameters(text: string): {
  params: ParameterEntry[];
  error: ErrParameters | null;
} {
  const out: ParameterEntry[] = [];
  const lines = text.split("\n");
  let currentComment = "";
  for (const line of lines) {
    // ignore empty lines
    if (!line.trim().length) continue;

    if (line.startsWith("#")) {
      // we have a comment
      currentComment += line.substring(1) + "\n";
    } else {
      // close the current comment or expr if needed
      if (currentComment.length) {
        out.push({
          Kind: ParameterEntryKind.Co,
          Data: currentComment.substring(0, currentComment.length - 1)
        });
        currentComment = "";
      }

      const i = line.indexOf("=");
      if (i == -1) {
        return {
          params: [],
          error: {
            Origin: line,
            Details: `Définition invalide (symbole '#' ou '=' manquant)`
          }
        };
      }
      const vars = line.substring(0, i).trim();
      const expression = line.substring(i + 1).trim();
      // differentiate between regular var and intrisic
      // by finding the number of variables
      if (vars.split(",").length > 1) {
        // we have an intrinsic
        out.push({ Kind: ParameterEntryKind.In, Data: line });
      } else {
        // we have a single variable
        const rp: Rp = { variable: { Name: 0, Indice: "" }, expression: "" };
        const varsParts = vars.split("_", 2);
        const name = varsParts[0].trim();
        if (name.length != 1) {
          return {
            params: [],
            error: {
              Origin: line,
              Details: `Nom de variable invalide : ${vars}`
            }
          };
        }
        rp.variable.Name = name.charCodeAt(0);
        rp.variable.Indice = varsParts[1] || "";
        rp.expression = expression;
        out.push({ Kind: ParameterEntryKind.Rp, Data: rp });
      }
    }
  }

  // close the current comment or expr if needed
  if (currentComment.length) {
    out.push({
      Kind: ParameterEntryKind.Co,
      Data: currentComment.substring(0, currentComment.length - 1)
    });
    currentComment = "";
  }

  return { params: out, error: null };
}

/** Returns the raw text that should be displayed, enforcing formatting */
export function parametersToString(params: Parameters) {
  params = params || [];
  const lines: string[] = [];

  params.forEach(entry => {
    let rp: Rp;
    let s: string;
    switch (entry.Kind) {
      case ParameterEntryKind.Co:
        s = entry.Data as string;
        lines.push(...s.split("\n").map(l => "# " + l.trim()));
        return;
      case ParameterEntryKind.Rp:
        rp = entry.Data as Rp;
        lines.push(`${variableToString(rp.variable)} = ${rp.expression}`);
        return;
      case ParameterEntryKind.In:
        s = entry.Data as string;
        lines.push(s);
        return;
    }
  });
  return lines.join("\n");
}
